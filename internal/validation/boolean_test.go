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

type booleanTestSuite struct {
	suite.Suite
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestBooleanValidator(t *testing.T) {
	suite.Run(t, new(booleanTestSuite))
}

func (s *booleanTestSuite) TestBooleanValidator() {
	s.Run("happy path when condition is true", func() {
		err := NewBooleanValidator(true, "error message").Validate()
		s.Assert().NoError(err)
	})
	s.Run("happy path when condition is false", func() {
		errMsg := "error message"
		err := NewBooleanValidator(false, errMsg).Validate()
		s.Assert().Error(err)
		s.Assert().EqualError(err, errMsg)
	})
}
