package main

import "strings"

type Filter interface {
	Accept(env EnvVar) bool
}

type PrefixFilter struct {
	prefix string
}

func (f *PrefixFilter) Accept(env EnvVar) bool {
	return !strings.HasPrefix(env.name, f.prefix)
}

type FilterHandler struct {
	filters []Filter
}

func (f *FilterHandler) Filter(envs []EnvVar) []EnvVar {
	filtered := make([]EnvVar, 0)
	for _, env := range envs {
		if f.Accept(env) {
			filtered = append(filtered, env)
		}
	}
	return filtered
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

func NewEmptyFilterChain() FilterHandler {
	return FilterHandler{filters: make([]Filter, 0)}
}
