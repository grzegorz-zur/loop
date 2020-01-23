package main

import (
	"path/filepath"
)

// Patterns represents list of file name patterns.
type Patterns []string

// Match checks if name matches any pattern.
func (patterns Patterns) Match(name string) (match bool, err error) {
	for _, pattern := range patterns {
		match, err := filepath.Match(pattern, name)
		if err != nil {
			return false, err
		}
		if match {
			return true, nil
		}
	}
	return false, nil
}
