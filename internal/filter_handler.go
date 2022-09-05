package internal

type Filter interface {
	Accept(env EnvVar) bool
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

func NewFilterHandler(filters []Filter) FilterHandler {
	return FilterHandler{filters: filters}
}
