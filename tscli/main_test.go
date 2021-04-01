package main

import (
	"testing"
)

func Test_timeFormatOK(t *testing.T) {
	tests := []struct {
		desc   string
		input  string
		expect bool
	}{
		{
			desc:   "correctly formatted current time",
			input:  "now",
			expect: true,
		},
		{
			desc:   "correctly formed hour",
			input:  "-1h",
			expect: true,
		},
		{
			desc:   "correctly formed minute",
			input:  "-1m",
			expect: true,
		},
		{
			desc:   "incorrectly formed, no unit",
			input:  "-1",
			expect: false,
		},
		{
			desc:   "incorrectly formed, no value",
			input:  "-m",
			expect: false,
		},
		{
			desc:   "incorrectly formed, no negative sign",
			input:  "1h",
			expect: false,
		},
		{
			desc:   "empty string",
			input:  "",
			expect: false,
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(tt *testing.T) {
			if got := timeFormatOK(test.input); got != test.expect {
				tt.Errorf("expected: %t, but got %t\n", test.expect, got)
			}
		})
	}
}
