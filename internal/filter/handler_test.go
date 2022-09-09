package filter

import (
	"reflect"
	"testing"

	"github.com/thomaszub/envls/internal/env"
)

func TestEmptyFilterHandler_Accept(t *testing.T) {
	f := FilterHandler{Filters: []Filter{}}
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
		act := f.accept(tt.input)
		if act != tt.exp {
			t.Errorf("Empty filter handler must accept all. Test: %+v", tt.input)
		}
	}
}

func TestFilterHandler_Accept(t *testing.T) {
	h := FilterHandler{Filters: []Filter{&NoPrefixFilter{Prefix: "_"}, &NoPrefixFilter{Prefix: "&"}}}
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
		act := h.accept(tt.input)
		if act != tt.exp {
			t.Errorf("Wrong Accept result for FilterHandler. Got %t for %+v", act, tt.input)
		}
	}
}

func TestFilterHandler_Accepted(t *testing.T) {
	h := FilterHandler{Filters: []Filter{&NoPrefixFilter{Prefix: "_"}, &NoPrefixFilter{Prefix: "&"}}}
	tests := []env.Var{
		{
			Name:  "SOMEVAR",
			Value: "&value",
		},
		{
			Name:  "&SOMEVAR",
			Value: "value",
		},
		{
			Name:  "_SOMEVAR",
			Value: "value",
		},
		{
			Name:  "SOMEVAR",
			Value: "_value",
		},
	}

	exp := []env.Var{
		{
			Name:  "SOMEVAR",
			Value: "&value",
		},
		{
			Name:  "SOMEVAR",
			Value: "_value",
		},
	}

	act := h.Accepted(tests)

	if !reflect.DeepEqual(exp, act) {
		t.Errorf("Accepted list contains not the excpected elements. Got: %+v", act)
	}
}
