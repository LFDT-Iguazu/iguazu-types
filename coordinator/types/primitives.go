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

package types

import (
	"errors"
	"math"
	"strconv"
)

// ErrOverflow is returned by Safe* arithmetic helpers when the result would
// exceed the uint64 range.
var ErrOverflow = errors.New("coordinator/types: uint64 overflow")

// ErrDivByZero is returned by Safe* division helpers when the divisor is zero.
var ErrDivByZero = errors.New("coordinator/types: division by zero")

// Slot represents a beacon-chain slot number.
type Slot uint64

// Epoch represents a beacon-chain epoch number.
type Epoch uint64

// ValidatorIndex is a sequential index of a validator in the beacon state.
type ValidatorIndex uint64

// CommitteeIndex is an index of a committee within a slot.
type CommitteeIndex uint64

// FarFutureEpoch is the sentinel value used for validators that have not yet
// activated or exited. It mirrors the consensus-spec constant of the same name.
const FarFutureEpoch = Epoch(math.MaxUint64)

// ----- Slot arithmetic ------------------------------------------------------

func (s Slot) Add(x uint64) Slot { return Slot(uint64(s) + x) }
func (s Slot) Sub(x uint64) Slot { return Slot(uint64(s) - x) }
func (s Slot) Mul(x uint64) Slot { return Slot(uint64(s) * x) }
func (s Slot) Div(x uint64) Slot { return Slot(uint64(s) / x) }
func (s Slot) Mod(x uint64) Slot { return Slot(uint64(s) % x) }

func (s Slot) SafeAdd(x uint64) (Slot, error) {
	r, err := safeAdd(uint64(s), x)
	return Slot(r), err
}
func (s Slot) SafeSub(x uint64) (Slot, error) {
	r, err := safeSub(uint64(s), x)
	return Slot(r), err
}
func (s Slot) SafeMul(x uint64) (Slot, error) {
	r, err := safeMul(uint64(s), x)
	return Slot(r), err
}
func (s Slot) SafeDiv(x uint64) (Slot, error) {
	r, err := safeDiv(uint64(s), x)
	return Slot(r), err
}

func (s Slot) String() string { return strconv.FormatUint(uint64(s), 10) }
func (s Slot) HashTreeRoot() ([32]byte, error) {
	return uint64HashTreeRoot(uint64(s)), nil
}

// ----- Epoch arithmetic -----------------------------------------------------

func (e Epoch) Add(x uint64) Epoch { return Epoch(uint64(e) + x) }
func (e Epoch) Sub(x uint64) Epoch { return Epoch(uint64(e) - x) }
func (e Epoch) Mul(x uint64) Epoch { return Epoch(uint64(e) * x) }
func (e Epoch) Div(x uint64) Epoch { return Epoch(uint64(e) / x) }
func (e Epoch) Mod(x uint64) Epoch { return Epoch(uint64(e) % x) }

func (e Epoch) SafeAdd(x uint64) (Epoch, error) {
	r, err := safeAdd(uint64(e), x)
	return Epoch(r), err
}
func (e Epoch) SafeSub(x uint64) (Epoch, error) {
	r, err := safeSub(uint64(e), x)
	return Epoch(r), err
}
func (e Epoch) SafeMul(x uint64) (Epoch, error) {
	r, err := safeMul(uint64(e), x)
	return Epoch(r), err
}
func (e Epoch) SafeDiv(x uint64) (Epoch, error) {
	r, err := safeDiv(uint64(e), x)
	return Epoch(r), err
}

func (e Epoch) String() string { return strconv.FormatUint(uint64(e), 10) }
func (e Epoch) HashTreeRoot() ([32]byte, error) {
	return uint64HashTreeRoot(uint64(e)), nil
}

// ----- ValidatorIndex / CommitteeIndex --------------------------------------

func (v ValidatorIndex) String() string { return strconv.FormatUint(uint64(v), 10) }
func (v ValidatorIndex) HashTreeRoot() ([32]byte, error) {
	return uint64HashTreeRoot(uint64(v)), nil
}

func (c CommitteeIndex) String() string { return strconv.FormatUint(uint64(c), 10) }
func (c CommitteeIndex) HashTreeRoot() ([32]byte, error) {
	return uint64HashTreeRoot(uint64(c)), nil
}

// ----- helpers --------------------------------------------------------------

func safeAdd(a, b uint64) (uint64, error) {
	if a > math.MaxUint64-b {
		return 0, ErrOverflow
	}
	return a + b, nil
}

func safeSub(a, b uint64) (uint64, error) {
	if b > a {
		return 0, ErrOverflow
	}
	return a - b, nil
}

func safeMul(a, b uint64) (uint64, error) {
	if a == 0 || b == 0 {
		return 0, nil
	}
	if a > math.MaxUint64/b {
		return 0, ErrOverflow
	}
	return a * b, nil
}

func safeDiv(a, b uint64) (uint64, error) {
	if b == 0 {
		return 0, ErrDivByZero
	}
	return a / b, nil
}

func uint64HashTreeRoot(v uint64) [32]byte {
	var out [32]byte
	// SSZ: uint64 serializes little-endian, right-padded with zero bytes.
	out[0] = byte(v)
	out[1] = byte(v >> 8)
	out[2] = byte(v >> 16)
	out[3] = byte(v >> 24)
	out[4] = byte(v >> 32)
	out[5] = byte(v >> 40)
	out[6] = byte(v >> 48)
	out[7] = byte(v >> 56)
	return out
}
