package internal

import (
	"regexp"
	"strings"
)

type Filter interface {
	Accept(env EnvVar) bool
}

type NoPrefixFilter struct {
	prefix string
}

type RegexFilter struct {
	regex *regexp.Regexp
}

func (f *NoPrefixFilter) Accept(env EnvVar) bool {
	return !strings.HasPrefix(env.Name, f.prefix)
}

func NewNoPrefixFilter(prefix string) Filter {
	return &NoPrefixFilter{prefix: prefix}
}

func (f *RegexFilter) Accept(env EnvVar) bool {
	accept := f.regex.MatchString(env.Name)
	if accept {
		return true
	}
	accept = f.regex.MatchString(env.Value)
	return accept
}

func NewRegexFilter(regex *regexp.Regexp) Filter {
	return &RegexFilter{regex: regex}
}
