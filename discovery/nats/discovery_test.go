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

package nats

import (
	"fmt"
	"net"
	"strconv"
	"testing"
	"time"

	natsserver "github.com/nats-io/nats-server/v2/server"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/travisjeffery/go-dynaport"
	"go.uber.org/atomic"

	"github.com/groupcache/discovery-go/discovery"
)

func startNatsServer(t *testing.T) *natsserver.Server {
	t.Helper()
	serv, err := natsserver.NewServer(&natsserver.Options{
		Host: "127.0.0.1",
		Port: -1,
	})

	require.NoError(t, err)

	ready := make(chan bool)
	go func() {
		ready <- true
		serv.Start()
	}()
	<-ready

	if !serv.ReadyForConnections(2 * time.Second) {
		t.Fatalf("nats-io server failed to start")
	}

	return serv
}

func newPeer(t *testing.T, serverAddr string) *Discovery {
	// generate the ports for the single node
	nodePorts := dynaport.Get(1)
	gossipPort := nodePorts[0]

	// create a Cluster node
	host := "127.0.0.1"
	// create the various config option
	natsSubject := "some-subject"

	// create the config
	config := &Config{
		Server:        fmt.Sprintf("nats://%s", serverAddr),
		Subject:       natsSubject,
		Host:          host,
		DiscoveryPort: gossipPort,
	}

	// create the instance of provider
	provider := NewDiscovery(config)

	// initialize
	err := provider.Initialize()
	require.NoError(t, err)
	// return the provider
	return provider
}

