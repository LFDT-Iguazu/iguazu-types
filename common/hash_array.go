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

package common

import (
	"sort"

	"golang.org/x/crypto/sha3"
)

// Ptr returns a pointer to a copy of h.
func (h Hash) Ptr() *Hash { return &h }

// HashArray represents a slice of Hashes.
type HashArray []Hash

// HashArrayFromBytes decodes the HashArray from byte representation.
func HashArrayFromBytes(data []byte) HashArray {
	count := len(data) / HashLength
	ha := make(HashArray, count)
	for i := 0; i < count; i++ {
		ha[i] = Hash{}
		ha[i].SetBytes(data[i*HashLength : (i+1)*HashLength])
	}
	return ha
}

// ToBytes encodes the HashArray to its byte representation.
func (ha HashArray) ToBytes() []byte {
	res := make([]byte, len(ha)*HashLength)
	for i := 0; i < len(ha); i++ {
		copy(res[i*HashLength:(i+1)*HashLength], ha[i][:])
	}
	return res
}

// Copy creates a copy of HashArray.
func (ha HashArray) Copy() HashArray {
	c := make(HashArray, len(ha))
	copy(c, ha)
	return c
}

// Intersection returns hashes present in both ha and hashArray.
func (ha HashArray) Intersection(hashArray HashArray) HashArray {
	c := make(HashArray, 0)
	m := make(map[Hash]bool)
	for _, item := range ha {
		m[item] = true
	}
	for _, item := range hashArray {
		if _, ok := m[item]; ok {
			c = append(c, item)
		}
	}
	return c
}

// SequenceIntersection returns the first contiguous sequence in ha that matches hashArray.
func (ha HashArray) SequenceIntersection(hashArray HashArray) HashArray {
	c := make(HashArray, 0)
	index := -1
	for _, item := range ha {
		i := hashArray.IndexOf(item)
		if index == -1 && i >= 0 {
			index = i + 1
			c = append(c, item)
			continue
		}
		if index >= 0 && index == i {
			index = i + 1
			c = append(c, item)
			continue
		}
		if index >= 0 && index != i {
			break
		}
	}
	return c
}

// Difference returns hashes in ha that are not in hashArray.
func (ha HashArray) Difference(hashArray HashArray) HashArray {
	mb := make(map[Hash]struct{})
	for _, x := range hashArray {
		mb[x] = struct{}{}
	}
	diff := []Hash{}
	for _, x := range ha {
		if _, found := mb[x]; !found {
			diff = append(diff, x)
		}
	}
	return diff
}

// Uniq returns a new HashArray with duplicate hashes removed.
func (ha HashArray) Uniq() HashArray {
	c := make(HashArray, 0)
	m := make(map[Hash]struct{})
	for _, item := range ha {
		if _, ok := m[item]; !ok {
			c = append(c, item)
			m[item] = struct{}{}
		}
	}
	return c
}

// Deduplicate removes duplicate hashes from ha in-place.
func (ha *HashArray) Deduplicate() {
	m := make(map[Hash]struct{})
	j := 0
	for _, item := range *ha {
		if _, ok := m[item]; !ok {
			(*ha)[j] = item
			j++
			m[item] = struct{}{}
		}
	}
	*ha = (*ha)[:j]
}

// IsUniq returns true if all hashes in ha are unique.
func (ha HashArray) IsUniq() bool {
	return len(ha) == len(ha.Copy().Uniq())
}

// Concat appends hashes to ha and returns the result.
func (ha HashArray) Concat(hashes HashArray) HashArray {
	ha = append(ha, hashes...)
	return ha
}

// IsEqualTo returns true if ha and hashArray have the same hashes in the same order.
func (ha HashArray) IsEqualTo(hashArray HashArray) bool {
	if len(ha) != len(hashArray) {
		return false
	}
	for i, v := range ha {
		if v != hashArray[i] {
			return false
		}
	}
	return true
}

// Sort sorts ha lexicographically and returns it.
func (ha HashArray) Sort() HashArray {
	strs := make([]string, len(ha))
	for i, hash := range ha {
		strs[i] = hash.String()
	}
	sort.Strings(strs)
	for i, str := range strs {
		ha[i] = HexToHash(str)
	}
	return ha
}

// Has returns true if ha contains hash.
func (ha HashArray) Has(hash Hash) bool {
	for _, h := range ha {
		if h == hash {
			return true
		}
	}
	return false
}

// IndexOf returns the index of hash in ha, or -1 if not found.
func (ha HashArray) IndexOf(hash Hash) int {
	for i, h := range ha {
		if h == hash {
			return i
		}
	}
	return -1
}

// Reverse returns a reversed copy of ha.
func (ha HashArray) Reverse() HashArray {
	cpy := ha.Copy()
	i := 0
	j := len(ha) - 1
	for i < j {
		cpy[i], cpy[j] = cpy[j], cpy[i]
		i++
		j--
	}
	return cpy
}

// Key computes a deterministic hash key from the unique, sorted hashes in ha.
func (ha HashArray) Key() Hash {
	c := ha.Copy()
	buf := c.Uniq().Sort().ToBytes()
	sha := sha3.NewLegacyKeccak256()
	sha.Write(buf[:])
	key := sha.Sum(nil)
	return BytesToHash(key)
}

// Contains reports whether desired is in addresses.
func Contains(addresses []Address, desired Address) (bool, int) {
	for i, a := range addresses {
		if a == desired {
			return true, i
		}
	}
	return false, -1
}

// SortAddresses sorts addresses lexicographically.
func SortAddresses(addresses []Address) []Address {
	res := make([]Address, len(addresses))
	strs := make([]string, len(addresses))
	for i, a := range addresses {
		strs[i] = a.String()
	}
	sort.Strings(strs)
	for i, str := range strs {
		res[i] = HexToAddress(str)
	}
	return res
}
