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

// WithdrawalOp represents a finalized withdrawal operation attached to a
// validator record.
type WithdrawalOp struct {
	Amount uint64
	Hash   []byte
	Slot   Slot
}

// Copy returns a deep copy of the withdrawal op.
func (w *WithdrawalOp) Copy() *WithdrawalOp {
	if w == nil {
		return nil
	}
	return &WithdrawalOp{
		Amount: w.Amount,
		Hash:   append([]byte(nil), w.Hash...),
		Slot:   w.Slot,
	}
}

// Validator mirrors the beacon-state validator record with Iguazu-specific
// extensions (creator address, withdrawal credentials as 20-byte address,
// activation/exit tx hashes, finalized withdrawal ops).
type Validator struct {
	PublicKey                  []byte
	CreatorAddress             []byte
	WithdrawalCredentials      []byte
	EffectiveBalance           uint64
	Slashed                    bool
	ActivationEligibilityEpoch Epoch
	ActivationEpoch            Epoch
	ExitEpoch                  Epoch
	WithdrawableEpoch          Epoch
	ActivationHash             []byte
	ExitHash                   []byte
	WithdrawalOps              []*WithdrawalOp
}

// Copy returns a deep copy of the validator record.
func (v *Validator) Copy() *Validator {
	if v == nil {
		return nil
	}
	cp := &Validator{
		PublicKey:                  append([]byte(nil), v.PublicKey...),
		CreatorAddress:             append([]byte(nil), v.CreatorAddress...),
		WithdrawalCredentials:      append([]byte(nil), v.WithdrawalCredentials...),
		EffectiveBalance:           v.EffectiveBalance,
		Slashed:                    v.Slashed,
		ActivationEligibilityEpoch: v.ActivationEligibilityEpoch,
		ActivationEpoch:            v.ActivationEpoch,
		ExitEpoch:                  v.ExitEpoch,
		WithdrawableEpoch:          v.WithdrawableEpoch,
		ActivationHash:             append([]byte(nil), v.ActivationHash...),
		ExitHash:                   append([]byte(nil), v.ExitHash...),
	}
	if v.WithdrawalOps != nil {
		cp.WithdrawalOps = make([]*WithdrawalOp, len(v.WithdrawalOps))
		for i, op := range v.WithdrawalOps {
			cp.WithdrawalOps[i] = op.Copy()
		}
	}
	return cp
}
