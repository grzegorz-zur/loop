package main

// Commands represents list of commands.
type Commands []*Command

// Execute executes all commands unless one fails.
func (cmds Commands) Execute() (ok bool, err error) {
	for _, cmd := range cmds {
		ok, err := cmd.Run()
		status(ok, cmd.String())
		if !ok {
			return false, err
		}
	}
	return true, nil
}
