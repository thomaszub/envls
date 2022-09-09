package format

import (
	"reflect"
	"testing"

	"github.com/thomaszub/envls/internal/env"
)

func TestDelimiterFormatter_Format(t *testing.T) {
	f := DelimiterFormatter{Delimiter: "->"}
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
	f := JsonFormatter{Pretty: false}
	tests := getTestCases()
	exp := `[{"Name":"SOMEVAR","Value":"someval"},{"Name":"OTHER","Value":"otherval"}]`

	act, err := f.Format(tests)

	if err != nil {
		t.Error(err.Error())
	}
	assert(t, exp, act)
}

func TestJsonFormatterPretty_Format(t *testing.T) {
	f := JsonFormatter{Pretty: true}
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

func TestJsonFormatterEmpty_Format(t *testing.T) {
	fPretty := JsonFormatter{Pretty: true}
	fCompact := JsonFormatter{Pretty: true}

	var tests []env.Var
	exp := "[]"

	actPretty, err := fPretty.Format(tests)
	if err != nil {
		t.Error(err.Error())
	}
	assert(t, exp, actPretty)

	actCompact, err := fCompact.Format(tests)
	if err != nil {
		t.Error(err.Error())
	}
	assert(t, exp, actCompact)
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
