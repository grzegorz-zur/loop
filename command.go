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

func (cmd *Command) String() string {
	text := strings.Builder{}
	if cmd.Dir != "" {
		text.WriteString(cmd.Dir)
		text.WriteString(" ")
	}
	if len(cmd.Env) != 0 {
		for name, value := range cmd.Env {
			text.WriteString(name)
			text.WriteString("=")
			text.WriteString(value)
		}
		text.WriteString(" ")
	}
	text.WriteString(cmd.Exec)
	if len(cmd.Args) != 0 {
		text.WriteString(" ")
		args := strings.Join(cmd.Args, " ")
		text.WriteString(args)
	}
	return text.String()
}

// Run starts and waits for command to finish.
func (cmd *Command) Run() (ok bool, err error) {
	err = cmd.Start()
	if err != nil {
		return false, fmt.Errorf("error running command %s: %w", cmd, err)
	}
	return cmd.Wait()
}

// Start starts the command.
func (cmd *Command) Start() error {
	exe, err := exec.LookPath(cmd.Exec)
	if err != nil {
		return fmt.Errorf("error looking for path of %s: %w", cmd.Exec, err)
	}
	args := make([]string, 0, len(cmd.Args)+1)
	args = append(args, cmd.Exec)
	args = append(args, cmd.Args...)
	attr := cmd.attr()
	cmd.process, err = os.StartProcess(exe, args, attr)
	if err != nil {
		return fmt.Errorf("error starting command %s: %w", cmd, err)
	}
	return nil
}

// Stop stops the running command.
func (cmd *Command) Stop() (ok bool, err error) {
	if cmd.process == nil {
		return true, nil
	}
	err = cmd.process.Signal(syscall.SIGTERM)
	if err != nil {
		return false, fmt.Errorf("error stoping process %d: %w", cmd.process.Pid, err)
	}
	ok, err = cmd.Wait()
	if err != nil {
		return false, fmt.Errorf("error waiting for process %d: %w", cmd.process.Pid, err)
	}
	return ok, nil
}

// Wait waits for process end returns result.
func (cmd *Command) Wait() (ok bool, err error) {
	if cmd.process == nil {
		return true, nil
	}
	defer func() {
		cmd.process = nil
	}()
	state, err := cmd.process.Wait()
	if err != nil {
		return false, fmt.Errorf("error waiting for process %d: %w", cmd.process.Pid, err)
	}
	ok = state.Success()
	return ok, nil
}

func (cmd *Command) attr() *os.ProcAttr {
	env := []string{}
	for name, value := range cmd.Env {
		entry := fmt.Sprintf("%s=%s", name, value)
		env = append(env, entry)
	}
	env = append(env, os.Environ()...)
	files := []*os.File{os.Stdin, os.Stdout, os.Stderr}
	return &os.ProcAttr{
		Dir:   cmd.Dir,
		Env:   env,
		Files: files,
	}
}
