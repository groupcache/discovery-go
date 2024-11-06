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

// conditionalValidator runs a validator when a condition is met
type conditionalValidator struct {
	c bool
	v Validator
}

// NewConditionalValidator creates a conditional validator, that runs the validator if the condition is true.
// This validator will help when performing data update
func NewConditionalValidator(condition bool, validator Validator) Validator {
	return &conditionalValidator{c: condition, v: validator}
}

// Validate runs the provided conditional validator
func (v conditionalValidator) Validate() error {
	if v.c {
		return v.v.Validate()
	}
	return nil
}
