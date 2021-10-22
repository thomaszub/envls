package internal

import "testing"

func TestNoPrefixFilter_Accept(t *testing.T) {
	f := NewNoPrefixFilter("_")
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
