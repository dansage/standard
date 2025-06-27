// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2025 Daniel Sage
// SPDX-FileType: SOURCE

//go:build unix

package systemd

import (
	"fmt"
	"net"
	"os"
)

// Notify reports the specified state to systemd using the configured notify socket.
func Notify(state string) error {
	// pull the configured socket address from the environment
	ns, ok := os.LookupEnv("NOTIFY_SOCKET")
	if !ok || ns == "" {
		return nil
	}

	// open the configured socket
	conn, err := net.Dial("unixgram", ns)
	if err != nil {
		return fmt.Errorf("failed to open notify socket: %w", err)
	}
	defer conn.Close()

	// write the service state to the socket
	if _, err := conn.Write([]byte(state)); err != nil {
		return fmt.Errorf("failed to write state to notify socket: %w", err)
	}
	return nil
}

// Notifyf reports the specified state (after formatting) to systemd using the configured notify socket.
func Notifyf(format string, args ...any) error {
	return Notify(fmt.Sprintf(format, args...))
}
