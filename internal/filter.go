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

type FilterHandler struct {
	filters []Filter
}

func (f *FilterHandler) Accepted(envs []EnvVar) []EnvVar {
	accepted := make([]EnvVar, 0)
	for _, env := range envs {
		if f.Accept(env) {
			accepted = append(accepted, env)
		}
	}
	return accepted
}

func (f *FilterHandler) Accept(env EnvVar) bool {
	for _, filter := range f.filters {
		if !filter.Accept(env) {
			return false
		}
	}
	return true
}

func (f *FilterHandler) AppendFilter(filter Filter) {
	f.filters = append(f.filters, filter)
}

func NewEmptyFilterHandler() FilterHandler {
	return FilterHandler{filters: make([]Filter, 0)}
}
