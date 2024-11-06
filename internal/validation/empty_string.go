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

import "fmt"

type emptyStringValidator struct {
	fieldValue string
	fieldName  string
}

// NewEmptyStringValidator creates a string a emptyString validator
func NewEmptyStringValidator(fieldName, fieldValue string) Validator {
	return emptyStringValidator{fieldValue: fieldValue, fieldName: fieldName}
}

// Validate checks whether the given string is empty or not
func (v emptyStringValidator) Validate() error {
	if v.fieldValue == "" {
		return fmt.Errorf("the [%s] is required", v.fieldName)
	}
	return nil
}
