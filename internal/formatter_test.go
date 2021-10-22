package internal

import (
	"reflect"
	"testing"
)

func TestDelimiterFormatter_Format(t *testing.T) {
	f := NewDelimiterFormatter("->")
	tests := []EnvVar{
		{
			Name:  "SOMEVAR",
			Value: "value",
		},
		{
			Name:  "OTHER",
			Value: "otherval",
		},
	}

	exp := []string{
		"SOMEVAR->value",
		"OTHER->otherval",
	}

	act := f.Format(tests)

	if !reflect.DeepEqual(exp, act) {
		t.Errorf("Formatting not correct. Got %s", act)
	}

}
