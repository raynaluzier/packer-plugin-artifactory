// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

// Reads the relative .env file with required environment variables
// Ensures the file doesn't include invalid characters
// For each listed variable, sets the key/value pair as an environment variable so it can be pulled
// into the project for use

package main

import (
	"fmt"
	"os"
	"strings"
)

func loadDotEnv() error {
	data, err := os.ReadFile(".env")
	if err != nil {
		return err
	}
	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "#") {
			continue
		}
		kvp := strings.SplitN(line, "=", 2)
		if len(kvp) != 2 || kvp[0] == "" || strings.Contains(kvp[0], " ") || strings.Contains(kvp[0], "'") || 
	strings.Contains(kvp[0], "\"") {
			fmt.Printf("Warning: invalid line in .env: %s\n", line)
			continue
		}
		key := kvp[0]
		value := kvp[1]
		if value[0] == '"' {
			if value[len(value)-1] != '"' {
				fmt.Printf("Warning: invalid line in .env: %s\n", line)
				continue
			}
			value = value[1 : len(value)-1]
		} else if value[0] == '\'' {
			if value[len(value)-1] != '\'' {
				fmt.Printf("Warning: invalid line in .env: %s\n", line)
				continue
			}
			value = value[: len(value)-1]
		}
		os.Setenv(key, value)
	}
	return nil
}

func init() {
	err := loadDotEnv()
	if err != nil {
		fmt.Println(err)
	}
}
