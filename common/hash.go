// Copyright 2026 Digital Clever Solution LLC
// SPDX-License-Identifier: Apache-2.0
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package common

import (
	"database/sql/driver"
	"fmt"
	"math/big"
	"reflect"

	"github.com/LFDT-Iguazu/iguazu-types/common/hexutil"
)

// HashLength is the length of a hash in bytes.
const HashLength = 32

var hashType = reflect.TypeOf(Hash{})

// Hash is a 32-byte hash value, for example a Keccak-256 digest.
type Hash [HashLength]byte

// BytesToHash returns b as a Hash, right-aligned. If b is longer than
// HashLength, its leading bytes are dropped.
func BytesToHash(b []byte) Hash {
	var h Hash
	h.SetBytes(b)
	return h
}

// BigToHash returns the big-endian bytes of i as a Hash.
func BigToHash(i *big.Int) Hash { return BytesToHash(i.Bytes()) }

// HexToHash parses s (with or without a 0x prefix) as a Hash.
func HexToHash(s string) Hash { return BytesToHash(FromHex(s)) }

// Bytes returns the hash as a byte slice.
func (h Hash) Bytes() []byte { return h[:] }

// Big returns the hash as a big-endian unsigned integer.
func (h Hash) Big() *big.Int { return new(big.Int).SetBytes(h[:]) }

// Hex returns the 0x-prefixed hex encoding of the hash.
func (h Hash) Hex() string { return hexutil.Encode(h[:]) }

// String implements fmt.Stringer.
func (h Hash) String() string { return h.Hex() }

// TerminalString returns an abbreviated form for log output.
func (h Hash) TerminalString() string { return terminalString(h[:]) }

// Format implements fmt.Formatter.
func (h Hash) Format(s fmt.State, verb rune) { formatFixed(s, verb, h[:], h.Hex(), "hash") }

// SetBytes sets the hash to b, right-aligned. If b is longer than HashLength,
// its leading bytes are dropped.
func (h *Hash) SetBytes(b []byte) { setRightAligned(h[:], b) }

// MarshalText implements encoding.TextMarshaler.
func (h Hash) MarshalText() ([]byte, error) { return hexutil.Bytes(h[:]).MarshalText() }

// UnmarshalText implements encoding.TextUnmarshaler.
func (h *Hash) UnmarshalText(input []byte) error {
	return hexutil.UnmarshalFixedText("Hash", input, h[:])
}

// UnmarshalJSON implements json.Unmarshaler.
func (h *Hash) UnmarshalJSON(input []byte) error {
	return hexutil.UnmarshalFixedJSON(hashType, input, h[:])
}

// Scan implements sql.Scanner.
func (h *Hash) Scan(src interface{}) error { return scanFixed(h[:], src, "Hash") }

// Value implements driver.Valuer.
func (h Hash) Value() (driver.Value, error) { return h[:], nil }
