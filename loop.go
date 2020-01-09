package main

// Loop structure.
type Loop struct {
	Watch    Watch
	Commands []*Command
	Run      *Command
}

// NewLoop creates new Loop object with default settings.
func NewLoop() *Loop {
	return &Loop{
		Watch: Watch{
			Dirs:     []string{"."},
			Patterns: []string{"*"},
		},
	}
}

// Loop runs endlessly.
func (l *Loop) Loop() error {
	defer l.Run.Stop()
	for {
		ok, err := l.Execute()
		if err != nil {
			return err
		}
		if l.Run != nil && ok {
			err = l.Run.Start()
			status(err == nil, l.Run.String())
			if err != nil {
				return err
			}
		}
		err = l.Watch.Wait()
		if err != nil {
			return err
		}
		if l.Run != nil && ok {
			_, err = l.Run.Stop()
			if err != nil {
				return err
			}
		}
	}
}

// Execute executes all commands.
func (l *Loop) Execute() (bool, error) {
	for _, c := range l.Commands {
		ok, err := c.Execute()
		status(ok, c.String())
		if !ok {
			return ok, err
		}
	}
	return true, nil
}
