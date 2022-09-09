package format

import (
	"encoding/json"
	"strings"

	"github.com/thomaszub/envls/internal/env"
)

type DelimiterFormatter struct {
	Delimiter string
}

func (d *DelimiterFormatter) Format(envs []env.Var) (string, error) {
	var output []string
	for _, env := range envs {
		output = append(output, env.Name+d.Delimiter+env.Value)
	}
	return strings.Join(output, "\n"), nil
}

type JsonFormatter struct {
	Pretty bool
}

func (d *JsonFormatter) Format(envs []env.Var) (string, error) {
	if len(envs) == 0 {
		return "[]", nil
	}
	var j []byte
	var err error
	if d.Pretty {
		j, err = json.MarshalIndent(envs, "", "    ")
	} else {
		j, err = json.Marshal(envs)
	}
	if err != nil {
		return "", err
	}
	return string(j), nil
}
