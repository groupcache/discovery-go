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

package dnssd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfig(t *testing.T) {
	t.Run("With valid configuration", func(t *testing.T) {
		config := &Config{
			DomainName: "google.com",
		}
		assert.NoError(t, config.Validate())
	})
	t.Run("With invalid configuration", func(t *testing.T) {
		config := &Config{
			DomainName: "",
		}
		assert.Error(t, config.Validate())
	})
}
