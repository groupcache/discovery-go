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

package discovery

import "errors"

var (
	// ErrAlreadyInitialized is used when attempting to re-initialize the discovery provider
	ErrAlreadyInitialized = errors.New("provider already initialized")
	// ErrNotInitialized is used when the provider is not initialized
	ErrNotInitialized = errors.New("provider not initialized")
	// ErrAlreadyRegistered is used when attempting to re-register the provider
	ErrAlreadyRegistered = errors.New("provider already registered")
	// ErrNotRegistered is used when attempting to de-register the provider
	ErrNotRegistered = errors.New("provider is not registered")
)
