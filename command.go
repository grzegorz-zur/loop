package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"syscall"
)

// Command structure.
type Command struct {
	Dir     string
	Env     map[string]string
	Exec    string
	Args    []string
	process *os.Process
}

func (c *Command) String() string {
	text := strings.Builder{}
	if c.Dir != "" {
		text.WriteString(c.Dir)
		text.WriteString(" ")
	}
	if len(c.Env) != 0 {
		for k, v := range c.Env {
			text.WriteString(k)
			text.WriteString("=")
			text.WriteString(v)
		}
		text.WriteString(" ")
	}
	text.WriteString(c.Exec)
	if len(c.Args) != 0 {
		text.WriteString(" ")
		text.WriteString(strings.Join(c.Args, " "))
	}
	return text.String()
}

// Execute executes sequance of commands from list.
//
// Returns true when all commands succeeded and false otherwise.
func (c *Command) Execute() (bool, error) {
	err := c.Start()
	if err != nil {
		return false, err
	}
	return c.Wait()
}

// Start starts the command set in `run` fields.
func (c *Command) Start() error {
	attr := c.attr()
	exe, err := exec.LookPath(c.Exec)
	if err != nil {
		return err
	}
	args := make([]string, 0, len(c.Args)+1)
	args = append(args, c.Exec)
	args = append(args, c.Args...)
	c.process, err = os.StartProcess(exe, args, attr)
	if err != nil {
		return err
	}
	return nil
}

func (c *Command) attr() *os.ProcAttr {
	env := []string{}
	for k, v := range c.Env {
		e := fmt.Sprintf("%s=%s", k, v)
		env = append(env, e)
	}
	env = append(env, os.Environ()...)
	files := []*os.File{os.Stdin, os.Stdout, os.Stderr}
	return &os.ProcAttr{
		Dir:   c.Dir,
		Env:   env,
		Files: files,
	}
}

// Stop stops the running command.
func (c *Command) Stop() (bool, error) {
	if c.process == nil {
		return true, nil
	}
	err := c.process.Signal(syscall.SIGTERM)
	if err != nil {
		return false, err
	}
	return c.Wait()
}

// Wait waits for process end returns result.
func (c *Command) Wait() (bool, error) {
	if c.process == nil {
		return true, nil
	}
	defer func() {
		c.process = nil
	}()
	state, err := c.process.Wait()
	ok := state.Success()
	return ok, err
}
