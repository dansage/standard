// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2025 Daniel Sage
// SPDX-FileType: SOURCE

//go:build !unix

package systemd

import "time"

// Watchdog is a background method that continuously polls the service to determine if it reports it is still healthy
// and sends a notification to systemd halfway through the configured interval period. If the service reports it is
// unhealthy, an error is returned. If the watchdog is not configured for the service, the method returns without an
// error. This method blocks the current thread and should be run in a goroutine for best effect.
func Watchdog(healthy func() bool) error {
	return nil
}

// WatchdogInterval pulls the configured watchdog interval from the environment. The interval specifies how often the
// service should notify systemd it is still alive and running. If the returned interval is `-1`, the watchdog is not
// enabled for this service.
func WatchdogInterval() (time.Duration, error) {
	return -1, nil
}
