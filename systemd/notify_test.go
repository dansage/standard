// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2025 Daniel Sage
// SPDX-FileType: SOURCE

//go:build unix

package systemd_test

import (
	"os"
	"strings"
	"testing"
	"time"

	"go.dsage.org/standard/systemd"
)

var (
	tests = []struct {
		Name     string
		Expected string
		Notify   func() error
		Verify   func(string) bool
	}{
		{
			Name:     "Ready",
			Expected: "READY=1",
			Notify:   systemd.NotifyReady,
			Verify:   verifyExact("READY=1"),
		},
		{
			Name:     "Reloading",
			Expected: "RELOADING=1\nMONOTONIC_USEC={time.Now().MicroUnix()}",
			Notify:   systemd.NotifyReloading,
			Verify: func(got string) bool {
				// verify the output starts with the correct prefix
				return strings.HasPrefix(got, "RELOADING=1\nMONOTONIC_USEC=")
			},
		},
		{
			Name:     "Status",
			Expected: "STATUS=The service is running normally",
			Notify: func() error {
				return systemd.NotifyStatus("The service is running normally")
			},
			Verify: verifyExact("STATUS=The service is running normally"),
		},
		{
			Name:     "Stopping",
			Expected: "STOPPING=1",
			Notify:   systemd.NotifyStopping,
			Verify:   verifyExact("STOPPING=1"),
		},
		{
			Name:     "Watchdog",
			Expected: "WATCHDOG=1",
			Notify:   systemd.NotifyWatchdog,
			Verify:   verifyExact("WATCHDOG=1"),
		},
	}
)

// TestNotify runs a variety of predefined tests to ensure each method reports the correct data to systemd.
func TestNotify(t *testing.T) {
	// loop through the predefined tests
	for _, test := range tests {
		// start a subtest to isolate the result from other tests
		t.Run(test.Name, func(t *testing.T) {
			var socket string
			output := make(chan string, 1)

			// start listening for notifications in a goroutine
			go listen(t, &socket, output)

			// delay 100ms to give the listener time to fully start
			time.Sleep(100 * time.Millisecond)

			// set the `NOTIFY_SOCKET` environment variable
			if err := os.Setenv("NOTIFY_SOCKET", socket); err != nil {
				t.Fatalf("os.Setenv failed: %v", err)
			}

			// call the notify method for the test
			if err := test.Notify(); err != nil {
				t.Fatalf("notify callback failed: %v", err)
			}

			// clear the `NOTIFY_SOCKET` environment variable
			if err := os.Unsetenv("NOTIFY_SOCKET"); err != nil {
				t.Fatalf("os.Unsetenv failed: %v", err)
			}

			// verify the output matches the expected value
			got := <-output
			if !test.Verify(got) {
				t.Fatalf("notification data is wrong,expected %q, got %q", test.Expected, got)
			}
		})
	}

	// clear the `NOTIFY_SOCKET` environment variable
	if err := os.Unsetenv("NOTIFY_SOCKET"); err != nil {
		t.Fatalf("os.Unsetenv failed: %v", err)
	}

	// verify no errors are raised when the `NOTIFY_SOCKET` variable is missing entirely
	if err := systemd.Notify("TEST=1"); err != nil {
		t.Fatalf("notify failed with unset variable: %v", err)
	}

	// set the `NOTIFY_SOCKET` environment variable to an empty string
	if err := os.Setenv("NOTIFY_SOCKET", ""); err != nil {
		t.Fatalf("os.Setenv failed: %v", err)
	}

	// verify no errors are raised when the `NOTIFY_SOCKET` variable is set but empty
	if err := systemd.Notify("TEST=1"); err != nil {
		t.Fatalf("notify failed with empty variable: %v", err)
	}

	// clear the `NOTIFY_SOCKET` environment variable
	if err := os.Unsetenv("NOTIFY_SOCKET"); err != nil {
		t.Fatalf("os.Unsetenv failed: %v", err)
	}
}
