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

// Package hexutil implements the "0x"-prefixed hexadecimal text encoding used
// by JSON-RPC APIs.
//
// Two kinds of values are supported:
//
//   - byte strings: "0x" followed by an even number of hex digits ("0x" is the
//     empty string);
//   - quantities (unsigned integers): "0x" followed by the shortest hex digits
//     representing the number, without leading zeros ("0x0" is zero).
//
// Decoding accepts both lower and upper case digits and an "0X" prefix.
package hexutil

import (
	"encoding/hex"
	"fmt"
	"math/big"
	"strconv"
)

// Decoding errors.
var (
	ErrEmptyString   = hexError("empty hex string")
	ErrSyntax        = hexError("invalid hex string")
	ErrMissingPrefix = hexError("hex string without 0x prefix")
	ErrOddLength     = hexError("hex string of odd length")
	ErrEmptyNumber   = hexError(`hex string "0x"`)
	ErrLeadingZero   = hexError("hex number with leading zero digits")
	ErrUint64Range   = hexError("hex number > 64 bits")
	ErrUint8Range    = hexError("hex number > 8 bits")
	ErrBig256Range   = hexError("hex number > 256 bits")
)

// hexError is the type of all decoding errors returned by this package.
type hexError string

func (e hexError) Error() string { return string(e) }

const prefix = "0x"

// hasPrefix reports whether s starts with "0x" or "0X".
func hasPrefix[T string | []byte](s T) bool {
	return len(s) >= 2 && s[0] == '0' && (s[1] == 'x' || s[1] == 'X')
}

// Encode returns b as a "0x"-prefixed hex string.
func Encode(b []byte) string {
	out := make([]byte, len(prefix)+hex.EncodedLen(len(b)))
	copy(out, prefix)
	hex.Encode(out[len(prefix):], b)
	return string(out)
}

// Decode decodes a "0x"-prefixed hex byte string.
func Decode(s string) ([]byte, error) {
	if s == "" {
		return nil, ErrEmptyString
	}
	if !hasPrefix(s) {
		return nil, ErrMissingPrefix
	}
	out, err := hex.DecodeString(s[len(prefix):])
	if err != nil {
		return nil, digitsError(err)
	}
	return out, nil
}

// MustDecode is like Decode but panics on error.
func MustDecode(s string) []byte {
	b, err := Decode(s)
	if err != nil {
		panic(err)
	}
	return b
}

// EncodeUint64 returns i as a hex quantity.
func EncodeUint64(i uint64) string {
	return prefix + strconv.FormatUint(i, 16)
}

// DecodeUint64 decodes a hex quantity that fits into 64 bits.
func DecodeUint64(s string) (uint64, error) {
	digits, err := quantityDigits([]byte(s), true)
	if err != nil {
		return 0, err
	}
	return parseUint(digits, 64, ErrUint64Range, true)
}

// EncodeBig returns i as a hex quantity. Negative numbers get a "-" sign in
// front of the prefix.
func EncodeBig(i *big.Int) string {
	if i.Sign() < 0 {
		return "-" + prefix + new(big.Int).Abs(i).Text(16)
	}
	return prefix + i.Text(16)
}

// DecodeBig decodes a hex quantity of at most 256 bits.
func DecodeBig(s string) (*big.Int, error) {
	digits, err := quantityDigits([]byte(s), true)
	if err != nil {
		return nil, err
	}
	return parseBig(digits)
}

// decodeDigits decodes an even number of hex digits.
func decodeDigits(digits []byte) ([]byte, error) {
	if len(digits)%2 != 0 {
		return nil, ErrOddLength
	}
	out := make([]byte, len(digits)/2)
	if _, err := hex.Decode(out, digits); err != nil {
		return nil, ErrSyntax
	}
	return out, nil
}

// byteStringDigits validates the text form of a byte string and returns its
// hex digits. Empty input is accepted and yields no digits. If requirePrefix
// is false, the "0x" prefix is optional.
func byteStringDigits(text []byte, requirePrefix bool) ([]byte, error) {
	if len(text) == 0 {
		return nil, nil
	}
	switch {
	case hasPrefix(text):
		text = text[len(prefix):]
	case requirePrefix:
		return nil, ErrMissingPrefix
	}
	if len(text)%2 != 0 {
		return nil, ErrOddLength
	}
	return text, nil
}

// quantityDigits validates the text form of a quantity and returns its hex
// digits. Empty input yields no digits (the zero value) unless strict is set,
// in which case it is an error.
func quantityDigits(text []byte, strict bool) ([]byte, error) {
	if len(text) == 0 {
		if strict {
			return nil, ErrEmptyString
		}
		return nil, nil
	}
	if !hasPrefix(text) {
		return nil, ErrMissingPrefix
	}
	digits := text[len(prefix):]
	if len(digits) == 0 {
		return nil, ErrEmptyNumber
	}
	if len(digits) > 1 && digits[0] == '0' {
		return nil, ErrLeadingZero
	}
	return digits, nil
}

// parseUint parses hex digits into an unsigned integer of the given bit size.
// If sequential is set, digits are checked one by one and the first problem
// found (an invalid digit or too many digits) is reported; otherwise too many
// digits are reported before any invalid digit.
func parseUint(digits []byte, bits int, rangeErr error, sequential bool) (uint64, error) {
	maxDigits := bits / 4
	if !sequential && len(digits) > maxDigits {
		return 0, rangeErr
	}
	var v uint64
	for i, c := range digits {
		n, ok := nibble(c)
		if !ok {
			return 0, ErrSyntax
		}
		if i >= maxDigits {
			return 0, rangeErr
		}
		v = v<<4 | uint64(n)
	}
	return v, nil
}

// parseBig parses hex digits into a big integer of at most 256 bits. No
// digits yield zero.
func parseBig(digits []byte) (*big.Int, error) {
	if len(digits) == 0 {
		return new(big.Int), nil
	}
	if len(digits) > 64 {
		return nil, ErrBig256Range
	}
	for _, c := range digits {
		if _, ok := nibble(c); !ok {
			return nil, ErrSyntax
		}
	}
	v, ok := new(big.Int).SetString(string(digits), 16)
	if !ok {
		return nil, ErrSyntax
	}
	return v, nil
}

// nibble returns the value of a single hex digit.
func nibble(c byte) (byte, bool) {
	switch {
	case '0' <= c && c <= '9':
		return c - '0', true
	case 'a' <= c && c <= 'f':
		return c - 'a' + 10, true
	case 'A' <= c && c <= 'F':
		return c - 'A' + 10, true
	}
	return 0, false
}

// digitsError maps errors of encoding/hex to the errors of this package. An
// invalid digit is reported before an odd length.
func digitsError(err error) error {
	if err == hex.ErrLength {
		return ErrOddLength
	}
	return ErrSyntax
}

// lengthError reports a fixed-size byte string of the wrong length.
func lengthError(gotDigits, wantBytes int, typeName string) error {
	return fmt.Errorf("hex string has length %d, want %d for %s", gotDigits, wantBytes*2, typeName)
}
