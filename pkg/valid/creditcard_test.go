// Copyright (c) 2019 Leonardo Faoro. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package valid

import "testing"

func TestIsCreditCard(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		number string
		want   bool
	}{
		{"empty", "", false},
		{"not numbers", "credit card", false},

		{"visa", "4220855426222389", true},
		{"visa spaces", "4220 8554 2622 2389", true},
		{"visa dashes", "4220-8554-2622-2389", true},
		{"mastercard", "5139288802098206", true},
		{"american express", "374953669708156", true},
		{"discover", "6011464355444102", true},
		{"jcb", "3548209662790989", true},

		// below should be valid, do they respect international standards?
		{"diners club international", "30190239451016", false},
		{"rupay", "6521674451993089", false},
		{"mir", "2204151414444676", false},
		{"china unionPay", "624356436327468104", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsCreditCard(tt.number); got != tt.want {
				t.Errorf("IsCreditCard(%v) = %v, want %v", tt.number, got, tt.want)
			}
		})
	}
}
