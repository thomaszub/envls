package filter

import (
	"regexp"
	"strings"

	"github.com/thomaszub/envls/internal/env"
)

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
