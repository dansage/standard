// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2025 Daniel Sage
// SPDX-FileType: SOURCE

//go:build linux

package systemd_test

import (
	"os"
	"testing"
	"time"

	"go.dsage.org/standard/systemd"
)

// TestWatchdog ensures the background watchdog process operates correctly.
func TestWatchdog(t *testing.T) {
	var socket string
	output := make(chan string, 1)

	// clear the `WATCHDOG_USEC` environment variable
	if err := os.Unsetenv("WATCHDOG_USEC"); err != nil {
		t.Fatalf("os.Unsetenv failed: %v", err)
	}

	// verify the watchdog exits if the interval is not configured
	if err := systemd.Watchdog(func() bool {
		return true
	}); err != nil {
		t.Fatalf("Watchdog failed with no interval: %v", err)
	}

	// set the `WATCHDOG_USEC` environment variable to reflect 2 seconds
	if err := os.Setenv("WATCHDOG_USEC", "2000000"); err != nil {
		t.Fatalf("os.Setenv failed: %v", err)
	}

	// start listening for notifications in a goroutine
	go listen(t, &socket, output)

	// delay 100ms to give the listener time to fully start
	time.Sleep(100 * time.Millisecond)

	// set the `NOTIFY_SOCKET` environment variable
	if err := os.Setenv("NOTIFY_SOCKET", socket); err != nil {
		t.Fatalf("os.Setenv failed: %v", err)
	}

	// start the watchdog, indicating the service is healthy
	var watchdogErr error
	healthy := true
	go func() {
		watchdogErr = systemd.Watchdog(func() bool {
			return healthy
		})
	}()

	// delay 1 second to give the watchdog time to report in
	time.Sleep(time.Second)

	// verify the watchdog reported in
	if <-output != "WATCHDOG=1" {
		t.Fatalf("watchdog did not notify on time")
	}

	// start listening for notifications in a goroutine again (the previous listener closed)
	go listen(t, &socket, output)

	// delay 100ms to give the listener time to fully start
	time.Sleep(100 * time.Millisecond)

	// set the `NOTIFY_SOCKET` environment variable
	if err := os.Setenv("NOTIFY_SOCKET", socket); err != nil {
		t.Fatalf("os.Setenv failed: %v", err)
	}

	// indicate the service is no longer healthy
	healthy = false

	// delay 1 second to give the watchdog time to report in
	time.Sleep(time.Second)

	// verify the watchdog stopped with an error
	if watchdogErr == nil {
		t.Fatalf("watchdog did not exit for unhealthy service")
	}

	// clear the `WATCHDOG_USEC` environment variable
	if err := os.Unsetenv("WATCHDOG_USEC"); err != nil {
		t.Fatalf("os.Unsetenv failed: %v", err)
	}
}

// TestWatchdogInterval ensures the watchdog interval is correctly determined based on the configured value.
func TestWatchdogInterval(t *testing.T) {
	// clear the `WATCHDOG_USEC` environment variable
	if err := os.Unsetenv("WATCHDOG_USEC"); err != nil {
		t.Fatalf("os.Unsetenv failed: %v", err)
	}

	// verify the watchdog interval is discovered as `-1` without error
	interval, err := systemd.WatchdogInterval()
	if err != nil {
		t.Fatalf("WatchdogInterval failed: %v", err)
	}
	if interval != -1 {
		t.Fatalf("WatchdogInterval returned %d, expected -1", interval)
	}

	// set the `WATCHDOG_USEC` environment variable to reflect 30 seconds
	if err := os.Setenv("WATCHDOG_USEC", "30000000"); err != nil {
		t.Fatalf("os.Setenv failed: %v", err)
	}

	// verify the watchdog interval is discovered as 30 seconds without error
	interval, err = systemd.WatchdogInterval()
	if err != nil {
		t.Fatalf("WatchdogInterval failed: %v", err)
	}
	if interval != time.Duration(30)*time.Second {
		t.Fatalf("WatchdogInterval returned %d, expected 30s", interval)
	}

	// set the `WATCHDOG_USEC` environment variable to a non-numeric value
	if err := os.Setenv("WATCHDOG_USEC", "invalid"); err != nil {
		t.Fatalf("os.Setenv failed: %v", err)
	}

	// verify the watchdog interval is not discovered and an error is returned
	interval, err = systemd.WatchdogInterval()
	if err == nil {
		t.Fatalf("WatchdogInterval succeeded with value %q: %v", "invalid", interval)
	}

	// clear the `WATCHDOG_USEC` environment variable
	if err := os.Unsetenv("WATCHDOG_USEC"); err != nil {
		t.Fatalf("os.Unsetenv failed: %v", err)
	}
}
