// Package storage defines storage interfaces and common errors.
package storage

import "errors"

var (
	// ErrURLNotFound is returned when a URL is not found in storage.
	ErrURLNotFound = errors.New("url not found")
	// ErrURLExists is returned when a URL with the given alias already exists.
	ErrURLExists = errors.New("url exists")
)
