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
	"encoding/hex"
	"fmt"
	"math/big"
	"reflect"

	"github.com/LFDT-Iguazu/iguazu-types/common/hexutil"
	"golang.org/x/crypto/sha3"
)

// AddressLength is the length of an account address in bytes.
const AddressLength = 20

var addressType = reflect.TypeOf(Address{})

// Address is a 20-byte account address.
type Address [AddressLength]byte

// BytesToAddress returns b as an Address, right-aligned. If b is longer than
// AddressLength, its leading bytes are dropped.
func BytesToAddress(b []byte) Address {
	var a Address
	a.SetBytes(b)
	return a
}

// BigToAddress returns the big-endian bytes of i as an Address.
func BigToAddress(i *big.Int) Address { return BytesToAddress(i.Bytes()) }

// HexToAddress parses s (with or without a 0x prefix) as an Address.
func HexToAddress(s string) Address { return BytesToAddress(FromHex(s)) }

// IsHexAddress reports whether s is a hex-encoded address, with or without a
// 0x prefix, in any letter case.
func IsHexAddress(s string) bool {
	if has0xPrefix(s) {
		s = s[2:]
	}
	return len(s) == 2*AddressLength && isHex(s)
}

// Bytes returns the address as a byte slice.
func (a Address) Bytes() []byte { return a[:] }

// Hash returns the address right-aligned in a Hash.
func (a Address) Hash() Hash { return BytesToHash(a[:]) }

// Hex returns the address in the mixed-case checksum encoding of EIP-55.
func (a Address) Hex() string { return string(a.checksumHex()) }

// String implements fmt.Stringer.
func (a Address) String() string { return a.Hex() }

// checksumHex returns the 0x-prefixed EIP-55 encoding: a hex letter is upper
// case when the matching nibble of keccak256(lowercase hex) is 8 or more.
func (a Address) checksumHex() []byte {
	digits := make([]byte, 2*AddressLength)
	hex.Encode(digits, a[:])

	hasher := sha3.NewLegacyKeccak256()
	hasher.Write(digits)
	digest := hasher.Sum(nil)

	for i, c := range digits {
		if c < 'a' {
			continue // a decimal digit
		}
		n := digest[i/2]
		if i%2 == 0 {
			n >>= 4
		}
		if n&0x0f >= 8 {
			digits[i] = c - 'a' + 'A'
		}
	}
	return append([]byte("0x"), digits...)
}

// Format implements fmt.Formatter.
func (a Address) Format(s fmt.State, verb rune) { formatFixed(s, verb, a[:], a.Hex(), "address") }

// SetBytes sets the address to b, right-aligned. If b is longer than
// AddressLength, its leading bytes are dropped.
func (a *Address) SetBytes(b []byte) { setRightAligned(a[:], b) }

// MarshalText implements encoding.TextMarshaler. The output is lower case.
func (a Address) MarshalText() ([]byte, error) { return hexutil.Bytes(a[:]).MarshalText() }

// UnmarshalText implements encoding.TextUnmarshaler.
func (a *Address) UnmarshalText(input []byte) error {
	return hexutil.UnmarshalFixedText("Address", input, a[:])
}

// UnmarshalJSON implements json.Unmarshaler.
func (a *Address) UnmarshalJSON(input []byte) error {
	return hexutil.UnmarshalFixedJSON(addressType, input, a[:])
}

// Scan implements sql.Scanner.
func (a *Address) Scan(src interface{}) error { return scanFixed(a[:], src, "Address") }

// Value implements driver.Valuer.
func (a Address) Value() (driver.Value, error) { return a[:], nil }
