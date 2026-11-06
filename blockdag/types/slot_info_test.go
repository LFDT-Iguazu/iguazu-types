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

import (
	"encoding/json"
	"math"
	"testing"
	"time"
)

const (
	// siGenesis is an arbitrary genesis time used by the slot info tests.
	siGenesis = uint64(1_700_000_000)
	// siForkSlot is the slot the slot duration changes from 6 to 3 seconds at,
	// which mirrors the mainnet ForkSlotUpSecondsPerSlot setup.
	siForkSlot = uint64(100)
)

// siForked returns slot info with the slot duration fork enabled at siForkSlot.
func siForked() *SlotInfo {
	return &SlotInfo{
		GenesisTime:              siGenesis,
		SecondsPerSlot:           6,
		SlotsPerEpoch:            32,
		ForkSlotUpSecondsPerSlot: siForkSlot,
		UpSecondsPerSlot:         3,
	}
}

// siPlain returns slot info with the same base settings as siForked but with the
// slot duration fork disabled.
func siPlain() *SlotInfo {
	return &SlotInfo{
		GenesisTime:    siGenesis,
		SecondsPerSlot: 6,
		SlotsPerEpoch:  32,
	}
}

func startSlotTime(t *testing.T, si *SlotInfo, slot uint64) uint64 {
	t.Helper()
	tm, err := si.StartSlotTime(slot)
	if err != nil {
		t.Fatalf("StartSlotTime(%d): unexpected error: %v", slot, err)
	}
	return uint64(tm.Unix())
}

func TestSlotInfoIsForkUpSecondsPerSlot(t *testing.T) {
	tests := []struct {
		name string
		si   *SlotInfo
		slot uint64
		want bool
	}{
		{name: "before fork", si: siForked(), slot: siForkSlot - 1, want: false},
		{name: "at fork", si: siForked(), slot: siForkSlot, want: true},
		{name: "after fork", si: siForked(), slot: siForkSlot + 1, want: true},
		{name: "fork disabled", si: siPlain(), slot: siForkSlot + 1, want: false},
		{
			name: "fork slot unset means fork from genesis",
			si:   &SlotInfo{GenesisTime: siGenesis, SecondsPerSlot: 6, UpSecondsPerSlot: 3},
			slot: 0,
			want: true,
		},
		{
			name: "up value unset disables fork from genesis",
			si:   &SlotInfo{GenesisTime: siGenesis, SecondsPerSlot: 6},
			slot: 0,
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.si.IsForkUpSecondsPerSlot(tt.slot); got != tt.want {
				t.Fatalf("IsForkUpSecondsPerSlot(%d) = %v, want %v", tt.slot, got, tt.want)
			}
		})
	}
}

func TestSlotInfoSecondsPerSlotBySlot(t *testing.T) {
	si := siForked()
	tests := []struct {
		slot uint64
		want uint64
	}{
		{slot: 0, want: 6},
		{slot: siForkSlot - 1, want: 6},
		{slot: siForkSlot, want: 3},
		{slot: siForkSlot + 1000, want: 3},
	}
	for _, tt := range tests {
		if got := si.SecondsPerSlotBySlot(tt.slot); got != tt.want {
			t.Fatalf("SecondsPerSlotBySlot(%d) = %d, want %d", tt.slot, got, tt.want)
		}
	}
	if got := siPlain().SecondsPerSlotBySlot(siForkSlot + 1000); got != 6 {
		t.Fatalf("fork disabled: SecondsPerSlotBySlot = %d, want 6", got)
	}
}

func TestSlotInfoStartSlotTimeNoFork(t *testing.T) {
	si := siPlain()
	for _, slot := range []uint64{0, 1, 99, 100, 101, 1000} {
		want := siGenesis + slot*6
		if got := startSlotTime(t, si, slot); got != want {
			t.Fatalf("StartSlotTime(%d) = %d, want %d", slot, got, want)
		}
	}
}

// TestSlotInfoStartSlotTimeForkDisabledByMaxSlot makes sure the math.MaxUint64 fork
// slot used by the chain configs to switch the fork off does not overflow into a
// shifted timeline.
func TestSlotInfoStartSlotTimeForkDisabledByMaxSlot(t *testing.T) {
	si := siPlain()
	si.ForkSlotUpSecondsPerSlot = math.MaxUint64
	si.UpSecondsPerSlot = 3

	for _, slot := range []uint64{0, 1, 100, 1_000_000} {
		want := siGenesis + slot*6
		if got := startSlotTime(t, si, slot); got != want {
			t.Fatalf("StartSlotTime(%d) = %d, want %d", slot, got, want)
		}
		if got := si.SlotByTimestamp(want); got != slot {
			t.Fatalf("SlotByTimestamp(%d) = %d, want %d", want, got, slot)
		}
	}
}

