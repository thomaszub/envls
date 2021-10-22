package internal

import "strings"

type Filter interface {
	Accept(env EnvVar) bool
}

type NoPrefixFilter struct {
	prefix string
}

func (f *NoPrefixFilter) Accept(env EnvVar) bool {
	return !strings.HasPrefix(env.Name, f.prefix)
}

func NewNoPrefixFilter(prefix string) Filter {
	return &NoPrefixFilter{prefix: prefix}
}
