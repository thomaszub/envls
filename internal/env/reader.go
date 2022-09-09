package env

import (
	"os"
	"strings"
)

type DefaultReader struct {
}

func (r *DefaultReader) Read() []Var {
	var envVars []Var
	for _, s := range os.Environ() {
		envVar := createEnvVar(s)
		envVars = append(envVars, envVar)
	}
	return envVars
}

func createEnvVar(s string) Var {
	keyValuePair := strings.Split(s, "=")
	return Var{
		Name:  keyValuePair[0],
		Value: keyValuePair[1],
	}
}
