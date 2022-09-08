package internal

import (
	"regexp"
	"strings"
)

type NoPrefixFilter struct {
	Prefix string
}

type RegexFilter struct {
	Regex *regexp.Regexp
}

func (f *NoPrefixFilter) Accept(env EnvVar) bool {
	return !strings.HasPrefix(env.Name, f.Prefix)
}

func (f *RegexFilter) Accept(env EnvVar) bool {
	r := f.Regex
	return r.MatchString(env.Name) || r.MatchString(env.Value)
}
