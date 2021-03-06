package internal

import (
	"reflect"
	"testing"
)

func TestEmptyFilterHandler_Accept(t *testing.T) {
	f := NewEmptyFilterHandler()
	tests := []struct {
		input EnvVar
		exp   bool
	}{
		{
			input: EnvVar{
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

func TestFilterHandler_Accept(t *testing.T) {
	f := NewEmptyFilterHandler()
	f.AppendFilter(NewNoPrefixFilter("_"))
	f.AppendFilter(NewNoPrefixFilter("&"))
	tests := []struct {
		input EnvVar
		exp   bool
	}{
		{
			input: EnvVar{
				Name:  "SOMEVAR",
				Value: "&value",
			},
			exp: true,
		},
		{
			input: EnvVar{
				Name:  "&SOMEVAR",
				Value: "value",
			},
			exp: false,
		},
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
				Value: "_value",
			},
			exp: true,
		},
	}

	for _, tt := range tests {
		act := f.Accept(tt.input)
		if act != tt.exp {
			t.Errorf("Wrong Accept result for FilterHandler. Got %t for %+v", act, tt.input)
		}
	}
}

func TestFilterHandler_Accepted(t *testing.T) {
	f := NewEmptyFilterHandler()
	f.AppendFilter(NewNoPrefixFilter("_"))
	f.AppendFilter(NewNoPrefixFilter("&"))
	tests := []EnvVar{
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

	exp := []EnvVar{
		{
			Name:  "SOMEVAR",
			Value: "&value",
		},
		{
			Name:  "SOMEVAR",
			Value: "_value",
		},
	}

	act := f.Accepted(tests)

	if !reflect.DeepEqual(exp, act) {
		t.Errorf("Accepted list contains not the excpected elements. Got: %+v", act)
	}
}
