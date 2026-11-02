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

const (
	// BlsPubKeyLength is the length of a compressed BLS12-381 public key in bytes.
	BlsPubKeyLength = 48
	// BlsSigLength is the length of a compressed BLS12-381 signature in bytes.
	BlsSigLength = 96
)

var (
	blsPubKeyType = reflect.TypeOf(BlsPubKey{})
	blsSigType    = reflect.TypeOf(BlsSignature{})
)

// BlsPubKey is a compressed BLS12-381 public key.
type BlsPubKey [BlsPubKeyLength]byte

// BytesToBlsPubKey returns b as a BlsPubKey, right-aligned. If b is longer
// than BlsPubKeyLength, its leading bytes are dropped.
func BytesToBlsPubKey(b []byte) BlsPubKey {
	var k BlsPubKey
	k.SetBytes(b)
	return k
}

// BigToBlsPubKey returns the big-endian bytes of i as a BlsPubKey.
func BigToBlsPubKey(i *big.Int) BlsPubKey { return BytesToBlsPubKey(i.Bytes()) }

// HexToBlsPubKey parses s (with or without a 0x prefix) as a BlsPubKey.
func HexToBlsPubKey(s string) BlsPubKey { return BytesToBlsPubKey(FromHex(s)) }

// Bytes returns the key as a byte slice.
func (k BlsPubKey) Bytes() []byte { return k[:] }

// Big returns the key as a big-endian unsigned integer.
func (k BlsPubKey) Big() *big.Int { return new(big.Int).SetBytes(k[:]) }

// Hex returns the 0x-prefixed hex encoding of the key.
func (k BlsPubKey) Hex() string { return hexutil.Encode(k[:]) }

// String implements fmt.Stringer.
func (k BlsPubKey) String() string { return k.Hex() }

// TerminalString returns an abbreviated form for log output.
func (k BlsPubKey) TerminalString() string { return terminalString(k[:]) }

// Format implements fmt.Formatter.
func (k BlsPubKey) Format(s fmt.State, verb rune) { formatFixed(s, verb, k[:], k.Hex(), "blspubkey") }

// SetBytes sets the key to b, right-aligned. If b is longer than
// BlsPubKeyLength, its leading bytes are dropped.
func (k *BlsPubKey) SetBytes(b []byte) { setRightAligned(k[:], b) }

// MarshalText implements encoding.TextMarshaler.
func (k BlsPubKey) MarshalText() ([]byte, error) { return hexutil.Bytes(k[:]).MarshalText() }

// UnmarshalText implements encoding.TextUnmarshaler.
func (k *BlsPubKey) UnmarshalText(input []byte) error {
	return hexutil.UnmarshalFixedText("BlsPubKey", input, k[:])
}

// UnmarshalJSON implements json.Unmarshaler.
func (k *BlsPubKey) UnmarshalJSON(input []byte) error {
	return hexutil.UnmarshalFixedJSON(blsPubKeyType, input, k[:])
}

// Scan implements sql.Scanner.
func (k *BlsPubKey) Scan(src interface{}) error { return scanFixed(k[:], src, "BlsPubKey") }

// Value implements driver.Valuer.
func (k BlsPubKey) Value() (driver.Value, error) { return k[:], nil }

// BlsSignature is a compressed BLS12-381 signature.
type BlsSignature [BlsSigLength]byte

// BytesToBlsSig returns b as a BlsSignature, right-aligned. If b is longer
// than BlsSigLength, its leading bytes are dropped.
func BytesToBlsSig(b []byte) BlsSignature {
	var s BlsSignature
	s.SetBytes(b)
	return s
}

// BigToBlsSig returns the big-endian bytes of i as a BlsSignature.
func BigToBlsSig(i *big.Int) BlsSignature { return BytesToBlsSig(i.Bytes()) }

// HexToBlsSig parses s (with or without a 0x prefix) as a BlsSignature.
func HexToBlsSig(s string) BlsSignature { return BytesToBlsSig(FromHex(s)) }

// Bytes returns the signature as a byte slice.
func (s BlsSignature) Bytes() []byte { return s[:] }

// Big returns the signature as a big-endian unsigned integer.
func (s BlsSignature) Big() *big.Int { return new(big.Int).SetBytes(s[:]) }

// Hex returns the 0x-prefixed hex encoding of the signature.
func (s BlsSignature) Hex() string { return hexutil.Encode(s[:]) }

// String implements fmt.Stringer.
func (s BlsSignature) String() string { return s.Hex() }

// TerminalString returns an abbreviated form for log output.
func (s BlsSignature) TerminalString() string { return terminalString(s[:]) }

// Format implements fmt.Formatter.
func (s BlsSignature) Format(st fmt.State, verb rune) {
	formatFixed(st, verb, s[:], s.Hex(), "blssignature")
}

// SetBytes sets the signature to b, right-aligned. If b is longer than
// BlsSigLength, its leading bytes are dropped.
func (s *BlsSignature) SetBytes(b []byte) { setRightAligned(s[:], b) }

// MarshalText implements encoding.TextMarshaler.
func (s BlsSignature) MarshalText() ([]byte, error) { return hexutil.Bytes(s[:]).MarshalText() }

// UnmarshalText implements encoding.TextUnmarshaler.
func (s *BlsSignature) UnmarshalText(input []byte) error {
	return hexutil.UnmarshalFixedText("BlsSignature", input, s[:])
}

// UnmarshalJSON implements json.Unmarshaler.
func (s *BlsSignature) UnmarshalJSON(input []byte) error {
	return hexutil.UnmarshalFixedJSON(blsSigType, input, s[:])
}

// Scan implements sql.Scanner.
func (s *BlsSignature) Scan(src interface{}) error { return scanFixed(s[:], src, "BlsSignature") }

// Value implements driver.Valuer.
func (s BlsSignature) Value() (driver.Value, error) { return s[:], nil }
