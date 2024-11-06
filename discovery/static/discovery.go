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

import "github.com/tochemey/gokv/discovery"

// Discovery represents the static discovery provider
type Discovery struct {
	config *Config
}

// enforce compilation error
var _ discovery.Provider = &Discovery{}

// NewDiscovery creates an instance of the static discovery provider
func NewDiscovery(config *Config) *Discovery {
	d := &Discovery{
		config: config,
	}

	return d
}

// ID returns the discovery provider identifier
func (d *Discovery) ID() string {
	return "static"
}

// Initialize the discovery provider
func (d *Discovery) Initialize() error {
	return d.config.Validate()
}

// Register registers this node to a service discovery directory.
func (d *Discovery) Register() error {
	return nil
}

// Deregister removes this node from a service discovery directory.
func (d *Discovery) Deregister() error {
	return nil
}

// Close closes the provider
func (d *Discovery) Close() error {
	return nil
}

// DiscoverPeers returns a list of known nodes.
func (d *Discovery) DiscoverPeers() ([]string, error) {
	return d.config.Hosts, nil
}
