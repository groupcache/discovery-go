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

package discovery

// Provider helps discover other running nodes in a cloud environment
type Provider interface {
	// ID returns the service discovery provider name
	ID() string
	// Initialize initializes the service discovery provider.
	Initialize() error
	// Register registers the service discovery provider.
	Register() error
	// Deregister de-registers the service discovery provider.
	Deregister() error
	// DiscoverPeers returns a list discovered nodes' addresses.
	DiscoverPeers() ([]string, error)
	// Close closes the provider
	Close() error
}
