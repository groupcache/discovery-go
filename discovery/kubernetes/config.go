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

package kubernetes

import "github.com/groupcache/discovery-go/internal/validation"

// Config defines the kubernetes discovery configuration
type Config struct {
	// Namespace specifies the kubernetes namespace
	Namespace string
	// DiscoveryPortName specifies the discovery port name
	DiscoveryPortName string
	// PortName specifies the client port name
	PortName string
	// PodLabels specifies the pod labels
	PodLabels map[string]string
}

// Validate checks whether the given discovery configuration is valid
func (x Config) Validate() error {
	return validation.New(validation.FailFast()).
		AddValidator(validation.NewEmptyStringValidator("Namespace", x.Namespace)).
		AddValidator(validation.NewEmptyStringValidator("DiscoveryPortName", x.DiscoveryPortName)).
		AddValidator(validation.NewEmptyStringValidator("PortName", x.PortName)).
		AddAssertion(len(x.PodLabels) > 0, "PodLabels are required").
		Validate()
}
