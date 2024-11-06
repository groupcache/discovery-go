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

package dnssd

import (
	"fmt"
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/atomic"

	"github.com/groupcache/discovery-go/discovery"
)

func TestDiscovery(t *testing.T) {
	t.Run("With a new instance", func(t *testing.T) {
		// create the instance of provider
		provider := NewDiscovery(nil)
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
		// cheap test
		// create the instance of provider
		provider := NewDiscovery(nil)
		require.NotNil(t, provider)
		assert.Equal(t, "dns-sd", provider.ID())
	})
	t.Run("With Initialize", func(t *testing.T) {
		// create the config
		config := &Config{
			DomainName: "google.com",
		}
		// create the instance of provider
		provider := NewDiscovery(config)
		assert.NoError(t, provider.Initialize())
	})
	t.Run("With Initialize: already initialized", func(t *testing.T) {
		// create the config
		config := &Config{
			DomainName: "google.com",
		}
		// create the instance of provider
		provider := NewDiscovery(config)
		provider.initialized = atomic.NewBool(true)
		assert.Error(t, provider.Initialize())
	})
	t.Run("With Register", func(t *testing.T) {
		// create the config
		config := &Config{
			DomainName: "google.com",
		}
		// create the instance of provider
		provider := NewDiscovery(config)
		require.NoError(t, provider.Initialize())
		require.NoError(t, provider.Register())

		time.Sleep(time.Second)
		require.True(t, provider.initialized.Load())
		require.NoError(t, provider.Deregister())
		time.Sleep(time.Second)
		assert.False(t, provider.initialized.Load())
	})
	t.Run("With Register when already registered", func(t *testing.T) {
		// create the config
		config := &Config{
			DomainName: "google.com",
		}
		// create the instance of provider
		provider := NewDiscovery(config)
		require.NoError(t, provider.Initialize())
		require.NoError(t, provider.Register())

		time.Sleep(time.Second)
		require.True(t, provider.initialized.Load())
		err := provider.Register()
		require.Error(t, err)
		require.EqualError(t, err, discovery.ErrAlreadyRegistered.Error())
		require.NoError(t, provider.Deregister())
		time.Sleep(time.Second)
		assert.False(t, provider.initialized.Load())
	})
	t.Run("With Deregister", func(t *testing.T) {
		// create the config
		config := &Config{
			DomainName: "google.com",
		}
		// create the instance of provider
		provider := NewDiscovery(config)
		// for the sake of the test
		provider.initialized = atomic.NewBool(true)
		assert.NoError(t, provider.Deregister())
	})
	t.Run("With Deregister when not initialized", func(t *testing.T) {
		// create the config
		config := &Config{
			DomainName: "google.com",
		}
		// create the instance of provider
		provider := NewDiscovery(config)
		// for the sake of the test
		provider.initialized = atomic.NewBool(false)
		err := provider.Deregister()
		assert.Error(t, err)
		assert.EqualError(t, err, discovery.ErrNotInitialized.Error())
	})
	t.Run("With DiscoverPeers: not initialized", func(t *testing.T) {
		// create the config
		config := &Config{
			DomainName: "google.com",
		}
		// create the instance of provider
		provider := NewDiscovery(config)
		peers, err := provider.DiscoverPeers()
		assert.Error(t, err)
		assert.Empty(t, peers)
		assert.EqualError(t, err, discovery.ErrNotInitialized.Error())
	})
	t.Run("With DiscoverPeers", func(t *testing.T) {
		// create the config
		config := &Config{
			DomainName: "google.com",
		}
		// create the instance of provider
		provider := NewDiscovery(config)

		require.NoError(t, provider.Initialize())
		require.NoError(t, provider.Register())

		// wait for registration to be completed
		time.Sleep(time.Second)
		require.True(t, provider.initialized.Load())

		// discover peers
		peers, err := provider.DiscoverPeers()
		require.NoError(t, err)
		require.NotEmpty(t, peers)
		require.NoError(t, validateAddrs(peers))
		assert.NoError(t, provider.Deregister())
		assert.NoError(t, provider.Close())
	})
	t.Run("With DiscoverPeers with IPV6", func(t *testing.T) {
		// create the config
		ipv6 := true
		config := &Config{
			DomainName: "google.com",
			IPv6:       &ipv6,
		}
		// create the instance of provider
		provider := NewDiscovery(config)

		require.NoError(t, provider.Initialize())
		require.NoError(t, provider.Register())

		// wait for registration to be completed
		time.Sleep(time.Second)
		require.True(t, provider.initialized.Load())

		// discover peers
		peers, err := provider.DiscoverPeers()
		require.NoError(t, err)
		require.NotEmpty(t, peers)
		require.NoError(t, validateAddrs(peers))
		assert.NoError(t, provider.Deregister())
		assert.NoError(t, provider.Close())
	})
}

func validateAddrs(addrs []string) error {
	for _, addr := range addrs {
		ipaddr := net.ParseIP(addr)
		if ipaddr == nil {
			return fmt.Errorf("invalid IP address format: %s", addr)
		}
	}
	return nil
}