// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2025 Daniel Sage
// SPDX-FileType: SOURCE

package build

import (
	"runtime/debug"
)

// DetectedVersion indicates the VCS commit info when the project was compiled.
var DetectedVersion string

// GoVersion indicates the version of Go used when the project was compiled.
var GoVersion string

// staticVersion is the project version injected into the binary during the build process.
var staticVersion string

// detectVersion uses the embedded build information to determine the version of the project. This automatically run
// during initialization and should never be required.
func detectVersion() {
	// initialize the versions as "unknown"
	DetectedVersion, GoVersion = "unknown", "unknown"

	// read the build information from the binary
	buildInfo, ok := debug.ReadBuildInfo()
	if !ok {
		// notest
		return
	}

	// store the Go version
	GoVersion = buildInfo.GoVersion

	// check if the static version was set during compilation
	if staticVersion != "" {
		// store the static version as the detected version
		DetectedVersion = staticVersion
		return
	}

	rev := "dev"
	modified := false

	// loop through the build settings to locate the VCS info
	for _, setting := range buildInfo.Settings {
		// notest
		// justification:
		// Go doesn't seem to store the VCS info in build info during testing, making this impossible for us to test as-is

		// check if this is the revision
		if setting.Key == "vcs.revision" {
			// use the first 10 characters of the hash
			rev = setting.Value[:10]
		}

		// check if this is the modified status
		if setting.Key == "vcs.modified" {
			modified = setting.Value == "true"
		}
	}

	// store the detected version
	if modified {
		// notest
		// justification:
		// Go doesn't seem to store the VCS info in build info during testing, making this impossible for us to test as-is

		DetectedVersion = rev + "-dirty"
	} else {
		// notest
		// justification:
		// Go doesn't seem to store the VCS info in build info during testing, making this impossible for us to test as-is

		DetectedVersion = rev
	}
}

func init() {
	// detect the version automatically if required
	detectVersion()
}
