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

// SpinesSeq is a packed sequence of 32-byte GWAT spine hashes concatenated
// into a single byte slice (len is a multiple of 32).
type SpinesSeq struct {
	Spines []byte
}

// SpineData carries the Iguazu DAG spine information stored inside the
// beacon state.
//
// Field semantics mirror coordinator/proto/prysm/v1alpha1/beacon_state.proto:
//   - Spines:       candidate GWAT spines for the current slot
//   - Prefix:       common prefix of spines across attestors
//   - Finalization: spines queued for finalization
//   - CpFinalized:  spines finalized between finalized and justified checkpoints
//   - ParentSpines: spine sequences of ancestors (excluding finalization/prefix)
type SpineData struct {
	Spines       []byte
	Prefix       []byte
	Finalization []byte
	CpFinalized  []byte
	ParentSpines []*SpinesSeq
}

// Copy returns a deep copy of the spine data.
func (s *SpineData) Copy() *SpineData {
	if s == nil {
		return nil
	}
	cp := &SpineData{
		Spines:       append([]byte(nil), s.Spines...),
		Prefix:       append([]byte(nil), s.Prefix...),
		Finalization: append([]byte(nil), s.Finalization...),
		CpFinalized:  append([]byte(nil), s.CpFinalized...),
	}
	if s.ParentSpines != nil {
		cp.ParentSpines = make([]*SpinesSeq, len(s.ParentSpines))
		for i, p := range s.ParentSpines {
			if p == nil {
				continue
			}
			cp.ParentSpines[i] = &SpinesSeq{Spines: append([]byte(nil), p.Spines...)}
		}
	}
	return cp
}
