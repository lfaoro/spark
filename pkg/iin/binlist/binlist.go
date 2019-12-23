// Copyright (c) 2019 Leonardo Faoro. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package binlist

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"path"
	"time"

	"github.com/pkg/errors"

	"github.com/lfaoro/spark/pkg/iin"
)

// BinList is the Service implementation for the https://binlist.net/ service.
type BinList struct {
	endpoint string
	key      string
}

func New(key string) iin.Service {
	return BinList{
		endpoint: "https://lookup.binlist.net",
		key:      key,
	}
}

func (b BinList) Lookup(number string) (*iin.CardMetadata, error) {
	if len(number) > 8 {
		number = number[:7]
	}
	if len(number) < 6 {
		return nil, errors.Errorf("invalid number: %v", number)
	}

	c := http.Client{Timeout: 1 * time.Second}
	endpoint, _ := url.Parse(b.endpoint)
	endpoint.Path = path.Join(endpoint.Path, number)
	req, err := http.NewRequest(http.MethodGet, endpoint.String(), nil)
	if err != nil {
		return nil, errors.Wrap(err, "failed binlist request")
	}
	req.Header.Add("Accept-Version", "3")
	// req.Header.Add("key", b.key)

	res, err := c.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "failed binlist request")
	}
	defer res.Body.Close()
	if res.StatusCode > 299 {
		return nil, errors.Errorf("failed binlist request with code %v", res.StatusCode)
	}

	result := Result{}
	err = json.NewDecoder(res.Body).Decode(&result)
	if err != nil {
		return nil, errors.Wrap(err, "failed binlist decoding")
	}

	return result.convert(), nil
}

type Result struct {
	Number struct {
		Length int  `json:"length"`
		Luhn   bool `json:"luhn"`
	} `json:"number"`
	Scheme  string `json:"scheme"`
	Type    string `json:"type"`
	Brand   string `json:"brand"`
	Prepaid bool   `json:"prepaid"`
	Country struct {
		Numeric   string `json:"numeric"`
		Alpha2    string `json:"alpha2"`
		Name      string `json:"name"`
		Emoji     string `json:"emoji"`
		Currency  string `json:"currency"`
		Latitude  int    `json:"latitude"`
		Longitude int    `json:"longitude"`
	} `json:"country"`
	Bank struct {
		Name  string `json:"name"`
		URL   string `json:"url"`
		Phone string `json:"phone"`
		City  string `json:"city"`
	} `json:"bank"`
}

func (r Result) convert() *iin.CardMetadata {
	mapsURL := fmt.Sprintf(
		"https://www.google.com/maps/search/?api=1&query=%d,%d",
		r.Country.Latitude, r.Country.Longitude)
	return &iin.CardMetadata{
		Issuer:            r.Bank.Name,
		IssuerURL:         r.Bank.URL,
		IssuerPhone:       r.Bank.Phone,
		IssuerCity:        r.Bank.City,
		IssuerCountry:     r.Country.Name,
		IssuerCountryCode: r.Country.Alpha2,
		CardScheme:        r.Scheme,
		CardBrand:         r.Brand,
		CardType:          r.Type,
		CardIsPrepaid:     r.Prepaid,
		CardCurrency:      r.Country.Currency,
		IssuerLatitude:    float32(r.Country.Latitude),
		IssuerLongitude:   float32(r.Country.Longitude),
		IssuerMap:         mapsURL,
	}
}
