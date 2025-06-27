// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2025 Daniel Sage
// SPDX-FileType: SOURCE

//go:build unix

package systemd_test

import (
	"net"
	"path/filepath"
	"testing"
)

// listen for notifications on a temporary socket.
func listen(t *testing.T, socket *string, output chan string) {
	// create a notification socket in a new temporary directory
	*socket = filepath.Join(t.TempDir(), "notify.socket")

	// listen for notifications from the calling test
	conn, err := net.ListenPacket("unixgram", *socket)
	if err != nil {
		t.Fatalf("net.ListenPacket failed: %v", err)
	}
	defer conn.Close()

	// listen for a single packet of data
	buf := make([]byte, 1024)
	n, _, err := conn.ReadFrom(buf)
	if err != nil {
		t.Fatalf("conn.ReadFrom failed: %v", err)
	}

	// write the packet to the output string
	output <- string(buf[:n])
}

// verifyExact the received value matches the expected value exactly.
func verifyExact(expected string) func(string) bool {
	return func(got string) bool {
		// verify the value matches
		return expected == got
	}
}
