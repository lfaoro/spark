// Copyright (c) 2019 Leonardo Faoro. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gcpkms

import (
	"testing"
)

var hsm = New("txexplorer", "europe-west4", Debug(true))

func Test_GCPKMS_Encrypt(t *testing.T) {
	// every encryption run, will generate a different set of bytes
	// we can test the correctness of the encryption only by decrypting the payload
	// to make this test repeatable, we check only that the bytes length is always
	// the same as long as the same encryption algorithm is used.
	var want = []byte("\n$\x00s\xcfG\x8e\x9a\x10F]6\x92\x14,kZ\xacnZ:\xa7\xc5\xf0[\a\xec%\x95~\x86\xd8\"\x1ad\xe4\x93o\x124\x122\n\f&\x9d\xe3\xf1`\x8dm\xfdr\xe4\x04\\\x12\x10$\xb836\x1d==\x13\xe8Ԟ4\x87\xc9~\xd9\x1a\x10ű\x97\x12a\a\b\x0e\x05\x8cZF\xea\xe1\xa6^")
	pt := "leonardo"
	ct, err := hsm.Encrypt([]byte(pt))
	if err != nil {
		t.Error(err)
	}
	if len(ct) != len(want) {
		t.Errorf("got %d, want %d", len(ct), len(want))
	}
}

func Test_GCPKMS_Decrypt(t *testing.T) {
	const want = "leonardo"
	var ct = []byte("\n$\x00s\xcfG\x8e\x9a\x10F]6\x92\x14,kZ\xacnZ:\xa7\xc5\xf0[\a\xec%\x95~\x86\xd8\"\x1ad\xe4\x93o\x124\x122\n\f&\x9d\xe3\xf1`\x8dm\xfdr\xe4\x04\\\x12\x10$\xb836\x1d==\x13\xe8Ԟ4\x87\xc9~\xd9\x1a\x10ű\x97\x12a\a\b\x0e\x05\x8cZF\xea\xe1\xa6^")
	pt, err := hsm.Decrypt(ct)
	if err != nil {
		t.Error(err)
	}
	if string(pt) != want {
		t.Errorf("got %s, want %s", pt, want)
	}
}

// BenchmarkGCPKMS_Encrypt-8   	       5	 217883063 ns/op
func BenchmarkGCPKMS_Encrypt(b *testing.B) {
	t := &testing.T{}
	for n := 0; n < b.N; n++ {
		Test_GCPKMS_Encrypt(t)
	}
}

// BenchmarkGCPKMS_Decrypt-8   	       8	 165241458 ns/op
func BenchmarkGCPKMS_Decrypt(b *testing.B) {
	t := &testing.T{}
	for n := 0; n < b.N; n++ {
		Test_GCPKMS_Decrypt(t)
	}
}
