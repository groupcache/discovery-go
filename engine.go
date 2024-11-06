/*
 * Copyright 2024 Arsene Tochemey Gandote
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      https://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package discoverygo

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"strconv"
	"sync"
	"time"

	goset "github.com/deckarep/golang-set/v2"
	"github.com/flowchartsman/retry"
	groupcache "github.com/groupcache/groupcache-go/v3"
	"github.com/groupcache/groupcache-go/v3/transport/peer"
	"github.com/hashicorp/memberlist"
	"go.uber.org/atomic"

	"github.com/groupcache/discovery-go/discovery"
	"github.com/groupcache/discovery-go/internal/errorschain"
	"github.com/groupcache/discovery-go/logger"
)

// Engine defines the discovery engine
type Engine struct {
	mconfig              *memberlist.Config
	mlist                *memberlist.Memberlist
	started              *atomic.Bool
	maxJoinAttempts      int
	shutdownTimeout      time.Duration
	maxJoinTimeout       time.Duration
	maxJoinRetryInterval time.Duration
	// specifies the discovery provider
	provider discovery.Provider
	// specifies the minimum number of cluster members
	// the default values is 1
	minimumPeersQuorum    uint
	stopEventsListenerSig chan struct{}
	eventsLock            *sync.Mutex
	lock                  *sync.Mutex
	// specifies the  hostNode
	hostNode *Peer

	logger logger.Logger

	deamon *groupcache.Daemon
}

// NewEngine creates an instance of the discovery Engine
func NewEngine(name string, provider discovery.Provider, daemon *groupcache.Daemon, host *Peer, opts ...Option) *Engine {
	// create an instance of Group
	engine := &Engine{
		provider:              provider,
		shutdownTimeout:       3 * time.Second,
		minimumPeersQuorum:    1,
		maxJoinTimeout:        time.Second,
		maxJoinRetryInterval:  time.Second,
		maxJoinAttempts:       10,
		lock:                  new(sync.Mutex),
		eventsLock:            new(sync.Mutex),
		stopEventsListenerSig: make(chan struct{}, 1),
		hostNode:              host,
		deamon:                daemon,
	}
	// apply the various options
	for _, opt := range opts {
		opt.Apply(engine)
	}

	return engine
}

// Start starts the Engine.
func (g *Engine) Start(ctx context.Context) (err error) {
	// extract the host and port of the host
	host, _, err := getHostPort(g.hostNode.Address)
	if err != nil {
		return fmt.Errorf("invalid host address: %w", err)
	}

	// create the memberlist configuration
	g.mconfig = memberlist.DefaultLANConfig()
	g.mconfig.BindAddr = host
	g.mconfig.BindPort = g.hostNode.DiscoveryPort
	g.mconfig.AdvertisePort = g.hostNode.DiscoveryPort
	g.mconfig.LogOutput = newLogWriter(g.logger)
	g.mconfig.Name = net.JoinHostPort(host, strconv.Itoa(g.hostNode.DiscoveryPort))

	// get the delegate
	delegate, err := g.newEngineDelegate()
	if err != nil {
		return fmt.Errorf("failed to create the discovery engine delegate: %w", err)
	}

	// set the delegate
	g.mconfig.Delegate = delegate

	// start process
	if err := errorschain.
		New(errorschain.ReturnFirst()).
		AddError(g.provider.Initialize()).
		AddError(g.provider.Register()).
		AddError(g.joinCluster(ctx)).
		Error(); err != nil {
		return err
	}

	// create enough buffer to house the cluster events
	// TODO: revisit this number
	eventsCh := make(chan memberlist.NodeEvent, 256)
	g.mconfig.Events = &memberlist.ChannelEventDelegate{
		Ch: eventsCh,
	}

	g.started.Store(true)
	// start listening to events
	go g.eventsListener(eventsCh)

	return nil
}

// Stop stops the Engine gracefully
func (g *Engine) Stop(ctx context.Context) error {
	if !g.started.Load() {
		return nil
	}

	// create a cancellation context
	ctx, cancelFn := context.WithTimeout(ctx, g.shutdownTimeout)
	defer cancelFn()

	// stop the events loop
	close(g.stopEventsListenerSig)

	if err := errorschain.
		New(errorschain.ReturnFirst()).
		AddError(g.mlist.Leave(g.shutdownTimeout)).
		AddError(g.provider.Deregister()).
		AddError(g.provider.Close()).
		AddError(g.mlist.Shutdown()).
		Error(); err != nil {
		return err
	}

	return nil
}

// Peers returns a channel containing the list of peers at a given time
func (g *Engine) Peers(ctx context.Context) ([]*Peer, error) {
	g.lock.Lock()
	members := g.mlist.Members()
	g.lock.Unlock()
	peers := make([]*Peer, 0, len(members))
	for _, member := range members {
		peer := new(Peer)
		if err := json.Unmarshal(member.Meta, &peer); err != nil {
			return nil, err
		}

		if peer != nil && !peer.IsSelf {
			peers = append(peers, peer)
		}
	}
	return peers, nil
}

// eventsListener listens to cluster events
func (g *Engine) eventsListener(eventsChan chan memberlist.NodeEvent) {
	for {
		select {
		case <-g.stopEventsListenerSig:
			// finish listening to cluster events
			return
		case event := <-eventsChan:
			var peerInfo *Peer
			if err := json.Unmarshal(event.Node.Meta, &peerInfo); err != nil {
				// TODO: add a logger to log errors
				continue
			}

			if peerInfo.IsSelf {
				continue
			}

			ctx := context.Background()
			// we need to add the new peers
			currentPeers, _ := g.Peers(ctx)
			peersSet := goset.NewSet[peer.Info]()
			peersSet.Add(peer.Info{
				Address: g.hostNode.Address,
				IsSelf:  true,
			})

			for _, xpeer := range currentPeers {
				peersSet.Add(peer.Info{
					Address: xpeer.Address,
					IsSelf:  xpeer.IsSelf,
				})
			}

			switch event.Event {
			case memberlist.NodeJoin:
				// add the joined node to the peers list
				// and set it to the daemon
				g.eventsLock.Lock()
				peersSet.Add(peer.Info{
					Address: peerInfo.Address,
					IsSelf:  peerInfo.IsSelf,
				})

				_ = g.deamon.SetPeers(ctx, peersSet.ToSlice())
				g.eventsLock.Unlock()

			case memberlist.NodeLeave:
				// remove the left node from the peers list
				// and set it to the daemon
				g.eventsLock.Lock()
				peersSet.Remove(peer.Info{
					Address: peerInfo.Address,
					IsSelf:  peerInfo.IsSelf,
				})
				_ = g.deamon.SetPeers(ctx, peersSet.ToSlice())
				g.eventsLock.Unlock()

			case memberlist.NodeUpdate:
				// TODO: need to handle that later
				continue
			}
		}
	}
}

// joinCluster attempts to join an existing cluster if peers are provided
func (g *Engine) joinCluster(ctx context.Context) error {
	var err error
	g.mlist, err = memberlist.Create(g.mconfig)
	if err != nil {
		return err
	}

	discoveryCtx, discoveryCancel := context.WithTimeout(ctx, g.maxJoinTimeout)
	defer discoveryCancel()

	var peers []string
	retrier := retry.NewRetrier(g.maxJoinAttempts, g.maxJoinRetryInterval, g.maxJoinRetryInterval)
	if err := retrier.RunContext(discoveryCtx, func(ctx context.Context) error {
		peers, err = g.provider.DiscoverPeers()
		if err != nil {
			return err
		}
		return nil
	}); err != nil {
		return err
	}

	if len(peers) > 0 {
		// check whether the cluster quorum is met to operate
		if g.minimumPeersQuorum < uint(len(peers)) {
			return ErrClusterQuorum
		}

		// attempt to join
		joinCtx, joinCancel := context.WithTimeout(ctx, g.maxJoinTimeout)
		defer joinCancel()
		joinRetrier := retry.NewRetrier(g.maxJoinAttempts, g.maxJoinRetryInterval, g.maxJoinRetryInterval)
		if err := joinRetrier.RunContext(joinCtx, func(ctx context.Context) error {
			if _, err := g.mlist.Join(peers); err != nil {
				return err
			}
			return nil
		}); err != nil {
			return err
		}
	}

	return nil
}

// getHostPort returns the actual ip address and port from a given address
func getHostPort(address string) (string, int, error) {
	// Get the address
	addr, err := net.ResolveTCPAddr("tcp", address)
	if err != nil {
		return "", 0, err
	}

	return addr.IP.String(), addr.Port, nil
}
