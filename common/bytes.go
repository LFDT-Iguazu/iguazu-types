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

// Package common contains shared Iguazu types and utilities.
package common

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
)

const (
	Uint64Size = 8
	Uint32Size = 4
	Uint16Size = 2
	Uint8Size  = 1
)

// FromHex decodes the hex string s, with or without a 0x prefix. An odd
// number of digits is treated as if it had a leading zero. Decoding stops at
// the first invalid digit.
func FromHex(s string) []byte {
	if has0xPrefix(s) {
		s = s[2:]
	}
	if len(s)%2 == 1 {
		s = "0" + s
	}
	return Hex2Bytes(s)
}

// CopyBytes returns a copy of b. A nil slice stays nil.
func CopyBytes(b []byte) []byte {
	if b == nil {
		return nil
	}
	return append([]byte{}, b...)
}

// has0xPrefix reports whether s starts with "0x" or "0X".
func has0xPrefix(s string) bool {
	return len(s) >= 2 && s[0] == '0' && (s[1] == 'x' || s[1] == 'X')
}

// isHex reports whether s is a non-empty, even-length string of hex digits.
func isHex(s string) bool {
	if len(s) == 0 || len(s)%2 != 0 {
		return false
	}
	for i := 0; i < len(s); i++ {
		if !isHexDigit(s[i]) {
			return false
		}
	}
	return true
}

// isHexDigit reports whether c is a hex digit in either letter case.
func isHexDigit(c byte) bool {
	return '0' <= c && c <= '9' || 'a' <= c && c <= 'f' || 'A' <= c && c <= 'F'
}

// Bytes2Hex returns the hex encoding of d without a prefix.
func Bytes2Hex(d []byte) string { return hex.EncodeToString(d) }

// Hex2Bytes decodes the unprefixed hex string s. Decoding stops at the first
// invalid digit and the bytes decoded so far are returned.
func Hex2Bytes(s string) []byte {
	b, _ := hex.DecodeString(s)
	return b
}

// Hex2BytesFixed decodes the unprefixed hex string s into exactly n bytes,
// right-aligned: shorter input is left-padded with zeros, longer input loses
// its leading bytes.
func Hex2BytesFixed(s string, n int) []byte {
	b := Hex2Bytes(s)
	if len(b) > n {
		return b[len(b)-n:]
	}
	return LeftPadBytes(b, n)
}

// RightPadBytes returns b extended with zeros to length n. If b already has
// at least n bytes it is returned unchanged.
func RightPadBytes(b []byte, n int) []byte {
	if n <= len(b) {
		return b
	}
	out := make([]byte, n)
	copy(out, b)
	return out
}

// LeftPadBytes returns b prefixed with zeros to length n. If b already has at
// least n bytes it is returned unchanged.
func LeftPadBytes(b []byte, n int) []byte {
	if n <= len(b) {
		return b
	}
	out := make([]byte, n)
	copy(out[n-len(b):], b)
	return out
}

// TrimLeftZeroes returns b without its leading zero bytes.
func TrimLeftZeroes(b []byte) []byte {
	i := 0
	for i < len(b) && b[i] == 0 {
		i++
	}
	return b[i:]
}

// TrimRightZeroes returns b without its trailing zero bytes.
func TrimRightZeroes(b []byte) []byte {
	return b[:len(bytes.TrimRight(b, "\x00"))]
}

// Uint64ToBytes converts a uint64 to an 8-byte little-endian slice.
func Uint64ToBytes(u uint64) []byte {
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, u)
	return b
}

// BytesToUint64 converts an 8-byte little-endian slice to a uint64.
func BytesToUint64(b []byte) uint64 {
	return binary.LittleEndian.Uint64(b)
}
