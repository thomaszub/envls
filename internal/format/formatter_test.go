package format

import (
	"reflect"
	"testing"

	"github.com/thomaszub/envls/internal/env"
)

func TestDelimiterFormatter_Format(t *testing.T) {
	f := DelimiterFormatter{delimiter: "->"}
	tests := getTestCases()
	exp := `SOMEVAR->someval
OTHER->otherval`

	act, err := f.Format(tests)

	if err != nil {
		t.Error(err.Error())
	}
	assert(t, exp, act)
}

func TestJsonFormatterCompact_Format(t *testing.T) {
	f := JsonFormatter{pretty: false}
	tests := getTestCases()
	exp := `[{"Name":"SOMEVAR","Value":"someval"},{"Name":"OTHER","Value":"otherval"}]`

	act, err := f.Format(tests)

	if err != nil {
		t.Error(err.Error())
	}
	assert(t, exp, act)
}

func TestJsonFormatterPretty_Format(t *testing.T) {
	f := JsonFormatter{pretty: true}
	tests := getTestCases()
	exp := `[
    {
        "Name": "SOMEVAR",
        "Value": "someval"
    },
    {
        "Name": "OTHER",
        "Value": "otherval"
    }
]`

	act, err := f.Format(tests)

	if err != nil {
		t.Error(err.Error())
	}
	assert(t, exp, act)
}

func getTestCases() []env.Var {
	return []env.Var{
		{
			Name:  "SOMEVAR",
			Value: "someval",
		},
		{
			Name:  "OTHER",
			Value: "otherval",
		},
	}
}

func assert(t *testing.T, exp string, act string) {
	if !reflect.DeepEqual(exp, act) {
		t.Errorf("Formatting not correct.\nExpected:\n%s\nGot:\n%s", exp, act)
	}
}
