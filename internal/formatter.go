package internal

import (
	"encoding/json"
	"fmt"
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
