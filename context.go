package main

import (
	"encoding/json"
	"os"
)

func loadContext(filename string) (map[string]interface{}, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var ctx map[string]interface{}
	if err := json.NewDecoder(f).Decode(&ctx); err != nil {
		return nil, err
	}

	return ctx, nil
}
