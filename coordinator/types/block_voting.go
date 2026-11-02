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

// CommitteeVote aggregates the votes of a committee that agreed on the same
// spine-candidate set for a given slot.
//
// AggregationBits is stored as a raw SSZ bitlist payload (the last byte
// encodes the bit-length sentinel).
type CommitteeVote struct {
	AggregationBits []byte
	Slot            Slot
	Index           CommitteeIndex
}

// Copy returns a deep copy of the committee vote.
func (c *CommitteeVote) Copy() *CommitteeVote {
	if c == nil {
		return nil
	}
	return &CommitteeVote{
		AggregationBits: append([]byte(nil), c.AggregationBits...),
		Slot:            c.Slot,
		Index:           c.Index,
	}
}

// BlockVoting represents an in-progress vote on the Iguazu spine
// candidates for a specific beacon block.
type BlockVoting struct {
	Root       []byte
	Slot       Slot
	Candidates []byte
	Votes      []*CommitteeVote
}

// Copy returns a deep copy of the block voting entry.
func (b *BlockVoting) Copy() *BlockVoting {
	if b == nil {
		return nil
	}
	cp := &BlockVoting{
		Root:       append([]byte(nil), b.Root...),
		Slot:       b.Slot,
		Candidates: append([]byte(nil), b.Candidates...),
	}
	if b.Votes != nil {
		cp.Votes = make([]*CommitteeVote, len(b.Votes))
		for i, v := range b.Votes {
			cp.Votes[i] = v.Copy()
		}
	}
	return cp
}
