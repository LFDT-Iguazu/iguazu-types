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

// Attestation exposes the Iguazu-relevant fields of a beacon attestation
// without importing Prysm proto types.
type Attestation interface {
	// BeaconBlockRoot is the block root that this attestation is voting for.
	BeaconBlockRoot() []byte
	// AggregationBits is the committee participation bitlist.
	AggregationBits() []byte
	// AttestationSlot is the slot the committee was assigned to attest for.
	AttestationSlot() types.Slot
	// CommitteeIndex is the index of the committee within the slot.
	CommitteeIndex() types.CommitteeIndex
}

// BeaconBlockBody exposes the fields from a beacon block body that
// iguazu-consensus processors consume.
type BeaconBlockBody interface {
	Attestations() []Attestation
	// Eth1Candidates returns the spine-candidate bytes stored in Eth1Data.
	Eth1Candidates() []byte
}

// BeaconBlock is the read-only view of a single beacon block.
type BeaconBlock interface {
	Slot() types.Slot
	ParentRoot() []byte
	Body() BeaconBlockBody
}

// SignedBeaconBlock is the top-level wrapper carrying a beacon block and its
// BLS signature.  iguazu-consensus only needs the inner block.
type SignedBeaconBlock interface {
	Block() BeaconBlock
}

// CommitteeCounter abstracts the calculation of how many attestors are
// assigned to a given slot across all committees.  The concrete
// implementation lives in the coordinator (which owns the full shuffling
// algorithm); iguazu-consensus injects this interface to stay free of Prysm
// dependencies.
type CommitteeCounter interface {
	// CountSlotAttestors returns the total number of validators expected to
	// attest in the given slot (sum of all committee sizes).
	CountSlotAttestors(ctx context.Context, slot types.Slot) (int, error)
}
