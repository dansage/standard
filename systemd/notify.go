// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2025 Daniel Sage
// SPDX-FileType: SOURCE

package systemd

import "time"

// NotifyReady reports that the service is ready to systemd using the configured notify socket.
func NotifyReady() error {
	return Notify("READY=1")
}

// NotifyReloading reports that the service is reloading to systemd using the configured notify socket. Once the reload
// is complete, the service should report it is ready using `NotifyReady` to avoid the reload job timing out.
func NotifyReloading() error {
	return Notifyf("RELOADING=1\nMONOTONIC_USEC=%d", time.Now().UnixMicro())
}

// NotifyStatus reports the specified user-readable status to system using the configured notify socket.
func NotifyStatus(status string) error {
	return Notifyf("STATUS=%s", status)
}

// NotifyStopping reports that the service is stopping to systemd using the configured notify socket.
func NotifyStopping() error {
	return Notify("STOPPING=1")
}

// NotifyWatchdog reports that the service is alive to systemd using the configured notify socket.
func NotifyWatchdog() error {
	return Notify("WATCHDOG=1")
}