func TestSlotInfoStartSlotTimeFork(t *testing.T) {
	si := siForked()
	forkTime := siGenesis + siForkSlot*6 // 600 seconds after genesis

	tests := []struct {
		name string
		slot uint64
		want uint64
	}{
		{name: "genesis slot", slot: 0, want: siGenesis},
		{name: "slot before fork", slot: siForkSlot - 1, want: siGenesis + 99*6},
		{name: "fork slot", slot: siForkSlot, want: forkTime},
		{name: "slot after fork", slot: siForkSlot + 1, want: forkTime + 3},
		{name: "far after fork", slot: siForkSlot + 100, want: forkTime + 300},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := startSlotTime(t, si, tt.slot); got != tt.want {
				t.Fatalf("StartSlotTime(%d) = %d, want %d", tt.slot, got, tt.want)
			}
		})
	}

	// the last pre-fork slot must last the full base slot duration
	if got := startSlotTime(t, si, siForkSlot) - startSlotTime(t, si, siForkSlot-1); got != 6 {
		t.Fatalf("duration of the last pre-fork slot = %d, want 6", got)
	}
}

// TestSlotInfoForkKeepsHistory is a regression test for the slot timeline being
// recalculated from genesis with the new slot duration: everything before the fork
// slot must stay exactly where it was before the fork was configured.
func TestSlotInfoForkKeepsHistory(t *testing.T) {
	forked, plain := siForked(), siPlain()

	for slot := uint64(0); slot < siForkSlot; slot++ {
		if got, want := startSlotTime(t, forked, slot), startSlotTime(t, plain, slot); got != want {
			t.Fatalf("StartSlotTime(%d) = %d with fork, %d without", slot, got, want)
		}
	}
	for ts := siGenesis; ts < siGenesis+siForkSlot*6; ts++ {
		if got, want := forked.SlotByTimestamp(ts), plain.SlotByTimestamp(ts); got != want {
			t.Fatalf("SlotByTimestamp(%d) = %d with fork, %d without", ts, got, want)
		}
	}
}

func TestSlotInfoSlotByTimestamp(t *testing.T) {
	si := siForked()
	forkTime := siGenesis + siForkSlot*6

	tests := []struct {
		name      string
		timestamp uint64
		want      uint64
	}{
		{name: "before genesis", timestamp: siGenesis - 1, want: 0},
		{name: "genesis", timestamp: siGenesis, want: 0},
		{name: "inside genesis slot", timestamp: siGenesis + 5, want: 0},
		{name: "second slot", timestamp: siGenesis + 6, want: 1},
		{name: "last second before fork", timestamp: forkTime - 1, want: siForkSlot - 1},
		{name: "fork time", timestamp: forkTime, want: siForkSlot},
		{name: "inside fork slot", timestamp: forkTime + 2, want: siForkSlot},
		{name: "slot after fork", timestamp: forkTime + 3, want: siForkSlot + 1},
		{name: "far after fork", timestamp: forkTime + 300, want: siForkSlot + 100},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := si.SlotByTimestamp(tt.timestamp); got != tt.want {
				t.Fatalf("SlotByTimestamp(%d) = %d, want %d", tt.timestamp, got, tt.want)
			}
		})
	}
}

// TestSlotInfoSlotTimeRoundTrip checks that SlotByTimestamp is the inverse of
// StartSlotTime on both sides of the fork and that slots do not overlap or leave gaps.
func TestSlotInfoSlotTimeRoundTrip(t *testing.T) {
	si := siForked()
	for slot := uint64(0); slot < siForkSlot*3; slot++ {
		start := startSlotTime(t, si, slot)
		next := startSlotTime(t, si, slot+1)

		if want := si.SecondsPerSlotBySlot(slot); next-start != want {
			t.Fatalf("slot %d lasts %d seconds, want %d", slot, next-start, want)
		}
		if got := si.SlotByTimestamp(start); got != slot {
			t.Fatalf("SlotByTimestamp(start of %d) = %d", slot, got)
		}
		if got := si.SlotByTimestamp(next - 1); got != slot {
			t.Fatalf("SlotByTimestamp(end of %d) = %d", slot, got)
		}
	}
}

func TestSlotInfoForkAtGenesis(t *testing.T) {
	si := siForked()
	si.ForkSlotUpSecondsPerSlot = 0

	for _, slot := range []uint64{0, 1, 100} {
		want := siGenesis + slot*3
		if got := startSlotTime(t, si, slot); got != want {
			t.Fatalf("StartSlotTime(%d) = %d, want %d", slot, got, want)
		}
		if got := si.SlotByTimestamp(want); got != slot {
			t.Fatalf("SlotByTimestamp(%d) = %d, want %d", want, got, slot)
		}
	}
}

