package format

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/thomaszub/envls/internal/env"
)

type Formatter interface {
	Format(envs []env.Var) (string, error)
}

type DelimiterFormatter struct {
	delimiter string
}

func (d *DelimiterFormatter) Format(envs []env.Var) (string, error) {
	var output []string
	for _, env := range envs {
		output = append(output, env.Name+d.delimiter+env.Value)
	}
	return strings.Join(output, "\n"), nil
}

type JsonFormatter struct {
	pretty bool
}

func (d *JsonFormatter) Format(envs []env.Var) (string, error) {
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

func GetFormatter(flag string) (Formatter, error) {
	splitted := strings.Split(flag, ",")
	switch splitted[0] {
	case "del":
		if len(splitted) != 2 {
			return nil, fmt.Errorf("formatter \"del\" takes exactly one configuration argument, e.g. del,=. Got: %s", flag)
		}
		return &DelimiterFormatter{delimiter: splitted[1]}, nil
	case "json":
		if len(splitted) != 2 {
			return nil, jsonFormatterError(flag)
		}
		switch splitted[1] {
		case "compact":
			return &JsonFormatter{pretty: false}, nil
		case "pretty":
			return &JsonFormatter{pretty: true}, nil
		default:
			return nil, jsonFormatterError(flag)
		}
	default:
		return nil, fmt.Errorf("unknown formatter: %s", flag)
	}
}

func jsonFormatterError(flag string) error {
	return fmt.Errorf("formatter \"json\" takes exactly one configuration argument of compact or pretty, e.g. json,compact. Got: %s", flag)
}
