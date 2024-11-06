package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"time"

	groupcache "github.com/groupcache/groupcache-go/v3"
	"github.com/groupcache/groupcache-go/v3/transport"
	"github.com/groupcache/groupcache-go/v3/transport/peer"

	discoverygo "github.com/groupcache/discovery-go"
	"github.com/groupcache/discovery-go/discovery/dnssd"
)

func ExampleDNSSD() {
	const domainName = "accounts"

	// create an instance of
	config := dnssd.Config{
		DomainName: domainName,
	}

	provider := dnssd.NewDiscovery(&config)

	hostAddress := "192.168.1.1:8080"
	hostDiscoveryPort := 8081

	hostNode := &discoverygo.Peer{
		Info: &peer.Info{
			Address: hostAddress,
			IsSelf:  true,
		},
		DiscoveryPort: hostDiscoveryPort,
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	// Starts an instance of groupcache with the provided transport
	daemon, err := groupcache.ListenAndServe(ctx, hostAddress, groupcache.Options{
		// If transport is nil, defaults to HttpTransport
		Transport: nil,
		Logger:    slog.Default(),
		Replicas:  50,
	})
	cancel()
	if err != nil {
		log.Fatal("while starting server on 192.168.1.1:8080")
	}

	// create an instance of the discovery engine on the node
	discoveryEngine := discoverygo.NewEngine(provider, daemon, hostNode)
	// start the discovery engine
	if err := discoveryEngine.Start(ctx); err != nil {
		log.Fatal(err)
	}

	// Create a new group cache with a max cache size of 3MB
	group, err := daemon.NewGroup("users", 3000000, groupcache.GetterFunc(
		func(ctx context.Context, id string, dest transport.Sink) error {
			// Set the user in the groupcache to expire after 5 minutes
			if err := dest.SetString("hello", time.Now().Add(time.Minute*5)); err != nil {
				return err
			}
			return nil
		},
	))
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel = context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	var value string
	if err := group.Get(ctx, "12345", transport.StringSink(&value)); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Value: %s\n", value)

	// Remove the key from the groupcache
	if err := group.Remove(ctx, "12345"); err != nil {
		fmt.Printf("Remove Err: %s\n", err)
		log.Fatal(err)
	}

	// Shutdown the discovery engine
	_ = discoveryEngine.Stop(ctx)
	// Shutdown the daemon
	_ = daemon.Shutdown(context.Background())
}
