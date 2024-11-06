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

package validation

import (
	"go.uber.org/multierr"
)

// Validator interface generalizes the validator implementations
type Validator interface {
	Validate() error
}

// Chain represents list of validators and is used to accumulate errors and return them as a single "error"
type Chain struct {
	failFast   bool
	validators []Validator
	violations error
}

// ChainOption configures a validation chain at creation time.
type ChainOption func(*Chain)

// New creates a new validation chain.
func New(opts ...ChainOption) *Chain {
	chain := &Chain{
		validators: make([]Validator, 0),
	}

	for _, opt := range opts {
		opt(chain)
	}

	return chain
}

// FailFast sets whether a chain should stop validation on first error.
func FailFast() ChainOption {
	return func(c *Chain) { c.failFast = true }
}

// AllErrors sets whether a chain should return all errors.
func AllErrors() ChainOption {
	return func(c *Chain) { c.failFast = false }
}

// AddValidator adds validator to the validation chain.
func (c *Chain) AddValidator(v Validator) *Chain {
	c.validators = append(c.validators, v)
	return c
}

// AddAssertion adds assertion to the validation chain.
func (c *Chain) AddAssertion(isTrue bool, message string) *Chain {
	c.validators = append(c.validators, NewBooleanValidator(isTrue, message))
	return c
}

// Validate runs validation chain and returns resulting error(s).
// It returns all validation error by default, use FailFast option to stop validation on first error.
func (c *Chain) Validate() error {
	for _, v := range c.validators {
		if violations := v.Validate(); violations != nil {
			if c.failFast {
				// just return the error
				return violations
			}
			// append error to the violations
			c.violations = multierr.Append(c.violations, violations)
		}
	}
	return c.violations
}
