package main

import (
	"os"
	"strings"
)

type action struct {
	inputs map[string]string
}

func newAction() *action {
	// New Action Object
	a := &action{}
	// Set Input Arguments map
	a.inputs = make(map[string]string)
	// Get input arguments
	for _, y := range os.Environ() {
		s := strings.SplitN(y, `=`, 2)
		if strings.HasPrefix(s[0], "INPUT_") {
			a.inputs[strings.TrimPrefix(s[0], "INPUT_")] = s[1]
		}
	}

	//
	return a
}

func (a *action) GetInput(key string) string {
	return a.inputs[strings.ToUpper(key)]
}
