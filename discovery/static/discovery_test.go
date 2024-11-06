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

package static

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/groupcache/discovery-go/discovery"
)

func TestDiscovery(t *testing.T) {
	t.Run("With new instance", func(t *testing.T) {
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
		assert.Equal(t, "static", provider.ID())
	})

	t.Run("With DiscoverPeers", func(t *testing.T) {
		// create the config
		config := Config{
			Hosts: []string{
				"192.168.0.1:3000",
				"192.168.0.1:3001",
				"192.168.0.2:3000",
			},
		}

		// create the instance of provider
		provider := NewDiscovery(&config)
		require.NoError(t, provider.Initialize())
		require.NoError(t, provider.Register())

		// discover peers
		peers, err := provider.DiscoverPeers()
		require.NoError(t, err)
		require.Len(t, peers, 3)

		assert.NoError(t, provider.Deregister())
		assert.NoError(t, provider.Close())
	})
}
