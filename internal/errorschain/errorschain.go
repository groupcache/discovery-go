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

package errorschain

import "github.com/groupcache/groupcache-go/v3"

// Chain defines an error chain
type Chain struct {
	returnFirst bool
	errs        []error
}

// ChainOption configures a validation chain at creation time.
type ChainOption func(*Chain)

// New creates a new error chain. All errors will be evaluated respectively
// according to their insertion order
func New(opts ...ChainOption) *Chain {
	chain := &Chain{
		errs: make([]error, 0),
	}

	for _, opt := range opts {
		opt(chain)
	}

	return chain
}

// AddError add an error to the chain
func (c *Chain) AddError(err error) *Chain {
	c.errs = append(c.errs, err)
	return c
}

// Error returns the error
func (c *Chain) Error() error {
	var merr = &groupcache.MultiError{}
	for _, v := range c.errs {
		if v != nil {
			if c.returnFirst {
				// just return the error
				return v
			}
			// append error to the violations
			merr.Add(v)
		}
	}
	return merr.NilOrError()
}

// ReturnFirst sets whether a chain should stop validation on first error.
func ReturnFirst() ChainOption {
	return func(c *Chain) { c.returnFirst = true }
}

// ReturnAll sets whether a chain should return all errors.
func ReturnAll() ChainOption {
	return func(c *Chain) { c.returnFirst = false }
}
