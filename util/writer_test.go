// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2025 Daniel Sage
// SPDX-FileType: SOURCE

package util

import (
	"bytes"
	"os"
	"testing"
)

func TestMultiWriter_Register(t *testing.T) {
	wr := NewMultiWriter()

	// verify no writers are registered already
	if len(wr.writers) > 0 {
		t.Fatalf("multi writer has writers registered")
	}

	// register os.Stdout twice
	wr.Register(os.Stdout)
	wr.Register(os.Stdout)

	// verify exactly one writer is registered and it is os.Stdout
	if len(wr.writers) != 1 {
		t.Fatalf("multi writer should have 1 writer, has: %d", len(wr.writers))
	}
	if wr.writers[0] != os.Stdout {
		t.Fatalf("multi writer should have registered os.Stdout, has: %s", wr.writers[0])
	}
}

func TestMultiWriter_Unregister(t *testing.T) {
	wr := NewMultiWriter(os.Stderr, os.Stdout)

	// verify exactly two writers are registered
	if len(wr.writers) != 2 {
		t.Fatalf("multi writer should have 2 writers, has: %d", len(wr.writers))
	}

	// unregister os.Stderr
	wr.Unregister(os.Stderr)

	// verify exactly one writer is registered and it is os.Stdout
	if len(wr.writers) != 1 {
		t.Fatalf("multi writer should have 1 writer, has: %d", len(wr.writers))
	}
	if wr.writers[0] != os.Stdout {
		t.Fatalf("multi writer should have registered os.Stdout, has: %s", wr.writers[0])
	}
}

func TestMultiWriter_Write(t *testing.T) {
	var buf1, buf2 bytes.Buffer
	wr := NewMultiWriter()

	// register both buffers as writers
	wr.Register(&buf1)
	wr.Register(&buf2)

	// write a sample buffer to the writer
	expected := []byte("hello world")
	if _, err := wr.Write(expected); err != nil {
		t.Fatalf("failed to write to both buffers: %v", err)
	}

	// verify the first buffer contains the correct data
	if !bytes.Equal(buf1.Bytes(), expected) {
		t.Fatalf("first buffer should have been %v, has: %v", expected, buf1.Bytes())
	}

	// verify the second buffer contains the correct data
	if !bytes.Equal(buf2.Bytes(), expected) {
		t.Fatalf("second buffer should have been %v, has: %v", expected, buf2.Bytes())
	}

	// unregister the first buffer
	wr.Unregister(&buf1)

	// clear both buffers
	buf1.Reset()
	buf2.Reset()

	// write a sample buffer to the writer
	expected = []byte("goodbye world")
	if _, err := wr.Write(expected); err != nil {
		t.Fatalf("failed to write to both buffers: %v", err)
	}

	// verify the first buffer is empty
	if buf1.Len() != 0 {
		t.Fatalf("first buffer should have been empty, has: %v", buf1.Bytes())
	}

	// verify the second buffer contains the correct data
	if !bytes.Equal(buf2.Bytes(), expected) {
		t.Fatalf("second buffer should have been %v, has: %v", expected, buf2.Bytes())
	}
}
