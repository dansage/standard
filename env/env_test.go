// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2025 Daniel Sage
// SPDX-FileType: SOURCE

package env_test

import (
	"os"
	"strings"
	"testing"

	"go.dsage.org/standard/env"
)

// TestDefaultEnvironmentIsProduction ensures the environment type is detected as production when either no value is set
// or the ENV variable is set to an unknown value
func TestDefaultEnvironmentIsProduction(t *testing.T) {
	// get the initial value of the ENV variable
	initial, initialOk := os.LookupEnv("ENV")

	// pull all environment variables from the environment
	all := os.Environ()

	// clear all environment variables
	os.Clearenv()

	// loop through the variables and restore all values except the environment
	for _, e := range all {
		// split the value by the equal sign
		split := strings.Split(e, "=")

		// split the key and value for easier access
		key := split[0]
		value := split[1]

		// ensure the value is restored correctly if there are more items in the slice
		if len(split) > 2 {
			// join all but the first element with equal signs
			value = strings.Join(split[1:], "=")
		}

		// if the key is not ENV, restore the value
		if key != "ENV" {
			if err := os.Setenv(key, value); err != nil {
				t.Fatalf("failed to restore environment variable %q: %v", key, err)
			}
		}
	}

	// detect the current environment type
	current := env.DetectEnvironment()

	// verify the environment is detected as Production
	if current != env.Production {
		t.Fatalf("environment type was detected as %q instead of %q", current, env.Production)
	}

	// intentionally set the ENV variable to an unknown value
	if err := os.Setenv("ENV", "unknown"); err != nil {
		t.Fatalf("failed to set environment variable \"ENV\": %v", err)
	}

	// re-detect the current environment type
	current = env.DetectEnvironment()

	// verify the environment is detected as Production
	if current != env.Production {
		t.Fatalf("environment type was detected as %q instead of %q", current, env.Production)
	}

	// if set, restore the initial value of the ENV variable
	if initialOk {
		if err := os.Setenv("ENV", initial); err != nil {
			t.Fatalf("failed to restore environment variable \"ENV\": %v", err)
		}
	}
}

// TestEnvironmentIsDevelopment ensures the environment type is detected as development for all valid values
func TestEnvironmentIsDevelopment(t *testing.T) {
	// get the initial value of the ENV variable
	initial, initialOk := os.LookupEnv("ENV")

	// set the ENV variable to "dev"
	if err := os.Setenv("ENV", "dev"); err != nil {
		t.Fatalf("failed to set environment variable \"ENV\": %v", err)
	}

	// detect the current environment type
	current := env.DetectEnvironment()

	// verify the environment is detected as Testing
	if current != env.Development {
		t.Fatalf("environment type was detected as %q instead of %q", current, env.Development)
	}

	// set the ENV variable to "development"
	if err := os.Setenv("ENV", "development"); err != nil {
		t.Fatalf("failed to set environment variable \"ENV\": %v", err)
	}

	// re-detect the current environment type
	current = env.DetectEnvironment()

	// verify the environment is detected as Development
	if current != env.Development {
		t.Fatalf("environment type was detected as %q instead of %q", current, env.Development)
	}

	// set the ENV variable to "local"
	if err := os.Setenv("ENV", "local"); err != nil {
		t.Fatalf("failed to set environment variable \"ENV\": %v", err)
	}

	// re-detect the current environment type
	current = env.DetectEnvironment()

	// verify the environment is detected as Development
	if current != env.Development {
		t.Fatalf("environment type was detected as %q instead of %q", current, env.Development)
	}

	// if set, restore the initial value of the ENV variable
	if initialOk {
		if err := os.Setenv("ENV", initial); err != nil {
			t.Fatalf("failed to restore environment variable \"ENV\": %v", err)
		}
	}
}

// TestEnvironmentIsSet ensures the current environment type is automatically set during initialization
func TestEnvironmentIsSet(t *testing.T) {
	// verify the current environment is automatically set
	if env.Current != env.DetectEnvironment() {
		t.Fatalf("current environment was not automatically set")
	}
}

// TestEnvironmentIsTest ensures the environment type is detected as testing initially and for all valid values
func TestEnvironmentIsTest(t *testing.T) {
	// get the initial value of the ENV variable
	initial, initialOk := os.LookupEnv("ENV")

	// detect the current environment type
	current := env.DetectEnvironment()

	// verify the environment is detected as Testing
	if current != env.Testing {
		t.Fatalf("environment type was detected as %q instead of %q", current, env.Testing)
	}

	// set the ENV variable to "test"
	if err := os.Setenv("ENV", "test"); err != nil {
		t.Fatalf("failed to set environment variable \"ENV\": %v", err)
	}

	// re-detect the current environment type
	current = env.DetectEnvironment()

	// verify the environment is detected as Testing
	if current != env.Testing {
		t.Fatalf("environment type was detected as %q instead of %q", current, env.Testing)
	}

	// set the ENV variable to "testing"
	if err := os.Setenv("ENV", "testing"); err != nil {
		t.Fatalf("failed to set environment variable \"ENV\": %v", err)
	}

	// re-detect the current environment type
	current = env.DetectEnvironment()

	// verify the environment is detected as Testing
	if current != env.Testing {
		t.Fatalf("environment type was detected as %q instead of %q", current, env.Testing)
	}

	// if set, restore the initial value of the ENV variable
	if initialOk {
		if err := os.Setenv("ENV", initial); err != nil {
			t.Fatalf("failed to restore environment variable \"ENV\": %v", err)
		}
	}
}