func TestDiscovery(t *testing.T) {
	t.Run("With a new instance", func(t *testing.T) {
		// start the NATS server
		srv := startNatsServer(t)

		// generate the ports for the single node
		nodePorts := dynaport.Get(1)
		gossipPort := nodePorts[0]

		// create a Cluster node
		host := "127.0.0.1"
		// create the various config option
		natsSubject := "some-subject"

		serverAddr := srv.Addr().String()
		// create the config
		config := &Config{
			Server:        fmt.Sprintf("nats://%s", serverAddr),
			Subject:       natsSubject,
			Host:          host,
			DiscoveryPort: gossipPort,
		}

		// create the instance of provider
		provider := NewDiscovery(config)
		require.NotNil(t, provider)
		// assert that provider implements the Discovery interface
		// this is a cheap test
		// assert the type of svc
		assert.IsType(t, &Discovery{}, provider)
		var p interface{} = provider
		_, ok := p.(discovery.Provider)
		assert.True(t, ok)
	})
	t.Run("With ID assertion", func(t *testing.T) {
		// start the NATS server
		srv := startNatsServer(t)

		// generate the ports for the single node
		nodePorts := dynaport.Get(1)
		gossipPort := nodePorts[0]
		// create a Cluster node
		host := "127.0.0.1"
		// create the various config option
		natsSubject := "some-subject"

		// create the config
		config := &Config{
			Server:        fmt.Sprintf("nats://%s", srv.Addr().String()),
			Subject:       natsSubject,
			Host:          host,
			DiscoveryPort: gossipPort,
		}

		// create the instance of provider
		provider := NewDiscovery(config)
		require.NotNil(t, provider)
		assert.Equal(t, "nats", provider.ID())
	})
	t.Run("With Initialize", func(t *testing.T) {
		// start the NATS server
		srv := startNatsServer(t)

		// generate the ports for the single node
		nodePorts := dynaport.Get(1)
		gossipPort := nodePorts[0]

		// create a Cluster node
		host := "127.0.0.1"

		// create the various config option
		natsSubject := "some-subject"

		// create the config
		config := &Config{
			Server:        fmt.Sprintf("nats://%s", srv.Addr().String()),
			Subject:       natsSubject,
			Host:          host,
			DiscoveryPort: gossipPort,
		}

		// create the instance of provider
		provider := NewDiscovery(config)

		// initialize
		err := provider.Initialize()
		assert.NoError(t, err)

		// stop the NATS server
		t.Cleanup(srv.Shutdown)
	})
	t.Run("With Initialize: already initialized", func(t *testing.T) {
		// start the NATS server
		srv := startNatsServer(t)

		// generate the ports for the single node
		nodePorts := dynaport.Get(1)
		gossipPort := nodePorts[0]

		// create a Cluster node
		host := "127.0.0.1"

		// create the various config option
		natsSubject := "some-subject"

		// create the config
		config := &Config{
			Server:        fmt.Sprintf("nats://%s", srv.Addr().String()),
			Subject:       natsSubject,
			Host:          host,
			DiscoveryPort: gossipPort,
		}

		// create the instance of provider
		provider := NewDiscovery(config)
		provider.initialized = atomic.NewBool(true)
		assert.Error(t, provider.Initialize())
	})
	t.Run("With Register: already registered", func(t *testing.T) {
		// start the NATS server
		srv := startNatsServer(t)

		// generate the ports for the single node
		nodePorts := dynaport.Get(1)
		gossipPort := nodePorts[0]

		// create a Cluster node
		host := "127.0.0.1"

		// create the various config option
		natsServer := srv.Addr().String()
		natsSubject := "some-subject"

		// create the config
		config := &Config{
			Server:        fmt.Sprintf("nats://%s", natsServer),
			Subject:       natsSubject,
			Host:          host,
			DiscoveryPort: gossipPort,
		}

		// create the instance of provider
		provider := NewDiscovery(config)
		provider.registered = atomic.NewBool(true)
		err := provider.Register()
		assert.Error(t, err)
		assert.EqualError(t, err, discovery.ErrAlreadyRegistered.Error())
	})
	t.Run("With Deregister: already not registered", func(t *testing.T) {
		// start the NATS server
		srv := startNatsServer(t)

		// generate the ports for the single node
		nodePorts := dynaport.Get(1)
		gossipPort := nodePorts[0]

		// create a Cluster node
		host := "127.0.0.1"

		natsServer := srv.Addr().String()
		natsSubject := "some-subject"

		// create the config
		config := &Config{
			Server:        fmt.Sprintf("nats://%s", natsServer),
			Subject:       natsSubject,
			Host:          host,
			DiscoveryPort: gossipPort,
		}

		// create the instance of provider
		provider := NewDiscovery(config)
		err := provider.Deregister()
		assert.Error(t, err)
		assert.EqualError(t, err, discovery.ErrNotRegistered.Error())
	})
	t.Run("With DiscoverPeers", func(t *testing.T) {
		// start the NATS server
		srv := startNatsServer(t)
		// create two peers
		client1 := newPeer(t, srv.Addr().String())
		client2 := newPeer(t, srv.Addr().String())

		// no discovery is allowed unless registered
		peers, err := client1.DiscoverPeers()
		require.Error(t, err)
		assert.EqualError(t, err, discovery.ErrNotRegistered.Error())
		require.Empty(t, peers)

		peers, err = client2.DiscoverPeers()
		require.Error(t, err)
		assert.EqualError(t, err, discovery.ErrNotRegistered.Error())
		require.Empty(t, peers)

		// register client 2
		require.NoError(t, client2.Register())
		peers, err = client2.DiscoverPeers()
		require.NoError(t, err)
		require.Empty(t, peers)

		// register client 1
		require.NoError(t, client1.Register())
		peers, err = client1.DiscoverPeers()
		require.NoError(t, err)
		require.NotEmpty(t, peers)
		require.Len(t, peers, 1)

		discoveredNodeAddr := net.JoinHostPort(client2.config.Host, strconv.Itoa(int(client2.config.DiscoveryPort)))
		require.Equal(t, peers[0], discoveredNodeAddr)

		// discover more peers from client 2
		peers, err = client2.DiscoverPeers()
		require.NoError(t, err)
		require.NotEmpty(t, peers)
		require.Len(t, peers, 1)

		discoveredNodeAddr = net.JoinHostPort(client1.config.Host, strconv.Itoa(int(client1.config.DiscoveryPort)))
		require.Equal(t, peers[0], discoveredNodeAddr)

		// de-register client 2 but it can see client1
		require.NoError(t, client2.Deregister())
		peers, err = client2.DiscoverPeers()
		require.NoError(t, err)
		require.NotEmpty(t, peers)
		discoveredNodeAddr = net.JoinHostPort(client1.config.Host, strconv.Itoa(int(client1.config.DiscoveryPort)))
		require.Equal(t, peers[0], discoveredNodeAddr)

		// client-1 cannot see the deregistered client
		peers, err = client1.DiscoverPeers()
		require.NoError(t, err)
		require.Empty(t, peers)

		require.NoError(t, client1.Close())
		require.NoError(t, client2.Close())

		// stop the NATS server
		t.Cleanup(srv.Shutdown)
	})
	t.Run("With DiscoverPeers: not initialized", func(t *testing.T) {
		// start the NATS server
		srv := startNatsServer(t)

		// generate the ports for the single node
		nodePorts := dynaport.Get(1)
		gossipPort := nodePorts[0]

		// create a Cluster node
		host := "127.0.0.1"

		// create the various config option
		natsServer := srv.Addr().String()
		natsSubject := "some-subject"

		// create the config
		config := &Config{
			Server:        fmt.Sprintf("nats://%s", natsServer),
			Subject:       natsSubject,
			Host:          host,
			DiscoveryPort: gossipPort,
		}

		provider := NewDiscovery(config)
		peers, err := provider.DiscoverPeers()
		assert.Error(t, err)
		assert.Empty(t, peers)
		assert.EqualError(t, err, discovery.ErrNotInitialized.Error())
	})
	t.Run("With Initialize: invalid config", func(t *testing.T) {
		config := &Config{}
		provider := NewDiscovery(config)

		// initialize
		err := provider.Initialize()
		assert.Error(t, err)
	})
}