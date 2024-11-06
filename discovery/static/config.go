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

import "github.com/groupcache/discovery-go/internal/validation"

// Config represents the static discovery provider configuration
type Config struct {
	// Hosts defines the list of hosts in the form of ip:port where the port is the  gossip port.
	Hosts []string
}

// Validate checks whether the given discovery configuration is valid
func (x Config) Validate() error {
	chain := validation.
		New(validation.FailFast()).
		AddAssertion(len(x.Hosts) != 0, "hosts are required")

	for _, host := range x.Hosts {
		chain = chain.AddValidator(validation.NewTCPAddressValidator(host))
	}

	return chain.Validate()
}
