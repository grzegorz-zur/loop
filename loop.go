package main

import (
	"encoding/json"
	"errors"
	"github.com/fsnotify/fsnotify"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"strings"
	"syscall"
)

const (
	Config = ".loop.json" // Configuration file name
	red    = "\033[31m"
	green  = "\033[32m"
	reset  = "\033[39;49m"
)

func main() {
	loop := NewLoop()
	data, err := ioutil.ReadFile(Config)
	if err != nil {
		log.Fatal(err)
	}
	err = json.Unmarshal(data, loop)
	if err != nil {
		log.Fatal(err)
	}
	err = loop.Loop()
	if err != nil {
		log.Fatal(err)
	}
}

// Loop structure.
type Loop struct {
	Watch    Watch
	Commands [][]string
	Run      []string
	run      *exec.Cmd
}

// Watch structure.
type Watch struct {
	Directories []string
	Patterns    []string
}

// NewLoop creates new Loop object with default settings.
func NewLoop() *Loop {
	return &Loop{
		Watch: Watch{
			Directories: []string{"."},
			Patterns:    []string{"*"},
		},
	}
}

// Loop runs endlessly.
func (l *Loop) Loop() error {
	defer l.Stop()
	for {
		ok, err := l.Execute()
		if err != nil {
			return err
		}
		if ok {
			err = l.Start()
			if err != nil {
				return err
			}
		}
		err = l.Wait()
		if err != nil {
			return err
		}
		if ok {
			err = l.Stop()
			if err != nil {
				return err
			}
		}
	}
}

// Execute executes seqance of commands from list.
//
// Returns true when all commands suceeded and false otherwise.
func (l *Loop) Execute() (bool, error) {
	for _, c := range l.Commands {
		text := strings.Join(c, " ")
		cmd := exec.Command(c[0], c[1:]...)
		data, err := cmd.CombinedOutput()
		if err != nil {
			failure(text)
			os.Stdout.Write(data)
			var exit *exec.ExitError
			if errors.As(err, &exit) {
				return false, nil
			}
			return false, err
		}
		success(text)
		os.Stdout.Write(data)
	}
	return true, nil
}

// Start starts the command set in `run` fields.
func (l *Loop) Start() error {
	if l.Run == nil {
		return nil
	}
	l.run = exec.Command(l.Run[0], l.Run[1:]...)

	i, err := l.run.StdinPipe()
	if err != nil {
		return err
	}
	go io.Copy(i, os.Stdin)

	o, err := l.run.StdoutPipe()
	if err != nil {
		return err
	}
	go io.Copy(os.Stdout, o)

	e, err := l.run.StderrPipe()
	if err != nil {
		return err
	}
	go io.Copy(os.Stderr, e)

	text := strings.Join(l.Run, " ")
	err = l.run.Start()
	if err != nil {
		failure(text)
		return err
	}

	success(text)
	return nil
}

// Stop stops the running command.
func (l *Loop) Stop() error {
	if l.run == nil || l.run.Process == nil {
		return nil
	}
	defer func() {
		l.run = nil
	}()
	err := l.run.Process.Signal(syscall.SIGTERM)
	if err != nil {
		return err
	}
	err = l.run.Wait()
	var exit *exec.ExitError
	if err != nil && !errors.As(err, &exit) {
		return err
	}
	return nil
}

// Wait waits for file changes.
func (l *Loop) Wait() error {
	w, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}
	defer w.Close()
	for _, d := range l.Watch.Directories {
		w.Add(d)
	}
	for {
		select {
		case e := <-w.Events:
			m, err := l.match(e)
			if err != nil {
				return err
			}
			if m {
				return nil
			}
		case err := <-w.Errors:
			return err
		}
	}
}

func (l *Loop) match(e fsnotify.Event) (bool, error) {
	_, f := path.Split(e.Name)
	for _, p := range l.Watch.Patterns {
		m, err := path.Match(p, f)
		if err != nil {
			return false, err
		}
		if m {
			return m, nil
		}
	}
	return false, nil
}

func success(text string) {
	log.Println(green + text + reset)
}

func failure(text string) {
	log.Println(red + text + reset)
}
