// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2025 Daniel Sage
// SPDX-FileType: SOURCE

package util

import (
	"bytes"
	"errors"
	"io"
	"os"
	"testing"
)

func TestMultiWriter_Register(t *testing.T) {
	wr := NewMultiWriter()

	// verify no writers are registered already
	if len(wr.writers) > 0 {
		t.Fatalf("multi writer has writers registered")
	}

	// register os.Stdout twice, and nil once
	wr.Register(os.Stdout)
	wr.Register(os.Stdout)
	wr.Register(nil)

	// verify exactly one writer is registered and it is os.Stdout
	if len(wr.writers) != 1 {
		t.Fatalf("multi writer should have 1 writer, has: %d", len(wr.writers))
	}
	if wr.writers[0] != os.Stdout {
		t.Fatalf("multi writer should have registered os.Stdout, has: %s", wr.writers[0])
	}

	// register os.Stdout twice with NewMultiWriter
	wr = NewMultiWriter(os.Stdout, os.Stdout)

	// verify exactly one writer is registered and it is os.Stdout
	if len(wr.writers) != 1 {
		t.Fatalf("multi writer should have 1 writer, has: %d", len(wr.writers))
	}
	if wr.writers[0] != os.Stdout {
		t.Fatalf("multi writer should have registered os.Stdout, has: %s", wr.writers[0])
	}
}

// notComparable is a test writer that cannot be compared using only standard operators.
type notComparable struct {
	b []byte
}

// Write discards the provided data and always reports the entire buffer was written without any errors.
func (n notComparable) Write(p []byte) (int, error) {
	return len(p), nil
}

func TestMultiWriter_Register_NotComparable(t *testing.T) {
	wr := NewMultiWriter()

	// verify no writers are registered already
	if len(wr.writers) > 0 {
		t.Fatalf("multi writer has writers registered")
	}

	// register os.Stdout and a noncomparable writer twice
	wr.Register(os.Stdout)
	wr.Register(notComparable{})
	wr.Register(notComparable{})

	// verify exactly one writer is registered and it is os.Stdout
	if len(wr.writers) != 1 {
		t.Fatalf("multi writer should have 1 writer, has: %d", len(wr.writers))
	}
	if wr.writers[0] != os.Stdout {
		t.Fatalf("multi writer should have registered os.Stdout, has: %s", wr.writers[0])
	}

	// register a comparable writer of a different type
	wr.Register(bytes.NewBuffer(nil))

	// verify exactly two writers are registered
	if len(wr.writers) != 2 {
		t.Fatalf("multi writer should have 2 writers, has: %d", len(wr.writers))
	}
}

func TestMultiWriter_Register_NilPointer(t *testing.T) {
	var nilBuf *bytes.Buffer
	var buf bytes.Buffer
	wr := NewMultiWriter()

	// verify no writers are registered already
	if len(wr.writers) > 0 {
		t.Fatalf("multi writer has writers registered")
	}

	// register a valid buffer and a nil pointer writer
	wr.Register(&buf)
	wr.Register(nilBuf)

	// verify exactly one writer is registered and it is the valid buffer
	if len(wr.writers) != 1 {
		t.Fatalf("multi writer should have 1 writer, has: %d", len(wr.writers))
	}
	if wr.writers[0] != &buf {
		t.Fatalf("multi writer should have registered the valid buffer, has: %v", wr.writers[0])
	}

	// verify writing does not panic on the nil pointer writer
	expected := []byte("hello world")
	if _, err := wr.Write(expected); err != nil {
		t.Fatalf("failed to write: %v", err)
	}
	if !bytes.Equal(buf.Bytes(), expected) {
		t.Fatalf("buffer should have been %v, has: %v", expected, buf.Bytes())
	}

	// verify NewMultiWriter also skips the nil pointer writer
	wr = NewMultiWriter(nilBuf)
	if len(wr.writers) != 0 {
		t.Fatalf("multi writer should have 0 writers, has: %d", len(wr.writers))
	}
}

