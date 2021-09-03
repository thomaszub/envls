package internal

import (
	"os"
	"strings"
)

type EnvReader interface {
	Read() []EnvVar //TODO Directory to read?
}

func NewDefaultReader() EnvReader {
	return &DefaultEnvReader{}
}

type DefaultEnvReader struct {
}

func (r *DefaultEnvReader) Read() []EnvVar {
	envVars := make([]EnvVar, 0)
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
