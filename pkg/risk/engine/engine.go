// Copyright (c) 2019 Leonardo Faoro. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package engine

import (
	"sync"

	"github.com/lfaoro/spark/pkg/risk"
)

type Engine struct {
	score float32

	m     sync.Mutex // protects below
	rules map[string]risk.Rule
}

func New() risk.Service {
	return &Engine{}
}

func NewWithRules() risk.Service {
	return &Engine{}
}

func (e *Engine) Add(r risk.Rule) {
	e.m.Lock()
	e.rules[r.Name] = r
	e.m.Unlock()
}
func (e *Engine) Remove(name string) {
	e.m.Lock()
	delete(e.rules, name)
	e.m.Unlock()
}

func (e *Engine) Calculate(rd risk.Data) {
	for _, r := range e.rules {
		if r.Fn(rd) {
			e.score += r.Weight
		}
	}
}

func test() {
	r1 := risk.Rule{
		Name:   "same_country",
		Weight: 20,
		Fn: func(rd risk.Data) bool {
			return false
		},
	}
	r := New()
	r.Add(r1)
}
