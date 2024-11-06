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
	"strings"
	"time"

	"github.com/groupcache/discovery-go/internal/validation"
)

// Config represents the nats provider discoConfig
type Config struct {
	// Server defines the nats server in the format nats://host:port
	Server string
	// Subject defines the custom NATS subject
	Subject string
	// Timeout defines the nodes discovery timeout
	Timeout time.Duration
	// MaxJoinAttempts denotes the maximum number of attempts to connect an existing NATs server
	// Default to 5
	MaxJoinAttempts int
	// ReconnectWait sets the time to backoff after attempting a reconnect
	// to a server that we were already connected to previously.
	// Defaults to 2s.
	ReconnectWait time.Duration
	// Host specifies the host
	Host string
	// DiscoveryPort specifies the node discovery port
	DiscoveryPort int
}

// Validate checks whether the given discovery configuration is valid
func (config Config) Validate() error {
	return validation.New(validation.FailFast()).
		AddValidator(validation.NewEmptyStringValidator("Server", config.Server)).
		AddValidator(NewServerAddrValidator(config.Server)).
		AddValidator(validation.NewEmptyStringValidator("Subject", config.Subject)).
		AddValidator(validation.NewEmptyStringValidator("Host", config.Host)).
		AddAssertion(config.DiscoveryPort > 0, "DiscoveryPort is invalid").
		Validate()
}

// ServerAddrValidator helps validates the NATs server address
type ServerAddrValidator struct {
	server string
}

// NewServerAddrValidator validates the nats server address
func NewServerAddrValidator(server string) validation.Validator {
	return &ServerAddrValidator{server: server}
}

// Validate execute the validation code
func (x *ServerAddrValidator) Validate() error {
	// make sure that the nats prefix is set in the server address
	if !strings.HasPrefix(x.server, "nats") {
		return fmt.Errorf("invalid nats server address: %s", x.server)
	}

	hostAndPort := strings.SplitN(x.server, "nats://", 2)[1]
	return validation.NewTCPAddressValidator(hostAndPort).Validate()
}

// enforce compilation error
var _ validation.Validator = (*ServerAddrValidator)(nil)
