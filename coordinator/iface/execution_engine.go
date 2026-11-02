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

import "context"

// ExecutionEngine is the subset of the Iguazu execution-layer RPC
// (iguazu-engine / gwat) consumed by iguazu-consensus and the iguazu-coordinator
// sidecar.
//
// The interface uses iguazu-types values so the consumer modules stay free of
// GPL/LGPL dependencies. Concrete implementations live in the coordinator
// and the iguazu-coordinator sidecar, both of which wrap the JSON-RPC client
// against the execution process.
//
// The initial surface exposes only a liveness probe; the 7 DAG-specific
// JSON-RPC endpoints (dag_finalize, dag_getOptimisticSpines,
// dag_getCandidates, dag_coordinatedState, dag_syncSlotInfo,
// dag_validateSpines, dag_syncSpines) are added as the corresponding call
// sites are migrated out of the coordinator.
type ExecutionEngine interface {
	Ping(ctx context.Context) error
}
