// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2025 Daniel Sage
// SPDX-FileType: SOURCE

package env

import (
	"os"
	"strings"
)

// Current indicates the environment the application is running in.
var Current Type

// Type represents the different type of application environments.
type Type string

const (
	// Production indicates the application is running in a production environment.
	Production Type = "production"

	// Testing indicates the application is running in a CI testing environment.
	Testing Type = "testing"

	// Development indicates the application is running in a local development environment.
	Development Type = "development"
)

// DetectEnvironment uses the environment variable ENV to determine what type of environment this is.
func DetectEnvironment() Type {
	// pull the environment string from the environment (how appropriate!)
	env, _ := os.LookupEnv("ENV")

	// verify the environment string is valid
	switch strings.ToLower(env) {

	// look for a local development environment
	case "dev", "development", "local":
		return Development

	// look for a testing environment
	case "test", "testing":
		return Testing
	}

	// fail closed by assuming production in every other circumstance
	return Production
}

func init() {
	// detect the current environment immediately
	Current = DetectEnvironment()
}
