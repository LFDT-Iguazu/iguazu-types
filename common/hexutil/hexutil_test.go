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
	"errors"
	"math/big"
	"testing"
)

func TestEncodeDecode(t *testing.T) {
	cases := []struct {
		in   string
		want []byte
		err  error
	}{
		{"0x", []byte{}, nil},
		{"0x00ff", []byte{0x00, 0xff}, nil},
		{"0XAbCd", []byte{0xab, 0xcd}, nil},
		{"", nil, ErrEmptyString},
		{"00ff", nil, ErrMissingPrefix},
		{"0x0", nil, ErrOddLength},
		{"0xzz", nil, ErrSyntax},
		{"0xg", nil, ErrSyntax},
	}
	for _, c := range cases {
		got, err := Decode(c.in)
		if !errors.Is(err, c.err) {
			t.Errorf("Decode(%q) error = %v, want %v", c.in, err, c.err)
			continue
		}
		if err == nil && string(got) != string(c.want) {
			t.Errorf("Decode(%q) = %x, want %x", c.in, got, c.want)
		}
		if err == nil && Encode(got) != lower(c.in) {
			t.Errorf("Encode(Decode(%q)) = %q", c.in, Encode(got))
		}
	}
}

func lower(s string) string {
	b := []byte(s)
	for i, c := range b {
		if 'A' <= c && c <= 'Z' {
			b[i] = c + 'a' - 'A'
		}
	}
	return string(b)
}

func TestQuantities(t *testing.T) {
	u64 := []struct {
		in   string
		want uint64
		err  error
	}{
		{"0x0", 0, nil},
		{"0x1f", 0x1f, nil},
		{"0xffffffffffffffff", 1<<64 - 1, nil},
		{"0x10000000000000000", 0, ErrUint64Range},
		{"0x01", 0, ErrLeadingZero},
		{"0x", 0, ErrEmptyNumber},
		{"1", 0, ErrMissingPrefix},
		{"0xz", 0, ErrSyntax},
	}
	for _, c := range u64 {
		got, err := DecodeUint64(c.in)
		if !errors.Is(err, c.err) || (err == nil && got != c.want) {
			t.Errorf("DecodeUint64(%q) = %d, %v; want %d, %v", c.in, got, err, c.want, c.err)
		}
	}
	if s := EncodeUint64(255); s != "0xff" {
		t.Errorf("EncodeUint64(255) = %q", s)
	}
	if s := EncodeBig(big.NewInt(-255)); s != "-0xff" {
		t.Errorf("EncodeBig(-255) = %q", s)
	}
	if s := EncodeBig(new(big.Int)); s != "0x0" {
		t.Errorf("EncodeBig(0) = %q", s)
	}
	max256 := "0x" + string(repeat('f', 64))
	if v, err := DecodeBig(max256); err != nil || v.BitLen() != 256 {
		t.Errorf("DecodeBig(2^256-1) = %v, %v", v, err)
	}
	if _, err := DecodeBig("0x1" + string(repeat('0', 64))); !errors.Is(err, ErrBig256Range) {
		t.Errorf("DecodeBig(2^256) error = %v", err)
	}
}

func repeat(c byte, n int) []byte {
	b := make([]byte, n)
	for i := range b {
		b[i] = c
	}
	return b
}

func TestJSON(t *testing.T) {
	type doc struct {
		B  Bytes  `json:"b"`
		U  Uint64 `json:"u"`
		U8 Uint8  `json:"u8"`
		N  *Big   `json:"n"`
	}
	in := doc{B: Bytes{1, 2}, U: 300, U8: 7, N: (*Big)(big.NewInt(1 << 40))}
	enc, err := json.Marshal(in)
	if err != nil {
		t.Fatal(err)
	}
	const want = `{"b":"0x0102","u":"0x12c","u8":"0x7","n":"0x10000000000"}`
	if string(enc) != want {
		t.Fatalf("Marshal = %s, want %s", enc, want)
	}
	var out doc
	if err := json.Unmarshal(enc, &out); err != nil {
		t.Fatal(err)
	}
	if string(out.B) != string(in.B) || out.U != in.U || out.U8 != in.U8 || out.N.ToInt().Cmp(in.N.ToInt()) != 0 {
		t.Fatalf("round trip = %+v, want %+v", out, in)
	}

	// An empty string decodes to the zero value; non-strings are type errors.
	var u Uint64
	if err := json.Unmarshal([]byte(`""`), &u); err != nil || u != 0 {
		t.Errorf(`Uint64 from "" = %d, %v`, u, err)
	}
	var typeErr *json.UnmarshalTypeError
	if err := json.Unmarshal([]byte(`5`), &u); !errors.As(err, &typeErr) {
		t.Errorf("Uint64 from number: error = %v, want a type error", err)
	}
	var u8 Uint8
	if err := json.Unmarshal([]byte(`"0x100"`), &u8); !errors.As(err, &typeErr) {
		t.Errorf("Uint8 out of range: error = %v, want a type error", err)
	}
}

func TestUnmarshalFixed(t *testing.T) {
	out := make([]byte, 2)
	if err := UnmarshalFixedText("T", []byte("0x0102"), out); err != nil || out[0] != 1 || out[1] != 2 {
		t.Errorf("UnmarshalFixedText = %x, %v", out, err)
	}
	if err := UnmarshalFixedText("T", []byte("0102"), out); !errors.Is(err, ErrMissingPrefix) {
		t.Errorf("UnmarshalFixedText without prefix: error = %v", err)
	}
	if err := UnmarshalFixedUnprefixedText("T", []byte("0304"), out); err != nil || out[0] != 3 || out[1] != 4 {
		t.Errorf("UnmarshalFixedUnprefixedText = %x, %v", out, err)
	}
	err := UnmarshalFixedText("T", []byte("0x01"), out)
	if err == nil || err.Error() != "hex string has length 2, want 4 for T" {
		t.Errorf("UnmarshalFixedText wrong length: error = %v", err)
	}
	if err := UnmarshalFixedJSON(bytesType, []byte(`"0x0506"`), out); err != nil || out[0] != 5 {
		t.Errorf("UnmarshalFixedJSON = %x, %v", out, err)
	}
}
