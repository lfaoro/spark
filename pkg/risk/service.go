// Copyright (c) 2019 Leonardo Faoro. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package risk

type Service interface {
	Calculate(data Data)
	Add(r Rule)
	Remove(name string)
}

type Rule struct {
	Name   string
	Weight float32
	Fn     func(data Data) bool
}

type Data struct{}
