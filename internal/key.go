// Copyright (c) 2019 Leonardo Faoro. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package internal

import (
	"strings"

	"github.com/google/uuid"
)

// GenKey generates a new API key.
func GenKey(isTest bool) string {
	prefix := "key_live_"
	if isTest {
		prefix += "key_test_"
	}
	uid := uuid.New().String()
	uid = strings.Replace(uid, "-", "", -1)
	pwd := prefix + uid
	return pwd[:32]
}
