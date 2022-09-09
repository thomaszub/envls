package filter

import (
	"regexp"
	"strings"

	"github.com/thomaszub/envls/internal/env"
)

type Filter interface {
	Accept(env env.Var) bool
}

type NoPrefixFilter struct {
	Prefix string
}

func (f *NoPrefixFilter) Accept(env env.Var) bool {
	return !strings.HasPrefix(env.Name, f.Prefix)
}

type RegexFilter struct {
	Regex *regexp.Regexp
}

func (f *RegexFilter) Accept(env env.Var) bool {
	r := f.Regex
	return r.MatchString(env.Name) || r.MatchString(env.Value)
}

type AndFilter struct {
	Filters []Filter
}

func (f *AndFilter) Accept(env env.Var) bool {
	for _, filter := range f.Filters {
		if !filter.Accept(env) {
			return false
		}
	}
	return true
}
