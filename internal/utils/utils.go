// Package utils provides internal utility functions for the Forge framework.
// This package is not importable by external projects.
package utils

// Internal utility functions will be added as needed

// Min returns the smaller of two integers.
func Min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
