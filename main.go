package main

import "fmt"

func main() {
	envReader := NewDefaultReader()
	envVars := envReader.Read()
	for _, envVar := range envVars {
		fmt.Println(envVar)
	}
}
