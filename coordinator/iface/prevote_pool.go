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

package iface

import (
	"context"

	"github.com/LFDT-Iguazu/iguazu-types/coordinator/types"
)

// PreVotePool is the in-memory cache of Iguazu prevote messages shared
// between the P2P validation path (sync layer), the RPC layer (queries) and
// the block-proposer path. The concrete implementation lives in iguazu-consensus
// and is consumed both by the coordinator process (P2P ingress) and the
// iguazu-coordinator sidecar (aggregation + query path).
type PreVotePool interface {
	// HasPrevote reports whether the pool already holds a prevote whose
	// aggregation bits cover those of pv (i.e. there is no new information
	// to record).
	HasPrevote(pv *types.PreVote) (bool, error)

	// SavePrevote inserts pv into the pool, deduplicating against the
	// "seen-bits" index so that repeated deliveries are ignored.
	SavePrevote(pv *types.PreVote) error

	// GetPrevoteBySlot returns all prevotes stored for the given slot.
	// The returned slice is a live view and should be treated as read-only.
	GetPrevoteBySlot(ctx context.Context, slot types.Slot) []*types.PreVote

	// PurgeOutdatedPrevote drops every cached prevote whose slot is strictly
	// less than curSlot, moving their aggregation bits into the seen-bits
	// index for the remainder of the cache TTL.
	PurgeOutdatedPrevote(curSlot types.Slot) error
}
