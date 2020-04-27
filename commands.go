package main

import (
	"fmt"
)

// Commands represents list of commands.
type Commands []*Command

// Execute executes all commands unless one fails.
func (cmds Commands) Execute() (ok bool, err error) {
	for _, cmd := range cmds {
		ok, err := cmd.Run()
		status(ok, cmd.String())
		if err != nil {
			return false, fmt.Errorf("error running command %s: %w", cmd, err)
		}
		if !ok {
			return false, nil
		}
	}
	return true, nil
}
