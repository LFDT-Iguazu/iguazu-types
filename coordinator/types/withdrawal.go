// Copyright 2026 Digital Clever Solution LLC
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

// Withdrawal is a block-level operation signaling a validator's intent to
// withdraw (reward-only if still active, full balance otherwise).
//
// InitTxHash is the hash of the initiating transaction (Iguazu extension)
// used to deduplicate withdrawals across re-orgs.
type Withdrawal struct {
	PublicKey      []byte
	ValidatorIndex ValidatorIndex
	Amount         uint64
	InitTxHash     []byte
	Epoch          Epoch
}

// Copy returns a deep copy of the withdrawal.
func (w *Withdrawal) Copy() *Withdrawal {
	if w == nil {
		return nil
	}
	return &Withdrawal{
		PublicKey:      append([]byte(nil), w.PublicKey...),
		ValidatorIndex: w.ValidatorIndex,
		Amount:         w.Amount,
		InitTxHash:     append([]byte(nil), w.InitTxHash...),
		Epoch:          w.Epoch,
	}
}
