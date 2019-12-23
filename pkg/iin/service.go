// Copyright (c) 2019 Leonardo Faoro. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package iin provides a common interface for Issuer Identification Number
// services.
// ref: https://en.wikipedia.org/wiki/Payment_card_number#Issuer_identification_number_.28IIN.29
package iin

type CardMetadata struct {
	Issuer            string
	IssuerURL         string
	IssuerPhone       string
	IssuerCity        string
	IssuerCountry     string
	IssuerCountryCode string // alpha2

	CardScheme    string // visa/mastercard
	CardBrand     string // Visa/Platinum
	CardType      string // debit/credit
	CardIsPrepaid bool   // false
	CardCurrency  string // DKK

	// GeoMapping
	IssuerLatitude  float32
	IssuerLongitude float32
	IssuerMap       string
}

type Service interface {
	Lookup(s string) (*CardMetadata, error)
}
