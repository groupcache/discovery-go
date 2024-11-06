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

package discoverygo

import (
	"encoding/json"

	"github.com/hashicorp/memberlist"
)

// delegate implements memberlist Delegate
type delegate struct {
	meta []byte
}

// enforce compilation error
var _ memberlist.Delegate = (*delegate)(nil)

// newEngineDelegate creates an instance of delegate instance
func (g *Engine) newEngineDelegate() (*delegate, error) {
	info := g.hostNode
	bytea, err := json.Marshal(info)
	if err != nil {
		return nil, err
	}
	return &delegate{bytea}, nil
}

func (deletgate *delegate) NodeMeta(limit int) []byte {
	return deletgate.meta
}

func (deletgate *delegate) NotifyMsg(bytes []byte) {
}

func (deletgate *delegate) GetBroadcasts(overhead, limit int) [][]byte {
	return nil
}

func (deletgate *delegate) LocalState(join bool) []byte {
	return nil
}

func (deletgate *delegate) MergeRemoteState(buf []byte, join bool) {
}
