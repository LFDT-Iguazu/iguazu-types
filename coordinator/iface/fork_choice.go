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

package iface

import (
	"context"

	"github.com/LFDT-Iguazu/iguazu-types/common"
	"github.com/LFDT-Iguazu/iguazu-types/coordinator/types"
)

// ForkChoiceNode exposes the Iguazu-relevant spine data stored on a single
// fork-choice tree node.  The concrete implementation wraps the protoarray Node.
type ForkChoiceNode interface {
	// Root is the beacon block root identifying this node.
	Root() [32]byte
	// Slot is the beacon slot of the block.
	Slot() types.Slot

	// SpinesFinalized returns the cp-finalized spine sequence (cpFinalized field).
	// May be empty for genesis nodes.
	SpinesFinalized() common.HashArray
	// SpinesFinalization returns the accumulated finalization spine sequence.
	SpinesFinalization() common.HashArray
	// SpinesPrefix returns the node's current spine prefix.
	SpinesPrefix() common.HashArray
	// SpinesPublished returns the spine candidates published in the block that
	// produced this node (SpinesData.Spines).
	SpinesPublished() common.HashArray
}

// ForkChoiceFork is a single contiguous chain in the fork-choice tree,
// represented as an ordered slice of block roots from the tip (most recent)
// down to the common ancestor (oldest).
type ForkChoiceFork interface {
	// Roots returns block roots in tip-to-ancestor order:
	// Roots()[0] is the leaf; Roots()[len-1] is the oldest ancestor.
	Roots() [][32]byte
	// NodeByRoot returns the ForkChoiceNode for the given root, or nil.
	NodeByRoot(root [32]byte) ForkChoiceNode
}

// ForkChoice is the minimal view of the fork-choice data structure that
// iguazu-consensus needs.  The concrete implementation lives in the coordinator's
// protoarray package; only the T(G)-tree scoring logic is ported here.
type ForkChoice interface {
	// HasNode reports whether the given block root is tracked by fork choice.
	HasNode(root [32]byte) bool

	// GetForks returns all distinct forks currently tracked by the store.
	GetForks() []ForkChoiceFork

	// HeadBySubset runs the LMD-GHOST algorithm restricted to the given set of
	// acceptable block roots, starting from jCpRoot as the justified-checkpoint
	// root.  It returns the best head root among the acceptable candidates.
	HeadBySubset(ctx context.Context, acceptableRoots map[[32]byte]struct{}, jCpRoot [32]byte) ([32]byte, error)
}
