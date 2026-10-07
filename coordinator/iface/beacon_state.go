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
	"github.com/LFDT-Iguazu/iguazu-types/coordinator/types"
)

// BeaconStateReader exposes the read-only Iguazu-relevant slice of the
// beacon state to iguazu-consensus logic.
//
// The surface intentionally excludes Ethereum-specific fields that
// iguazu-consensus does not need to inspect; additional getters are added as
// consensus processors are migrated from the coordinator.
type BeaconStateReader interface {
	Slot() types.Slot

	SpineData() *types.SpineData
	BlockVoting() []*types.BlockVoting

	FinalizedCheckpoint() *types.Checkpoint
	CurrentJustifiedCheckpoint() *types.Checkpoint
	PreviousJustifiedCheckpoint() *types.Checkpoint

	NumValidators() int
	ValidatorAtIndex(idx types.ValidatorIndex) (*types.Validator, error)
	ValidatorAtIndexReadOnly(idx types.ValidatorIndex) (*types.Validator, error)
	Balances() []uint64
	BalanceAtIndex(idx types.ValidatorIndex) (uint64, error)
}

// BeaconStateWriter exposes the mutating operations that iguazu-consensus needs
// to apply block- or epoch-level state transitions.
type BeaconStateWriter interface {
	SetSpineData(*types.SpineData) error
	SetBlockVoting([]*types.BlockVoting) error

	SetFinalizedCheckpoint(*types.Checkpoint) error
	SetCurrentJustifiedCheckpoint(*types.Checkpoint) error

	UpdateValidatorAtIndex(idx types.ValidatorIndex, val *types.Validator) error
	UpdateBalancesAtIndex(idx types.ValidatorIndex, balance uint64) error
}

// BeaconState is the full read/write view used by iguazu-consensus processors.
type BeaconState interface {
	BeaconStateReader
	BeaconStateWriter
}

// DagBeaconState extends BeaconState with the Eth1Candidates field needed by
// the DAG-consensus processor.  The Eth1Data.Candidates field stores the
// current slot's spine candidates and is updated by ProcessDagConsensus.
type DagBeaconState interface {
	BeaconState
	Eth1Candidates() []byte
	SetEth1Candidates([]byte) error
}
