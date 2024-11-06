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

import "github.com/groupcache/discovery-go/internal/validation"

// Config defines the discovery configuration
type Config struct {
	// Domain specifies the dns name
	DomainName string
	// IPv6 states whether to fetch ipv6 address instead of ipv4
	// if it is false then all addresses are extracted
	IPv6 *bool
}

// Validate checks whether the given discovery configuration is valid
func (x Config) Validate() error {
	return validation.New(validation.FailFast()).
		AddValidator(validation.NewEmptyStringValidator("Namespace", x.DomainName)).
		Validate()
}
