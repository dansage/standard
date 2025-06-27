// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2025 Daniel Sage
// SPDX-FileType: SOURCE

//go:build unix

package systemd

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

// Watchdog is a background method that continuously polls the service to determine if it reports it is still healthy
// and sends a notification to systemd halfway through the configured interval period. If the service reports it is
// unhealthy, an error is returned. If the watchdog is not configured for the service, the method returns without an
// error. This method blocks the current thread and should be run in a goroutine for best effect.
func Watchdog(healthy func() bool) error {
	// pull the configured watchdog interval
	interval, err := WatchdogInterval()
	if err != nil {
		return err
	}

	// verify the interval is greater than zero
	if interval <= 0 {
		return nil
	}

	for {
		// wait for half of the interval (recommended by `sd_watchdog_enabled(3)`)
		time.Sleep(interval / 2)

		// verify the service is still healthy
		if !healthy() {
			return fmt.Errorf("the service reported it is unhealthy")
		}

		// notify systemd the service is still alive and running
		if err = NotifyWatchdog(); err != nil {
			return fmt.Errorf("failed to send watchdog notification: %w", err)
		}
	}
}

// WatchdogInterval pulls the configured watchdog interval from the environment. The interval specifies how often the
// service should notify systemd it is still alive and running. If the returned interval is `-1`, the watchdog is not
// enabled for this service.
func WatchdogInterval() (time.Duration, error) {
	// pull the configured watchdog interval (in microseconds) from the environment
	usec, ok := os.LookupEnv("WATCHDOG_USEC")
	if !ok {
		return -1, nil
	}

	// convert the value into a usable uint64
	u, err := strconv.ParseUint(usec, 10, 64)
	if err != nil {
		return -1, fmt.Errorf("failed to parse WATCHDOG_USEC value %q: %w", usec, err)
	}
	return time.Duration(u) * time.Microsecond, nil
}
