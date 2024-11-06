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
	"errors"
	"fmt"
	"net"
	"strconv"
	"strings"
)

var errFmt = "invalid address=(%s): %w"

// TCPAddressValidator helps validate a TCP address
type TCPAddressValidator struct {
	address string
}

// making sure the given struct implements the given interface
var _ Validator = (*TCPAddressValidator)(nil)

// NewTCPAddressValidator creates an instance of TCPAddressValidator
func NewTCPAddressValidator(address string) *TCPAddressValidator {
	return &TCPAddressValidator{address: address}
}

// Validate implements validation.Validator.
func (a *TCPAddressValidator) Validate() error {
	host, port, err := net.SplitHostPort(strings.TrimSpace(a.address))
	if err != nil {
		return fmt.Errorf(errFmt, a.address, err)
	}

	// let us validate the port number
	portNum, err := strconv.Atoi(port)
	if err != nil {
		return fmt.Errorf(errFmt, a.address, err)
	}

	// TODO: maybe we only need to check port number not to be negative
	if host == "" || portNum > 65535 || portNum < 0 {
		return fmt.Errorf(errFmt, a.address, errors.New("invalid address"))
	}

	return nil
}
