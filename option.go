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
	"time"

	"github.com/groupcache/discovery-go/logger"
)

// Option is the interface that applies to the Node
type Option interface {
	// Apply sets the Option value of a config.
	Apply(*Engine)
}

var _ Option = OptionFunc(nil)

// OptionFunc implements the GroupOtion interface.
type OptionFunc func(node *Engine)

// Apply applies the Node's option
func (f OptionFunc) Apply(node *Engine) {
	f(node)
}

// WithLogger sets the cacheLogger
func WithLogger(logger logger.Logger) Option {
	return OptionFunc(func(node *Engine) {
		node.logger = logger
	})
}

// WithShutdownTimeout sets the Node shutdown timeout.
func WithShutdownTimeout(timeout time.Duration) Option {
	return OptionFunc(func(node *Engine) {
		node.shutdownTimeout = timeout
	})
}

// WithMinimumPeersQuorum sets the minimum number of nodes to form a quorum
func WithMinimumPeersQuorum(minimumQuorum uint) Option {
	return OptionFunc(func(node *Engine) {
		node.minimumPeersQuorum = minimumQuorum
	})
}

// WithMaxJoinTimeout sets the max join timeout
func WithMaxJoinTimeout(timeout time.Duration) Option {
	return OptionFunc(func(node *Engine) {
		node.maxJoinTimeout = timeout
	})
}

// WithMaxJoinAttempts sets the max join attempts
func WithMaxJoinAttempts(attempts int) Option {
	return OptionFunc(func(node *Engine) {
		node.maxJoinAttempts = attempts
	})
}

// WithMaxJoinRetryInterval sets the max join retry interval
func WithMaxJoinRetryInterval(retryInterval time.Duration) Option {
	return OptionFunc(func(node *Engine) {
		node.maxJoinRetryInterval = retryInterval
	})
}
