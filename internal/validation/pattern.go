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
	"errors"
	"regexp"
)

// patternValidator is used to perform a validation
// provided a given pattern
type patternValidator struct {
	pattern    string
	expression string
	customErr  error
}

var _ Validator = (*patternValidator)(nil)

// NewPatternValidator creates an instance of the validator
// The given pattern should be valid regular expression
func NewPatternValidator(pattern, expression string, customErr error) Validator {
	return &patternValidator{
		pattern:    pattern,
		expression: expression,
		customErr:  customErr,
	}
}

// Validate executes the validation
func (x *patternValidator) Validate() error {
	if match, _ := regexp.MatchString(x.pattern, x.expression); !match {
		if x.customErr != nil {
			return x.customErr
		}
		return errors.New("invalid expression")
	}
	return nil
}
