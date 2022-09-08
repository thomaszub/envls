package internal

type Filter interface {
	Accept(env EnvVar) bool
}

type FilterHandler struct {
	Filters []Filter
}

func (f *FilterHandler) Accepted(envs []EnvVar) []EnvVar {
	var accepted []EnvVar
	for _, env := range envs {
		if f.accept(env) {
			accepted = append(accepted, env)
		}
	}
	return accepted
}

func (f *FilterHandler) accept(env EnvVar) bool {
	for _, filter := range f.Filters {
		if !filter.Accept(env) {
			return false
		}
	}
	return true
}
