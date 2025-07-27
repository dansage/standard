// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2025 Daniel Sage
// SPDX-FileType: SOURCE

package util

import (
	"io"
	"sync"
)

// MultiWriter is a custom io.Writer that directs all output to all registered writers.
type MultiWriter struct {
	// writers is the collection of registered writers.
	writers []io.Writer

	// mutex ensures safe access to the collection of registered writers.
	mutex sync.RWMutex
}

// NewMultiWriter initializes a new log writer with the specified writers pre-registered.
func NewMultiWriter(writers ...io.Writer) *MultiWriter {
	return &MultiWriter{
		writers: writers,
	}
}

// Register the specified writer for all future write operations.
func (w *MultiWriter) Register(wr io.Writer) {
	// ensure no write operations occur while changing the registered writers
	w.mutex.Lock()
	defer w.mutex.Unlock()

	// loop through the registered writers
	for _, writer := range w.writers {
		// check if the writer is already registered
		if writer == wr {
			return
		}
	}

	// add the writer to the collection
	w.writers = append(w.writers, wr)
}

// Unregister the specified writer for all future write operations.
func (w *MultiWriter) Unregister(wr io.Writer) {
	// ensure no write operations occur while changing the registered writers
	w.mutex.Lock()
	defer w.mutex.Unlock()

	// loop through the registered writers
	for i, writer := range w.writers {
		// check if the writer is identical
		if writer == wr {
			// remove the writer from the collection
			w.writers = append(w.writers[:i], w.writers[i+1:]...)
			return
		}
	}
}

// Write writes len(p) from p to all registered writers.
func (w *MultiWriter) Write(p []byte) (n int, err error) {
	// ensure the registered writers do not change
	w.mutex.RLock()
	defer w.mutex.RUnlock()

	// loop through the registered writers
	for _, writer := range w.writers {
		// write the data to the stream
		if n, err := writer.Write(p); err != nil {
			return n, err
		}
	}

	return n, nil
}
