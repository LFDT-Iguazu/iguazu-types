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
	"encoding/hex"
	"fmt"
	"strconv"
)

// setRightAligned copies src into the fixed-size array dst, aligned to the
// right. If src is longer than dst, its leading bytes are dropped.
func setRightAligned(dst, src []byte) {
	if len(src) > len(dst) {
		src = src[len(src)-len(dst):]
	}
	copy(dst[len(dst)-len(src):], src)
}

// scanFixed implements sql.Scanner for fixed-size byte arrays.
func scanFixed(dst []byte, src interface{}, typeName string) error {
	b, ok := src.([]byte)
	if !ok {
		return fmt.Errorf("can't scan %T into %s", src, typeName)
	}
	if len(b) != len(dst) {
		return fmt.Errorf("can't scan []byte of len %d into %s, want %d", len(b), typeName, len(dst))
	}
	copy(dst, b)
	return nil
}

// terminalString abbreviates b to its first and last three bytes.
func terminalString(b []byte) string {
	if len(b) <= 6 {
		return hex.EncodeToString(b)
	}
	return hex.EncodeToString(b[:3]) + ".." + hex.EncodeToString(b[len(b)-3:])
}

// formatFixed implements fmt.Formatter for fixed-size byte arrays:
//
//	%v, %s  the display form (for example the 0x-prefixed hex string)
//	%q      the quoted display form
//	%x, %X  hex digits in lower/upper case, with a 0x/0X prefix for %#x/%#X
//	%d      the bytes as a list of decimal numbers
func formatFixed(s fmt.State, verb rune, b []byte, display, kind string) {
	switch verb {
	case 'v', 's':
		_, _ = s.Write([]byte(display))
	case 'q':
		_, _ = s.Write([]byte(strconv.Quote(display)))
	case 'x', 'X':
		digits := hex.EncodeToString(b)
		prefix := "0x"
		if verb == 'X' {
			digits = toUpperHex(digits)
			prefix = "0X"
		}
		if s.Flag('#') {
			_, _ = s.Write([]byte(prefix))
		}
		_, _ = s.Write([]byte(digits))
	case 'd':
		// Print the bytes as a list of decimal numbers, not as a string.
		_, _ = fmt.Fprint(s, b) //nolint:staticcheck // QF1010: converting to string would change the output
	default:
		_, _ = fmt.Fprintf(s, "%%!%c(%s=%x)", verb, kind, b)
	}
}

func toUpperHex(s string) string {
	out := []byte(s)
	for i, c := range out {
		if 'a' <= c && c <= 'f' {
			out[i] = c - 'a' + 'A'
		}
	}
	return string(out)
}
