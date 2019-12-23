// Copyright (c) 2019 Leonardo Faoro. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package valid

import (
	"regexp"
	"strconv"
)

// IsCreditCard applies regular expression and luhn algorithm validation
// to validate the correct structure of a payment card number.
//
// inspired by Alex Saskevich `govalidator` package:
// ref: https://github.com/asaskevich/govalidator
func IsCreditCard(str string) bool {
	var CreditCard string = "^(?:4[0-9]{12}(?:[0-9]{3})?|5[1-5][0-9]{14}|6(?:011|5[0-9][0-9])[0-9]{12}|3[47][0-9]{13}|3(?:0[0-5]|[68][0-9])[0-9]{11}|(?:2131|1800|35\\d{3})\\d{11})$"
	var notNumberRegexp = regexp.MustCompile("[^0-9]+")
	var rxCreditCard = regexp.MustCompile(CreditCard)

	sanitized := notNumberRegexp.ReplaceAllString(str, "")
	if !rxCreditCard.MatchString(sanitized) {
		return false
	}

	return Luhn(sanitized)
}

// Luhn implements the Luhn checksum formula that validates
// identification numbers. It was designed to protect against accidental
// errors, not malicious attacks.
//
// https://en.wikipedia.org/wiki/Luhn_algorithm
func Luhn(s string) bool {
	var sum int64
	var digit string
	var tmpNum int64
	var shouldDouble bool
	for i := len(s) - 1; i >= 0; i-- {
		digit = s[i:(i + 1)]
		num, _ := strconv.Atoi(digit)
		tmpNum = int64(num)
		if shouldDouble {
			tmpNum *= 2
			if tmpNum >= 10 {
				sum += (tmpNum % 10) + 1
			} else {
				sum += tmpNum
			}
		} else {
			sum += tmpNum
		}
		shouldDouble = !shouldDouble
	}
	return sum%10 == 0
}
