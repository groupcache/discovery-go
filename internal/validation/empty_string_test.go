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

type emptyStringTestSuite struct {
	suite.Suite
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestEmptyStringValidator(t *testing.T) {
	suite.Run(t, new(emptyStringTestSuite))
}

func (s *emptyStringTestSuite) TestEmptyStringValidator() {
	s.Run("with string value set", func() {
		validator := NewEmptyStringValidator("field", "some-value")
		s.Assert().NotNil(validator)
		err := validator.Validate()
		s.Assert().NoError(err)
	})
	s.Run("with string value not set", func() {
		validator := NewEmptyStringValidator("field", "")
		s.Assert().NotNil(validator)
		err := validator.Validate()
		s.Assert().Error(err)
		s.Assert().EqualError(err, "the [field] is required")
	})
}
