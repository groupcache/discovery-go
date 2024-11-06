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

type validationTestSuite struct {
	suite.Suite
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestValidation(t *testing.T) {
	suite.Run(t, new(validationTestSuite))
}

func (s *validationTestSuite) TestNewChain() {
	s.Run("new chain without option", func() {
		chain := New()
		s.Assert().NotNil(chain)
	})
	s.Run("new chain with options", func() {
		chain := New(FailFast())
		s.Assert().NotNil(chain)
		s.Assert().True(chain.failFast)
		chain2 := New(AllErrors())
		s.Assert().NotNil(chain2)
		s.Assert().False(chain2.failFast)
	})
}

func (s *validationTestSuite) TestAddValidator() {
	chain := New()
	s.Assert().NotNil(chain)
	s.Assert().Empty(chain.validators)
	chain.AddValidator(NewBooleanValidator(true, ""))
	s.Assert().NotEmpty(chain.validators)
	s.Assert().Equal(1, len(chain.validators))
}

func (s *validationTestSuite) TestAddAssertion() {
	chain := New()
	s.Assert().NotNil(chain)
	s.Assert().Empty(chain.validators)
	chain.AddAssertion(true, "")
	s.Assert().NotEmpty(chain.validators)
	s.Assert().Equal(1, len(chain.validators))
}

func (s *validationTestSuite) TestValidate() {
	s.Run("with single validator", func() {
		chain := New()
		s.Assert().NotNil(chain)
		chain.AddValidator(NewEmptyStringValidator("field", ""))
		s.Assert().Nil(chain.violations)
		err := chain.Validate()
		s.Assert().NotNil(chain.violations)
		s.Assert().Error(err)
		s.Assert().EqualError(err, "the [field] is required")
	})
	s.Run("with multiple validators and FailFast option", func() {
		chain := New(FailFast())
		s.Assert().NotNil(chain)
		chain.
			AddValidator(NewEmptyStringValidator("field", "")).
			AddAssertion(false, "this is false")
		s.Assert().Nil(chain.violations)
		err := chain.Validate()
		s.Assert().Nil(chain.violations)
		s.Assert().Error(err)
		s.Assert().EqualError(err, "the [field] is required")
	})
	s.Run("with multiple validators and AllErrors option", func() {
		chain := New(AllErrors())
		s.Assert().NotNil(chain)
		chain.
			AddValidator(NewEmptyStringValidator("field", "")).
			AddAssertion(false, "this is false")
		s.Assert().Nil(chain.violations)
		err := chain.Validate()
		s.Assert().NotNil(chain.violations)
		s.Assert().Error(err)
		s.Assert().EqualError(err, "the [field] is required; this is false")
	})
}
