package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
)

const (
	Config = ".loop.json"
)

func main() {
	var loop Loop
	data, err := ioutil.ReadFile(Config)
	if err != nil {
		log.Fatal(err)
	}
	err = json.Unmarshal(data, &loop)
	if err != nil {
		log.Fatal(err)
	}
	loop.Loop()
}

type Loop struct {
	Commands [][]string
	Run      []string
}

func (l *Loop) Loop() {
	for {
		l.Execute()
		l.Start()
		l.Wait()
		l.Stop()
	}
}

func (l *Loop) Execute() {
	for _, c := range l.Commands {
		log.Println(strings.Join(c, " "))
		cmd := exec.Command(c[0], c[1:]...)
		data, err := cmd.CombinedOutput()
		os.Stdout.Write(data)
		if err != nil {
			break
		}
	}
}

func (l *Loop) Start() {
}

func (l *Loop) Stop() {
}

func (l *Loop) Wait() {
}
