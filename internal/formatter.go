package internal

import (
	"encoding/json"
	"strings"
)

type Formatter interface {
	Format(envs []EnvVar) (string, error)
}

type DelimiterFormatter struct {
	delimiter string
}

func (d *DelimiterFormatter) Format(envs []EnvVar) (string, error) {
	output := make([]string, 0)
	for _, env := range envs {
		output = append(output, env.Name+d.delimiter+env.Value)
	}
	return strings.Join(output, "\n"), nil
}

func NewDelimiterFormatter(delimiter string) Formatter {
	return &DelimiterFormatter{delimiter: delimiter}
}

type JsonFormatter struct {
	pretty bool
}

func (d *JsonFormatter) Format(envs []EnvVar) (string, error) {
	var j []byte
	var err error
	if d.pretty {
		j, err = json.MarshalIndent(envs, "", "    ")
	} else {
		j, err = json.Marshal(envs)
	}
	if err != nil {
		return "", err
	}
	return string(j), nil
}

func NewJsonFormatter(pretty bool) Formatter {
	return &JsonFormatter{pretty: pretty}
}
