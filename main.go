package main

import "fmt"

func main() {
	envReader := NewDefaultReader()
	filterChain := NewEmptyFilterChain()
	filterChain.AppendFilter(&PrefixFilter{prefix: "_"})
	formatter := DelimiterFormatter{delimiter: " -> "}

	envs := envReader.Read()
	filteredEnvs := filterChain.Filter(envs)
	formattedEnvs := formatter.Format(filteredEnvs)
	for _, formattedEnv := range formattedEnvs {
		fmt.Println(formattedEnv)
	}
}
