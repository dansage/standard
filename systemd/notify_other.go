// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2025 Daniel Sage
// SPDX-FileType: SOURCE

//go:build !unix

package systemd

// Notify reports the specified state to systemd using the configured notify socket.
func Notify(state string) error {
	return nil
}

// Notifyf reports the specified state (after formatting) to systemd using the configured notify socket.
func Notifyf(format string, args ...any) error {
	return nil
}
