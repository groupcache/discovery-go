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

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestErrorsChain(t *testing.T) {
	t.Run("With ReturnFirst", func(t *testing.T) {
		e1 := errors.New("err1")
		e2 := errors.New("err2")
		e3 := errors.New("err3")

		chain := New(ReturnFirst()).AddError(e1).AddError(e2).AddError(e3)
		actual := chain.Error()
		assert.True(t, errors.Is(actual, e1))
	})
	t.Run("With ReturnAll", func(t *testing.T) {
		e1 := errors.New("err1")
		e2 := errors.New("err2")
		e3 := errors.New("err3")

		chain := New(ReturnAll()).AddError(e1).AddError(e2).AddError(e3)
		actual := chain.Error()
		assert.EqualError(t, actual, "err1; err2; err3")
	})
}
