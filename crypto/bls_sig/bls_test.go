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

package bls_sig_test

import (
	"bytes"
	"testing"

	"github.com/LFDT-Iguazu/iguazu-types/common"
	"github.com/LFDT-Iguazu/iguazu-types/crypto/bls_sig"
	blst "github.com/supranational/blst/bindings/go"
)

var dst = []byte("BLS_SIG_BLS12381G2_XMD:SHA-256_SSWU_RO_POP_")

// signed returns a compressed public key and a signature of msg made with a
// key derived from seed.
func signed(t *testing.T, seed byte, msg []byte) (pub, sig []byte) {
	t.Helper()
	ikm := bytes.Repeat([]byte{seed}, 32)
	sk := blst.KeyGen(ikm)
	pub = new(blst.P1Affine).From(sk).Compress()
	sig = new(blst.P2Affine).Sign(sk, msg, dst).Compress()
	return pub, sig
}

func TestSignatureRoundTrip(t *testing.T) {
	msg := []byte("iguazu")
	pubBin, sigBin := signed(t, 1, msg)

	pub, err := bls_sig.PublicKeyFromBytes(pubBin)
	if err != nil {
		t.Fatal(err)
	}
	sig, err := bls_sig.SignatureFromBytes(sigBin)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(pub.Marshal(), pubBin) || !bytes.Equal(sig.Marshal(), sigBin) {
		t.Fatal("Marshal does not return the decoded bytes")
	}
	if !sig.Verify(*pub, msg) {
		t.Error("Verify: valid signature rejected")
	}
	if sig.Verify(*pub, []byte("other")) {
		t.Error("Verify: signature accepted for another message")
	}
	if !bls_sig.VerifyCompressed(sigBin, pubBin, msg) {
		t.Error("VerifyCompressed: valid signature rejected")
	}

	otherPub, _ := signed(t, 2, msg)
	if bls_sig.VerifyCompressed(sigBin, otherPub, msg) {
		t.Error("VerifyCompressed: signature accepted for another key")
	}
	if pub.IsInfinite() {
		t.Error("IsInfinite: regular key reported as infinite")
	}
}

func TestKnownVectors(t *testing.T) {
	// Vectors from the Waterfall network test suite.
	cases := []struct {
		pub, msg, sig string
	}{
		{
			pub: "0xa7914b432b6795c5fab4a860dfa802f98eeaf8e3885c29309572b2f25fdfd8a5f65330454371bf36e4323807bdcf9adf",
			msg: "0x68656c6c6f", // "hello"
			sig: "0x92e59bab188211166b6aef639320166a0ead5451dc4ecbeda244fcf9c146552d46bfdccdb15cc4651c99c0b671f42c9300bee3fe9c761dd34a7a9d0f21a20c5572d5c0d8cf6492c02ea6d59fc97ac6256eebd61fbe98d760e104354674c803d1",
		},
		{
			pub: "0xa491d1b0ecd9bb917989f0e74f0dea0422eac4a873e5e2644f368dffb9a6e20fd6e10c1b77654d067c0618f6e5a7f79a",
			msg: "0x0000000000000000000000000000000000000000000000000000000000000000",
			sig: "0xb6ed936746e01f8ecf281f020953fbf1f01debd5657c4a383940b020b26507f6076334f91e2366c96e9ab279fb5158090352ea1c5b0c9274504f4f0e7053af24802e51e4568d164fe986834f41e55c8e850ce1f98458c0cfc9ab380b55285a55",
		},
	}
	for _, c := range cases {
		pubBin, msg, sigBin := common.FromHex(c.pub), common.FromHex(c.msg), common.FromHex(c.sig)
		pub, err := bls_sig.PublicKeyFromBytes(pubBin)
		if err != nil {
			t.Fatal(err)
		}
		sig, err := bls_sig.SignatureFromBytes(sigBin)
		if err != nil {
			t.Fatal(err)
		}
		if !sig.Verify(*pub, msg) || !bls_sig.VerifyCompressed(sigBin, pubBin, msg) {
			t.Errorf("vector %s did not verify", c.pub)
		}
	}
}

func TestDecodeErrors(t *testing.T) {
	pubCases := map[string][]byte{
		"public key must be 48 bytes":               make([]byte, 47),
		"could not unmarshal bytes into public key": make([]byte, 48),
	}
	for want, in := range pubCases {
		if _, err := bls_sig.PublicKeyFromBytes(in); err == nil || err.Error() != want {
			t.Errorf("PublicKeyFromBytes(%d bytes) error = %v, want %q", len(in), err, want)
		}
	}
	// The compressed point at infinity decodes but is not a valid key.
	infinity := append([]byte{0xc0}, make([]byte, 47)...)
	if _, err := bls_sig.PublicKeyFromBytes(infinity); err == nil {
		t.Error("PublicKeyFromBytes accepted the point at infinity")
	}

	sigCases := map[string][]byte{
		"signature must be 96 bytes":               make([]byte, 95),
		"could not unmarshal bytes into signature": make([]byte, 96),
	}
	for want, in := range sigCases {
		if _, err := bls_sig.SignatureFromBytes(in); err == nil || err.Error() != want {
			t.Errorf("SignatureFromBytes(%d bytes) error = %v, want %q", len(in), err, want)
		}
	}
	if bls_sig.VerifyCompressed(make([]byte, 96), make([]byte, 48), []byte("x")) {
		t.Error("VerifyCompressed accepted invalid encodings")
	}
}

func TestCopy(t *testing.T) {
	pubBin, sigBin := signed(t, 3, []byte("copy"))
	pub, _ := bls_sig.PublicKeyFromBytes(pubBin)
	sig, _ := bls_sig.SignatureFromBytes(sigBin)
	if pc := pub.Copy(); pc == pub || !bytes.Equal(pc.Marshal(), pub.Marshal()) {
		t.Error("PublicKey.Copy is not an equal, independent copy")
	}
	if sc := sig.Copy(); sc == sig || !bytes.Equal(sc.Marshal(), sig.Marshal()) {
		t.Error("Signature.Copy is not an equal, independent copy")
	}
}
