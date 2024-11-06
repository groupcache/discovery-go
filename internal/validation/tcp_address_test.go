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

	"github.com/stretchr/testify/assert"
)

func TestTCPAddressValidator(t *testing.T) {
	t.Run("With happy path", func(t *testing.T) {
		addr := "127.0.0.1:3222"
		assert.NoError(t, NewTCPAddressValidator(addr).Validate())
	})
	t.Run("With invalid port number: case 1", func(t *testing.T) {
		addr := "127.0.0.1:-1"
		assert.Error(t, NewTCPAddressValidator(addr).Validate())
	})
	t.Run("With invalid port number: case 2", func(t *testing.T) {
		addr := "127.0.0.1:655387"
		assert.Error(t, NewTCPAddressValidator(addr).Validate())
	})
	t.Run("With  zero port number: case 3", func(t *testing.T) {
		addr := "127.0.0.1:0"
		assert.NoError(t, NewTCPAddressValidator(addr).Validate())
	})
	t.Run("With invalid host", func(t *testing.T) {
		addr := ":3222"
		assert.Error(t, NewTCPAddressValidator(addr).Validate())
	})
}
