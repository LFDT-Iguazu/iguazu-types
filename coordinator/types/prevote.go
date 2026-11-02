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

// PreVoteData is an attestor's view of the spine-candidate sequence that
// they intend to vote on for a given slot/committee.
type PreVoteData struct {
	Slot       Slot
	Index      CommitteeIndex
	Candidates []byte
}

// Copy returns a deep copy of the prevote data.
func (d *PreVoteData) Copy() *PreVoteData {
	if d == nil {
		return nil
	}
	return &PreVoteData{
		Slot:       d.Slot,
		Index:      d.Index,
		Candidates: append([]byte(nil), d.Candidates...),
	}
}

// PreVote is an aggregated prevote message signed by a committee.
//
// AggregationBits is stored as a raw SSZ bitlist payload (the last byte
// encodes the bit-length sentinel). Signature is a 96-byte BLS aggregate.
type PreVote struct {
	AggregationBits []byte
	Data            *PreVoteData
	Signature       []byte
}

// Copy returns a deep copy of the prevote.
func (p *PreVote) Copy() *PreVote {
	if p == nil {
		return nil
	}
	return &PreVote{
		AggregationBits: append([]byte(nil), p.AggregationBits...),
		Data:            p.Data.Copy(),
		Signature:       append([]byte(nil), p.Signature...),
	}
}

// IndexedPreVote is the indexed form of a prevote carrying the explicit list
// of attesting validator indices (used for signature verification).
type IndexedPreVote struct {
	AttestingIndices []uint64
	Data             *PreVoteData
	Signature        []byte
}

// Copy returns a deep copy of the indexed prevote.
func (p *IndexedPreVote) Copy() *IndexedPreVote {
	if p == nil {
		return nil
	}
	return &IndexedPreVote{
		AttestingIndices: append([]uint64(nil), p.AttestingIndices...),
		Data:             p.Data.Copy(),
		Signature:        append([]byte(nil), p.Signature...),
	}
}
