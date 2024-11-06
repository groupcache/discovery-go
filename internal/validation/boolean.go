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

import "errors"

// booleanValidator implements Validator.
type booleanValidator struct {
	boolCheck  bool
	errMessage string
}

// NewBooleanValidator creates a new boolean validator that returns an error message if condition is false
// This validator will come handy when dealing with conditional validation
func NewBooleanValidator(boolCheck bool, errMessage string) Validator {
	return &booleanValidator{boolCheck: boolCheck, errMessage: errMessage}
}

// Validate returns an error if boolean check is false
func (v booleanValidator) Validate() error {
	if !v.boolCheck {
		return errors.New(v.errMessage)
	}
	return nil
}
