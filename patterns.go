package main

import (
	"fmt"
	"path/filepath"
)

// Patterns represents list of file name patterns.
type Patterns []string

// Match checks if name matches any pattern.
func (patterns Patterns) Match(name string) (match bool, err error) {
	for _, pattern := range patterns {
		match, err := filepath.Match(pattern, name)
		if err != nil {
			return false, fmt.Errorf("error matching pattern %s with name %s: %w", pattern, name, err)
		}
		if match {
			return true, nil
		}
	}
	return false, nil
}