func TestSlotInfoStartSlotTimeOverflow(t *testing.T) {
	if _, err := siPlain().StartSlotTime(math.MaxUint64); err == nil {
		t.Fatal("expected an overflow error without the fork")
	}
	si := siForked()
	if _, err := si.StartSlotTime(math.MaxUint64); err == nil {
		t.Fatal("expected an overflow error with the fork enabled")
	}
}

func TestSlotInfoCurrentSlot(t *testing.T) {
	// place genesis 1000 seconds in the past: 600 of them are pre-fork slots (100
	// slots of 6 seconds), the remaining 400 are post-fork slots of 3 seconds
	si := siForked()
	si.GenesisTime = uint64(time.Now().Unix()) - 1000
	want := siForkSlot + 400/3

	got := si.CurrentSlot()
	if got < want-1 || got > want+1 {
		t.Fatalf("CurrentSlot() = %d, want ~%d", got, want)
	}

	si.GenesisTime = uint64(time.Now().Unix()) + 1000
	if got := si.CurrentSlot(); got != 0 {
		t.Fatalf("CurrentSlot() before genesis = %d, want 0", got)
	}
}

func TestSlotInfoCopy(t *testing.T) {
	si := siForked()
	cpy := si.Copy()
	if cpy == si {
		t.Fatal("Copy returned the same pointer")
	}
	if *cpy != *si {
		t.Fatalf("Copy() = %+v, want %+v", cpy, si)
	}

	cpy.UpSecondsPerSlot = 1
	if si.UpSecondsPerSlot != 3 {
		t.Fatal("Copy shares state with the original")
	}

	var nilSi *SlotInfo
	if nilSi.Copy() != nil {
		t.Fatal("Copy of nil slot info must be nil")
	}
}

func TestSlotInfoMatchCoordinatorParams(t *testing.T) {
	si := siForked()

	tests := []struct {
		name  string
		coord *SlotInfo
		slot  uint64
		want  bool
	}{
		{
			name:  "coordinator reports base slot duration",
			coord: &SlotInfo{GenesisTime: siGenesis, SecondsPerSlot: 6, SlotsPerEpoch: 32},
			slot:  0,
			want:  true,
		},
		{
			name:  "coordinator reports base slot duration after the fork",
			coord: &SlotInfo{GenesisTime: siGenesis, SecondsPerSlot: 6, SlotsPerEpoch: 32},
			slot:  siForkSlot + 1,
			want:  true,
		},
		{
			name:  "coordinator switched to the new slot duration after the fork",
			coord: &SlotInfo{GenesisTime: siGenesis, SecondsPerSlot: 3, SlotsPerEpoch: 32},
			slot:  siForkSlot + 1,
			want:  true,
		},
		{
			name:  "coordinator reports the new slot duration before the fork",
			coord: &SlotInfo{GenesisTime: siGenesis, SecondsPerSlot: 3, SlotsPerEpoch: 32},
			slot:  siForkSlot - 1,
			want:  false,
		},
		{
			name:  "genesis time mismatch",
			coord: &SlotInfo{GenesisTime: siGenesis + 1, SecondsPerSlot: 6, SlotsPerEpoch: 32},
			slot:  0,
			want:  false,
		},
		{
			name:  "slots per epoch mismatch",
			coord: &SlotInfo{GenesisTime: siGenesis, SecondsPerSlot: 6, SlotsPerEpoch: 16},
			slot:  0,
			want:  false,
		},
		{
			name:  "unrelated slot duration",
			coord: &SlotInfo{GenesisTime: siGenesis, SecondsPerSlot: 4, SlotsPerEpoch: 32},
			slot:  siForkSlot + 1,
			want:  false,
		},
		{name: "nil coordinator info", coord: nil, slot: 0, want: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := si.MatchCoordinatorParams(tt.coord, tt.slot); got != tt.want {
				t.Fatalf("MatchCoordinatorParams() = %v, want %v", got, tt.want)
			}
		})
	}

	var nilSi *SlotInfo
	if nilSi.MatchCoordinatorParams(siForked(), 0) {
		t.Fatal("nil slot info must not match")
	}
}

// TestSlotInfoJSONOmitsDisabledFork makes sure a slot info without the fork keeps the
// JSON encoding it had before the fork fields were added.
func TestSlotInfoJSONOmitsDisabledFork(t *testing.T) {
	got, err := json.Marshal(siPlain())
	if err != nil {
		t.Fatalf("Marshal: unexpected error: %v", err)
	}
	want := `{"genesisTime":1700000000,"secondsPerSlot":6,"slotsPerEpoch":32}`
	if string(got) != want {
		t.Fatalf("Marshal() = %s, want %s", got, want)
	}

	var decoded SlotInfo
	data, err := json.Marshal(siForked())
	if err != nil {
		t.Fatalf("Marshal: unexpected error: %v", err)
	}
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("Unmarshal: unexpected error: %v", err)
	}
	if decoded != *siForked() {
		t.Fatalf("round trip = %+v, want %+v", decoded, *siForked())
	}
}
