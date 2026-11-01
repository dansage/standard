// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2025 Daniel Sage
// SPDX-FileType: SOURCE

package util

import (
	"io"
	"reflect"
	"slices"
	"sync"
)

// MultiWriter is a custom io.Writer that directs all output to all registered writers.
type MultiWriter struct {
	// writers is the collection of registered writers.
	writers []io.Writer

	// mutex ensures safe access to the collection of registered writers.
	mutex sync.Mutex
}

// NewMultiWriter initializes a new multi writer with the specified writers pre-registered. Writers should be provided
// as pointers to ensure they can be reliably identified later on.
func NewMultiWriter(writers ...io.Writer) *MultiWriter {
	w := new(MultiWriter)

	// loop through the specified writers, registering each and ignoring duplicates
	for _, writer := range writers {
		w.Register(writer)
	}

	return w
}

// compareWriter determines if the specified writers are the same, always failing if either writer is not comparable.
func compareWriter(a, b io.Writer) bool {
	// use reflection to determine the types
	ta := reflect.TypeOf(a)
	tb := reflect.TypeOf(b)

	// verify the types are identical and comparable
	if ta != tb || !ta.Comparable() || !tb.Comparable() {
		return false
	}
	return a == b
}

// Register the specified writer for all future write operations. Writers should be provided as pointers to ensure they
// can be reliably identified later on. If the writer is not comparable, it will not be registered.
func (w *MultiWriter) Register(wr io.Writer) {
	// ensure no write operations occur while changing the registered writers
	w.mutex.Lock()
	defer w.mutex.Unlock()

	// verify the writer is not nil
	if wr == nil {
		return
	}
	if v := reflect.ValueOf(wr); v.Kind() == reflect.Pointer && v.IsNil() {
		return
	}

	// verify the writer is a comparable type
	if !reflect.TypeOf(wr).Comparable() {
		return
	}

	// verify the writer is not already registered
	if slices.ContainsFunc(w.writers, func(writer io.Writer) bool {
		return compareWriter(wr, writer)
	}) {
		return
	}

	// add the writer to the collection
	w.writers = append(w.writers, wr)
}

// Unregister the specified writer for all future write operations.
func (w *MultiWriter) Unregister(wr io.Writer) {
	// ensure no write operations occur while changing the registered writers
	w.mutex.Lock()
	defer w.mutex.Unlock()

	// remove the writer, if it is registered
	w.writers = slices.DeleteFunc(w.writers, func(writer io.Writer) bool {
		return compareWriter(wr, writer)
	})
}

// Write writes len(p) from p to all registered writers. If an error occurs while writing, an error is returned
// immediately and all remaining writers are skipped.
func (w *MultiWriter) Write(p []byte) (n int, err error) {
	// ensure no changes occur to the registered writers
	w.mutex.Lock()
	defer w.mutex.Unlock()

	// loop through the registered writers
	for _, writer := range w.writers {
		// write the data to the stream
		n, err = writer.Write(p)
		if err != nil {
			return n, err
		}
		if n < len(p) {
			return n, io.ErrShortWrite
		}
	}

	return len(p), nil
}
