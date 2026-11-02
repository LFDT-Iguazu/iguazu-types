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

// SpineStore is the minimal persistence surface that iguazu-consensus / the
// iguazu-coordinator sidecar need to load and save the GWAT spine cache.
// Concrete implementations wrap the coordinator's BoltDB store or the
// sidecar's own key-value backend.
type SpineStore interface {
	SaveSpine(ctx context.Context, key [32]byte, value []byte) error
	Spine(ctx context.Context, key [32]byte) ([]byte, error)
	DeleteSpine(ctx context.Context, key [32]byte) error
}
