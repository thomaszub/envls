package internal

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/thomaszub/envls/internal/env"
	"github.com/thomaszub/envls/internal/filter"
	"github.com/thomaszub/envls/internal/format"
)

type EnvReader interface {
	Read() []env.Var
}

type Formatter interface {
	Format(envs []env.Var) (string, error)
}

type Pipeline struct {
	reader    EnvReader
	filter    filter.Filter
	formatter Formatter
}

func (p *Pipeline) Execute() (string, error) {
	envs := p.reader.Read()
	var accepted []env.Var
	for _, env := range envs {
		if p.filter.Accept(env) {
			accepted = append(accepted, env)
		}
	}
	formatted, err := p.formatter.Format(accepted)
	if err != nil {
		return "", err
	}
	return formatted, nil
}

type config struct {
	all       bool
	formatter Formatter
	search    []*regexp.Regexp
}

type Config func(config *config) error

func WithFormat(flag string) Config {
	return func(config *config) error {
		splitted := strings.Split(flag, ",")
		switch splitted[0] {
		case "del":
			if len(splitted) != 2 {
				return fmt.Errorf("formatter \"del\" takes exactly one configuration argument, e.g. del,=. Got: %s", flag)
			}
			config.formatter = &format.DelimiterFormatter{Delimiter: splitted[1]}
		case "json":
			if len(splitted) != 2 {
				return jsonFormatterError(flag)
			}
			switch splitted[1] {
			case "compact":
				config.formatter = &format.JsonFormatter{Pretty: false}
			case "pretty":
				config.formatter = &format.JsonFormatter{Pretty: true}
			default:
				return jsonFormatterError(flag)
			}
		default:
			return fmt.Errorf("unknown formatter: %s", flag)
		}
		return nil
	}
}

func jsonFormatterError(flag string) error {
	return fmt.Errorf("formatter \"json\" takes exactly one configuration argument of compact or pretty, e.g. json,compact. Got: %s", flag)
}

func WithAll() Config {
	return func(config *config) error {
		config.all = true
		return nil
	}
}

func WithSearch(search string) Config {
	return func(config *config) error {
		regex, err := regexp.Compile(search)
		if err != nil {
			return err
		}
		config.search = append(config.search, regex)
		return nil
	}
}

func NewPipeline(configs ...Config) (*Pipeline, error) {
	cfg := config{}
	for _, c := range configs {
		err := c(&cfg)
		if err != nil {
			return nil, err
		}
	}

	r := env.DefaultReader{}

	var filters []filter.Filter
	if !cfg.all {
		filters = append(filters, &filter.NoPrefixFilter{Prefix: "_"})
	}
	for _, s := range cfg.search {
		filters = append(filters, &filter.RegexFilter{Regex: s})
	}
	h := filter.AndFilter{Filters: filters}

	return &Pipeline{&r, &h, cfg.formatter}, nil
}
