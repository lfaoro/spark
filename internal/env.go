// Copyright (c) 2019 Leonardo Faoro. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package internal

import (
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

func LoadEnv() {
	wd, _ := os.Getwd()
	err := godotenv.Load(filepath.Join(wd, ".env"))
	if err != nil {
		log.Printf("failed loading environment variables: %v\n", err)
	}
}
