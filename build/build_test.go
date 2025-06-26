// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2025 Daniel Sage
// SPDX-FileType: SOURCE

package build

import (
	"runtime"
	"testing"
)

// TestGoVersionIsDetected ensures the Go runtime version is automatically detected and is correct
func TestGoVersionIsDetected(t *testing.T) {
	// get the actual Go version from the runtime
	ver := runtime.Version()

	// verify the detected Go version matches the runtime version
	if GoVersion != ver {
		t.Fatalf("the detected Go version does not match: expected %q got %q", ver, GoVersion)
	}
}

// TestStaticVersionIsDetected ensures the static version overrides the detection logic, if present
func TestStaticVersionIsDetected(t *testing.T) {
	// save the static version to restore it later
	static := staticVersion

	// set the static version to a marker value
	staticVersion = "test-marker"

	// re-detect the version
	detectVersion()

	// verify the detected version matches the marker value
	if DetectedVersion != staticVersion {
		t.Fatalf("the detected version does not match: expected %q got %q", staticVersion, DetectedVersion)
	}

	// restore the original static version value
	staticVersion = static

	// re-detect the version (to restore pre-test state)
	detectVersion()
}
