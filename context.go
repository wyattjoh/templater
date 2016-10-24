package main

import (
	"encoding/json"
	"os"
	"strings"
)

func loadContext(filename string) (*Context, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var ctx Context
	if err := json.NewDecoder(f).Decode(&ctx); err != nil {
		return nil, err
	}

	return &ctx, nil
}

// Code here is sourced from the dockerize repository: https://github.com/jwilder/dockerize

// Context contains the application context.
type Context map[string]interface{}

// Env loads the OS environment variable.
func (c *Context) Env() map[string]string {
	env := make(map[string]string)
	for _, i := range os.Environ() {
		sep := strings.Index(i, "=")
		env[i[0:sep]] = i[sep+1:]
	}
	return env
}
