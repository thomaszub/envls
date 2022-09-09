package filter

import "github.com/thomaszub/envls/internal/env"

type Filter interface {
	Accept(env env.Var) bool
}

type FilterHandler struct {
	Filters []Filter
}

func (f *FilterHandler) Accepted(envs []env.Var) []env.Var {
	var accepted []env.Var
	for _, env := range envs {
		if f.accept(env) {
			accepted = append(accepted, env)
		}
	}
	return accepted
}

func (f *FilterHandler) accept(env env.Var) bool {
	for _, filter := range f.Filters {
		if !filter.Accept(env) {
			return false
		}
	}
	return true
}
