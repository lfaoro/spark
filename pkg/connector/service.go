// Copyright (c) 2019 Leonardo Faoro. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package connector

// Connector is the interface all integrated processing banks must
// adhere to.
// todo: wip
type Connector interface {
	Process(token string, amount float32) error
	Chargeback(tx string) error
	Void(tx string) error
}
