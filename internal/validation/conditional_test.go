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
	"testing"

	"github.com/stretchr/testify/suite"
)

type conditionalTestSuite struct {
	suite.Suite
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestConditionalValidator(t *testing.T) {
	suite.Run(t, new(conditionalTestSuite))
}

func (s *conditionalTestSuite) TestConditionalValidator() {
	s.Run("with condition set to true", func() {
		fieldName := "field"
		fieldValue := ""
		validator := NewConditionalValidator(true, NewEmptyStringValidator(fieldName, fieldValue))
		err := validator.Validate()
		s.Assert().Error(err)
		s.Assert().EqualError(err, "the [field] is required")
	})
	s.Run("with condition set to false", func() {
		fieldName := "field"
		fieldValue := ""
		validator := NewConditionalValidator(false, NewEmptyStringValidator(fieldName, fieldValue))
		err := validator.Validate()
		s.Assert().NoError(err)
	})
}
