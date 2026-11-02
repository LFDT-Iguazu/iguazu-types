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

package hexutil

import (
	"encoding/json"
	"math/big"
	"reflect"
)

// Bytes is a byte slice that marshals to and from a "0x"-prefixed hex string.
type Bytes []byte

// Big is a big integer that marshals to and from a hex quantity of at most
// 256 bits.
type Big big.Int

// Uint64 marshals to and from a hex quantity of at most 64 bits.
type Uint64 uint64

// Uint8 marshals to and from a hex quantity of at most 8 bits.
type Uint8 uint8

var (
	bytesType  = reflect.TypeOf(Bytes(nil))
	bigType    = reflect.TypeOf((*Big)(nil))
	uint64Type = reflect.TypeOf(Uint64(0))
	uint8Type  = reflect.TypeOf(Uint8(0))
)

// jsonString returns the contents of a JSON string literal, or a type error
// if input is not a string.
func jsonString(input []byte, typ reflect.Type) ([]byte, error) {
	if len(input) < 2 || input[0] != '"' || input[len(input)-1] != '"' {
		return nil, &json.UnmarshalTypeError{Value: "non-string", Type: typ}
	}
	return input[1 : len(input)-1], nil
}

// asTypeError converts decoding errors of this package into JSON type errors,
// so that the JSON decoder reports the offending Go type.
func asTypeError(err error, typ reflect.Type) error {
	if _, ok := err.(hexError); ok {
		return &json.UnmarshalTypeError{Value: err.Error(), Type: typ}
	}
	return err
}

// unmarshalJSONText decodes a JSON string with the given text decoder.
func unmarshalJSONText(input []byte, typ reflect.Type, decode func([]byte) error) error {
	text, err := jsonString(input, typ)
	if err != nil {
		return err
	}
	return asTypeError(decode(text), typ)
}

// MarshalText implements encoding.TextMarshaler.
func (b Bytes) MarshalText() ([]byte, error) {
	return []byte(Encode(b)), nil
}

// UnmarshalJSON implements json.Unmarshaler.
func (b *Bytes) UnmarshalJSON(input []byte) error {
	return unmarshalJSONText(input, bytesType, b.UnmarshalText)
}

// UnmarshalText implements encoding.TextUnmarshaler.
func (b *Bytes) UnmarshalText(input []byte) error {
	digits, err := byteStringDigits(input, true)
	if err != nil {
		return err
	}
	out, err := decodeDigits(digits)
	if err != nil {
		return err
	}
	*b = out
	return nil
}

// String returns the hex encoding of b.
func (b Bytes) String() string {
	return Encode(b)
}

// MarshalText implements encoding.TextMarshaler.
func (b Big) MarshalText() ([]byte, error) {
	return []byte(EncodeBig((*big.Int)(&b))), nil
}

// UnmarshalJSON implements json.Unmarshaler.
func (b *Big) UnmarshalJSON(input []byte) error {
	return unmarshalJSONText(input, bigType, b.UnmarshalText)
}

// UnmarshalText implements encoding.TextUnmarshaler.
func (b *Big) UnmarshalText(input []byte) error {
	digits, err := quantityDigits(input, false)
	if err != nil {
		return err
	}
	v, err := parseBig(digits)
	if err != nil {
		return err
	}
	(*big.Int)(b).Set(v)
	return nil
}

// ToInt converts b to a *big.Int.
func (b *Big) ToInt() *big.Int {
	return (*big.Int)(b)
}

// String returns the hex encoding of b.
func (b *Big) String() string {
	return EncodeBig(b.ToInt())
}

// MarshalText implements encoding.TextMarshaler.
func (i Uint64) MarshalText() ([]byte, error) {
	return []byte(EncodeUint64(uint64(i))), nil
}

// UnmarshalJSON implements json.Unmarshaler.
func (i *Uint64) UnmarshalJSON(input []byte) error {
	return unmarshalJSONText(input, uint64Type, i.UnmarshalText)
}

// UnmarshalText implements encoding.TextUnmarshaler.
func (i *Uint64) UnmarshalText(input []byte) error {
	digits, err := quantityDigits(input, false)
	if err != nil {
		return err
	}
	v, err := parseUint(digits, 64, ErrUint64Range, false)
	if err != nil {
		return err
	}
	*i = Uint64(v)
	return nil
}

// String returns the hex encoding of i.
func (i Uint64) String() string {
	return EncodeUint64(uint64(i))
}

// MarshalText implements encoding.TextMarshaler.
func (i Uint8) MarshalText() ([]byte, error) {
	return []byte(EncodeUint64(uint64(i))), nil
}

// UnmarshalJSON implements json.Unmarshaler.
func (i *Uint8) UnmarshalJSON(input []byte) error {
	return unmarshalJSONText(input, uint8Type, i.UnmarshalText)
}

// UnmarshalText implements encoding.TextUnmarshaler.
func (i *Uint8) UnmarshalText(input []byte) error {
	digits, err := quantityDigits(input, false)
	if err != nil {
		return err
	}
	// Parse as a 64-bit quantity first, then check the 8-bit range.
	v, err := parseUint(digits, 64, ErrUint8Range, false)
	if err != nil {
		return err
	}
	if v > 0xff {
		return ErrUint8Range
	}
	*i = Uint8(v)
	return nil
}

// String returns the hex encoding of i.
func (i Uint8) String() string {
	return EncodeUint64(uint64(i))
}

// UnmarshalFixedJSON decodes a JSON hex string into out, which must have
// exactly the length of the decoded bytes. typ is reported in errors.
func UnmarshalFixedJSON(typ reflect.Type, input, out []byte) error {
	return unmarshalJSONText(input, typ, func(text []byte) error {
		return UnmarshalFixedText(typ.String(), text, out)
	})
}

// UnmarshalFixedText decodes a "0x"-prefixed hex string into out, which must
// have exactly the length of the decoded bytes.
func UnmarshalFixedText(typeName string, input, out []byte) error {
	return unmarshalFixed(typeName, input, out, true)
}

// UnmarshalFixedUnprefixedText is like UnmarshalFixedText but the "0x" prefix
// is optional.
func UnmarshalFixedUnprefixedText(typeName string, input, out []byte) error {
	return unmarshalFixed(typeName, input, out, false)
}

func unmarshalFixed(typeName string, input, out []byte, requirePrefix bool) error {
	digits, err := byteStringDigits(input, requirePrefix)
	if err != nil {
		return err
	}
	if len(digits) != 2*len(out) {
		return lengthError(len(digits), len(out), typeName)
	}
	decoded, err := decodeDigits(digits)
	if err != nil {
		return err
	}
	copy(out, decoded)
	return nil
}
