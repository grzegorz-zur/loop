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
	Config = ".loop.json"
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

type Loop struct {
	Watch    Watch
	Commands [][]string
	Run      []string
	run      *exec.Cmd
}

type Watch struct {
	Directories []string
	Patterns    []string
}

func NewLoop() *Loop {
	return &Loop{
		Watch: Watch{
			Directories: []string{"."},
			Patterns:    []string{"*"},
		},
	}
}

func (l *Loop) Loop() error {
	defer l.Stop()
	for {
		err := l.Execute()
		if err != nil {
			return err
		}
		err = l.Start()
		if err != nil {
			return err
		}
		err = l.Wait()
		if err != nil {
			return err
		}
		err = l.Stop()
		if err != nil {
			return err
		}
	}
}

func (l *Loop) Execute() error {
	for _, c := range l.Commands {
		log.Println(strings.Join(c, " "))
		cmd := exec.Command(c[0], c[1:]...)
		data, err := cmd.CombinedOutput()
		os.Stdout.Write(data)
		if err != nil {
			var exit *exec.ExitError
			if errors.As(err, &exit) {
				break
			} else {
				return err
			}
		}
	}
	return nil
}

func (l *Loop) Start() error {
	if l.Run == nil {
		return nil
	}
	log.Println(strings.Join(l.Run, " "))
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

	err = l.run.Start()
	if err != nil {
		var exit *exec.ExitError
		if !errors.As(err, &exit) {
			return err
		}
	}

	return nil
}

func (l *Loop) Stop() error {
	if l.run == nil || l.run.Process == nil {
		return nil
	}
	err := l.run.Process.Signal(syscall.SIGTERM)
	if err != nil {
		return err
	}
	err = l.run.Wait()
	if err != nil {
		return err
	}
	l.run = nil
	return nil
}

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
