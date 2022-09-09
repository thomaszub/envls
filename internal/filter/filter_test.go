package filter

import (
	"regexp"
	"testing"

	"github.com/thomaszub/envls/internal/env"
)

func TestNoPrefixFilter_Accept(t *testing.T) {
	f := NoPrefixFilter{Prefix: "_"}
	tests := []struct {
		input env.Var
		exp   bool
	}{
		{
			input: env.Var{
				Name:  "_SOMEVAR",
				Value: "value",
			},
			exp: false,
		},
		{
			input: env.Var{
				Name:  "SOMEVAR",
				Value: "value",
			},
			exp: true,
		}, {
			input: env.Var{
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
		input env.Var
		exp   bool
	}{
		{
			input: env.Var{
				Name:  "_SOMEWHATVAR",
				Value: "value",
			},
			exp: true,
		},
		{
			input: env.Var{
				Name:  "SOMEVAR",
				Value: "valWHATue",
			},
			exp: true,
		}, {
			input: env.Var{
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

func TestEmptyAndFilter_Accept(t *testing.T) {
	f := AndFilter{Filters: []Filter{}}
	tests := []struct {
		input env.Var
		exp   bool
	}{
		{
			input: env.Var{
				Name:  "_&%$abc",
				Value: "&value",
			},
			exp: true,
		},
	}

	for _, tt := range tests {
		act := f.Accept(tt.input)
		if act != tt.exp {
			t.Errorf("Empty filter handler must accept all. Test: %+v", tt.input)
		}
	}
}

func TestAndFilter_Accept(t *testing.T) {
	h := AndFilter{Filters: []Filter{&NoPrefixFilter{Prefix: "_"}, &NoPrefixFilter{Prefix: "&"}}}
	tests := []struct {
		input env.Var
		exp   bool
	}{
		{
			input: env.Var{
				Name:  "SOMEVAR",
				Value: "&value",
			},
			exp: true,
		},
		{
			input: env.Var{
				Name:  "&SOMEVAR",
				Value: "value",
			},
			exp: false,
		},
		{
			input: env.Var{
				Name:  "_SOMEVAR",
				Value: "value",
			},
			exp: false,
		},
		{
			input: env.Var{
				Name:  "SOMEVAR",
				Value: "_value",
			},
			exp: true,
		},
	}

	for _, tt := range tests {
		act := h.Accept(tt.input)
		if act != tt.exp {
			t.Errorf("Wrong Accept result for FilterHandler. Got %t for %+v", act, tt.input)
		}
	}
}
