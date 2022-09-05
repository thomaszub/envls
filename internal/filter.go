package internal

import (
	"regexp"
	"strings"
)

type NoPrefixFilter struct {
	prefix string
}

type RegexFilter struct {
	regex *regexp.Regexp
}

func (f *NoPrefixFilter) Accept(env EnvVar) bool {
	return !strings.HasPrefix(env.Name, f.prefix)
}

func NewNoPrefixFilter(prefix string) NoPrefixFilter {
	return NoPrefixFilter{prefix: prefix}
}

func (f *RegexFilter) Accept(env EnvVar) bool {
	r := f.regex
	return r.MatchString(env.Name) || r.MatchString(env.Value)
}

func NewRegexFilter(regex *regexp.Regexp) RegexFilter {
	return RegexFilter{regex: regex}
}
