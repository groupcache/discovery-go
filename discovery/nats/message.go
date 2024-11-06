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

package nats

// MessageType defines the nats message type
type MessageType int

const (
	Register MessageType = iota
	Deregister
	Request
	Response
)

// Message is used internally by the nats discovery provider
// to communicate with nodes in the cluster
type Message struct {
	// Host defines the node host address
	Host string `json:"host"`
	// Port defines the node port
	Port int `json:"port"`
	// Name defines the node name
	Name string `json:"name"`
	// Type defines the type of message sent by a given node
	Type MessageType `json:"type"`
}