func TestMultiWriter_Unregister(t *testing.T) {
	var buf bytes.Buffer
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

	// unregister the unused buffer
	wr.Unregister(&buf)

	// verify exactly one writer is registered and it is os.Stdout
	if len(wr.writers) != 1 {
		t.Fatalf("multi writer should have 1 writer, has: %d", len(wr.writers))
	}
	if wr.writers[0] != os.Stdout {
		t.Fatalf("multi writer should have registered os.Stdout, has: %s", wr.writers[0])
	}

	// unregister a noncomparable writer
	wr.Unregister(notComparable{})

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
	n, err := wr.Write(expected)
	if err != nil {
		t.Fatalf("failed to write to both buffers: %v", err)
	}
	if n != len(expected) {
		t.Fatalf("failed to write to both buffers: expected %d bytes, got %d", len(expected), n)
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
		t.Fatalf("failed to write to buffer: %v", err)
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

// errorWriter is a test writer that always reports errors.
type errorWriter struct{}

// Write discards the provided data and always reports the entire buffer was written with an io.ErrUnexpectedEOF error.
func (e errorWriter) Write(p []byte) (int, error) {
	return len(p), io.ErrUnexpectedEOF
}

func TestMultiWriter_Write_Error(t *testing.T) {
	var buf bytes.Buffer
	wr := NewMultiWriter()

	// register a special writer that always reports a failure, and a normal buffer
	wr.Register(errorWriter{})
	wr.Register(&buf)

	// write a sample buffer to the writer
	expected := []byte("hello world")
	n, err := wr.Write(expected)
	if !errors.Is(err, io.ErrUnexpectedEOF) {
		t.Fatalf("expected unexpected EOF error, got: %v", err)
	}
	if n != len(expected) {
		t.Fatalf("expected %d bytes, got %d", len(expected), n)
	}

	// verify the normal buffer remains empty
	if buf.Len() != 0 {
		t.Fatalf("second buffer should have been empty, has %d bytes", buf.Len())
	}
}

func TestMultiWriter_Write_NoWriters(t *testing.T) {
	wr := NewMultiWriter()

	// write a sample buffer to the writer
	expected := []byte("hello world")
	n, err := wr.Write(expected)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if n != len(expected) {
		t.Fatalf("expected %d bytes, got %d", len(expected), n)
	}
}

// shortWriter is a test writer that reports short writes without errors.
type shortWriter struct{}

// Write discards the provided data and always reports it wrote 1 byte less than the provided buffer, or zero, without
// any errors.
func (shortWriter) Write(p []byte) (int, error) {
	return max(0, len(p)-1), nil
}

func TestMultiWriter_Write_ShortWrite(t *testing.T) {
	var buf bytes.Buffer
	wr := NewMultiWriter()

	// register a special writer that should always fail with non-empty buffers, and a normal buffer
	wr.Register(shortWriter{})
	wr.Register(&buf)

	// write a sample buffer to the writer
	expected := []byte("hello world")
	n, err := wr.Write(expected)
	if !errors.Is(err, io.ErrShortWrite) {
		t.Fatalf("expected short write error, got: %v", err)
	}
	if n != len(expected)-1 {
		t.Fatalf("expected %d bytes, got %d", len(expected)-1, n)
	}

	// verify the normal buffer remains empty
	if buf.Len() != 0 {
		t.Fatalf("second buffer should have been empty, has %d bytes", buf.Len())
	}
}

func TestMultiWriter_Write_Zero(t *testing.T) {
	var buf bytes.Buffer
	wr := MultiWriter{}

	// register the buffer as a writer
	wr.Register(&buf)

	// write a sample buffer to the writer
	expected := []byte("hello world")
	n, err := wr.Write(expected)
	if err != nil {
		t.Fatalf("failed to write to buffer: %v", err)
	}
	if n != len(expected) {
		t.Fatalf("failed to write to buffer: expected %d bytes, got %d", len(expected), n)
	}

	// verify the buffer contains the correct data
	if !bytes.Equal(buf.Bytes(), expected) {
		t.Fatalf("buffer should have been %v, has: %v", expected, buf.Bytes())
	}
}
