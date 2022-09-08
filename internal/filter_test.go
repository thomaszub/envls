package internal

import (
	"regexp"
	"testing"
)

func TestNoPrefixFilter_Accept(t *testing.T) {
	f := NoPrefixFilter{Prefix: "_"}
	tests := []struct {
		input EnvVar
		exp   bool
	}{
		{
			input: EnvVar{
				Name:  "_SOMEVAR",
				Value: "value",
			},
			exp: false,
		},
		{
			input: EnvVar{
				Name:  "SOMEVAR",
				Value: "value",
			},
			exp: true,
		}, {
			input: EnvVar{
				Name:  "SOMEVAR",
				Value: "_value",
			},
			exp: true,
		},
	}

	for _, tt := range tests {
		act := f.Accept(tt.input)
		if act != tt.exp {
			t.Errorf("Wrong Accept result. Got %t for %+v", act, tt.input)
		}
	}
}

func TestRegexFilter_Accept(t *testing.T) {
	regex, err := regexp.Compile(".*WHAT.*")
	if err != nil {
		t.Fatal(err)
	}
	f := RegexFilter{Regex: regex}
	tests := []struct {
		input EnvVar
		exp   bool
	}{
		{
			input: EnvVar{
				Name:  "_SOMEWHATVAR",
				Value: "value",
			},
			exp: true,
		},
		{
			input: EnvVar{
				Name:  "SOMEVAR",
				Value: "valWHATue",
			},
			exp: true,
		}, {
			input: EnvVar{
				Name:  "SOMEVAR",
				Value: "_value",
			},
			exp: false,
		},
	}

	for _, tt := range tests {
		act := f.Accept(tt.input)
		if act != tt.exp {
			t.Errorf("Wrong Accept result. Got %t for %+v", act, tt.input)
		}
	}
}
