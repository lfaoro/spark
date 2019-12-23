// Copyright (c) 2019 Leonardo Faoro. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package aesgcm

import (
	"testing"
)

var hsm, _ = New("oUnmr6cm9G7yeRBMsZrqk84w8M22jcxM")

func Test_AESGCM_Encrypt(t *testing.T) {
	var want = []byte("<\xb9\xdd\x01_\x03\x1aƁڝ\xc2\xf1\x91b\x9c\x967\xec\xe0\x03\x94}O5}\x96$\xc1}\xb9\xae\xec\x8d\xf1|")
	pt := "leonardo"
	ct, err := hsm.Encrypt([]byte(pt))
	if err != nil {
		t.Error(err)
	}
	if len(ct) != len(want) {
		t.Errorf("got %d, want %d", len(ct), len(want))
	}
}

func Test_AESGCM_Decrypt(t *testing.T) {
	const want = "leonardo"
	var ct = []byte("<\xb9\xdd\x01_\x03\x1aƁڝ\xc2\xf1\x91b\x9c\x967\xec\xe0\x03\x94}O5}\x96$\xc1}\xb9\xae\xec\x8d\xf1|")
	pt, err := hsm.Decrypt(ct)
	if err != nil {
		t.Error(err)
	}
	if string(pt) != want {
		t.Errorf("got %s, want %s", pt, want)
	}
}

func BenchmarkGCPKMS_Encrypt(b *testing.B) {
	t := &testing.T{}
	for n := 0; n < b.N; n++ {
		Test_AESGCM_Encrypt(t)
	}
}

func BenchmarkGCPKMS_Decrypt(b *testing.B) {
	t := &testing.T{}
	for n := 0; n < b.N; n++ {
		Test_AESGCM_Decrypt(t)
	}
}
