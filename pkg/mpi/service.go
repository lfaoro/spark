// Copyright (c) 2019 Leonardo Faoro. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mpi

type Service interface {
	Enrolled(data Request) Response
}

type Request struct {
	// ISO 4217 3-letter currency code.
	Currency string
	// Primary account number of card to charge.
	CardNumber string
	// Expiry month of card to charge.
	CardExpMonth string
	// Expiry year of card to charge.
	CardExpYear string
}

type Response struct {
	Enrolled bool
	// Electronic Commerce Indicator
	//
	// Values:
	// ECI Visa	ECI MC	Status
	// 5	2	Authentication Successful
	// 6	1	Attempts Processing Performed
	// 7	0	Authentication Failed
	// 7	1	Authentication Could Not Be Performed
	// 7	0	Error
	ECI int
	// Access Control Server URL
	ACS string
	// Payment Authentication Request
	PAR string
}
