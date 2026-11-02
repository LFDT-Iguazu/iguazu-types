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

//go:build ((linux && amd64) || (linux && arm64) || (darwin && amd64) || (darwin && arm64) || (windows && amd64)) && !blst_disabled

// Package bls_sig verifies BLS12-381 signatures in the minimal-pubkey-size
// variant (public keys in G1, signatures in G2) with the proof-of-possession
// ciphersuite, using the supranational/blst library.
package bls_sig

import (
	"errors"
	"fmt"

	"github.com/LFDT-Iguazu/iguazu-types/common"
	blst "github.com/supranational/blst/bindings/go"
)

// cipherSuite is the domain separation tag of the proof-of-possession
// ciphersuite (BLS_SIG_BLS12381G2_XMD:SHA-256_SSWU_RO_POP_).
var cipherSuite = []byte("BLS_SIG_BLS12381G2_XMD:SHA-256_SSWU_RO_POP_")

// PublicKey is a validated BLS public key.
type PublicKey struct {
	point *blst.P1Affine
}

// PublicKeyFromBytes decodes a compressed public key. The key must be on the
// curve, in the right subgroup and not the point at infinity.
func PublicKeyFromBytes(b []byte) (*PublicKey, error) {
	if len(b) != common.BlsPubKeyLength {
		return nil, fmt.Errorf("public key must be %d bytes", common.BlsPubKeyLength)
	}
	point := new(blst.P1Affine).Uncompress(b)
	if point == nil {
		return nil, errors.New("could not unmarshal bytes into public key")
	}
	// KeyValidate checks both the subgroup and the point at infinity.
	if !point.KeyValidate() {
		return nil, errors.New("ErrInfinitePubKey")
	}
	return &PublicKey{point: point}, nil
}

// Marshal returns the compressed encoding of the key.
func (k *PublicKey) Marshal() []byte {
	return k.point.Compress()
}

// Copy returns an independent copy of the key.
func (k *PublicKey) Copy() *PublicKey {
	point := *k.point
	return &PublicKey{point: &point}
}

// IsInfinite reports whether the key is the point at infinity.
func (k *PublicKey) IsInfinite() bool {
	return k.point.Equals(new(blst.P1Affine))
}

// Signature is a BLS signature that has passed the subgroup check.
type Signature struct {
	point *blst.P2Affine
}

// SignatureFromBytes decodes a compressed signature and checks that it is in
// the right subgroup. The point at infinity is accepted, as aggregated
// signatures may legitimately be infinite.
func SignatureFromBytes(b []byte) (*Signature, error) {
	if len(b) != common.BlsSigLength {
		return nil, fmt.Errorf("signature must be %d bytes", common.BlsSigLength)
	}
	point := new(blst.P2Affine).Uncompress(b)
	if point == nil {
		return nil, errors.New("could not unmarshal bytes into signature")
	}
	if !point.SigValidate(false) {
		return nil, errors.New("signature not in group")
	}
	return &Signature{point: point}, nil
}

// Verify reports whether the signature is valid for msg under pubKey. Both
// were already validated when they were decoded.
func (s *Signature) Verify(pubKey PublicKey, msg []byte) bool {
	return s.point.Verify(false, pubKey.point, false, msg, cipherSuite)
}

// Marshal returns the compressed encoding of the signature.
func (s *Signature) Marshal() []byte {
	return s.point.Compress()
}

// Copy returns an independent copy of the signature.
func (s *Signature) Copy() *Signature {
	point := *s.point
	return &Signature{point: &point}
}

// VerifyCompressed decodes a compressed signature and public key, validates
// both and reports whether the signature is valid for msg.
func VerifyCompressed(signature, pubKey, msg []byte) bool {
	return new(blst.P2Affine).VerifyCompressed(signature, true, pubKey, true, msg, cipherSuite)
}
