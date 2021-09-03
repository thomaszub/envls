package main

import "fmt"

func main() {
	envReader := NewDefaultReader()
	envs := envReader.Read()
	filterChain := NewEmptyFilterChain()
	filterChain.AppendFilter(&PrefixFilter{prefix: "_"})
	filteredEnvs := filterChain.Filter(envs)
	for _, env := range filteredEnvs {
		fmt.Println(env)
	}
}
