package internal

import (
	"os"
	"strings"
)

type DefaultEnvReader struct {
}

func (r *DefaultEnvReader) Read() []EnvVar {
	var envVars []EnvVar
	for _, s := range os.Environ() {
		envVar := createEnvVar(s)
		envVars = append(envVars, envVar)
	}
	return envVars
}

func createEnvVar(s string) EnvVar {
	keyValuePair := strings.Split(s, "=")
	return EnvVar{
		Name:  keyValuePair[0],
		Value: keyValuePair[1],
	}
}
