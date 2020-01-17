package main

import (
	"io/ioutil"
	"log"
	"os"
	"testing"
	"time"
)

func TestExecute(t *testing.T) {
	defer quiet()()
	loop := &Loop{}
	loop.Commands = Commands{
		{Exec: "true"},
	}
	ok, err := loop.Execute()
	if !ok {
		t.Error("not ok")
	}
	if err != nil {
		t.Fatal(err)
	}
}

func TestExecuteFail(t *testing.T) {
	defer quiet()()
	loop := &Loop{}
	loop.Commands = Commands{
		{Exec: "false"},
		{Exec: "true"},
	}
	ok, err := loop.Execute()
	if ok {
		t.Error("ok")
	}
	if err != nil {
		t.Fatal(err)
	}
}

func TestExecuteInvalid(t *testing.T) {
	defer quiet()()
	loop := &Loop{}
	loop.Commands = Commands{
		{Exec: "abcdefghijklmnopqrstuwxyz"},
	}
	ok, err := loop.Execute()
	if ok {
		t.Error("ok")
	}
	if err == nil {
		t.Error("no error")
	}
}

func TestStartStopTerminated(t *testing.T) {
	defer quiet()()
	loop := &Loop{}
	loop.Command = &Command{
		Exec: "true",
	}
	err := loop.Start()
	if err != nil {
		t.Fatal(err)
	}
	_, err = loop.Stop()
	if err != nil {
		t.Fatal(err)
	}
}

func TestStartStopDaemon(t *testing.T) {
	defer quiet()()
	loop := &Loop{}
	loop.Command = &Command{
		Exec: "sleep",
		Args: []string{"1m"},
	}
	err := loop.Start()
	if err != nil {
		t.Fatal(err)
	}
	_, err = loop.Stop()
	if err != nil {
		t.Fatal(err)
	}
}

func TestStartStopInvalid(t *testing.T) {
	defer quiet()()
	loop := &Loop{}
	loop.Command = &Command{
		Exec: "abcdefghijklmnopqrstuwxyz",
	}
	err := loop.Start()
	if err == nil {
		t.Fatal("error expected")
	}
	_, err = loop.Stop()
	if err != nil {
		t.Fatal(err)
	}
}

func TestWatch(t *testing.T) {
	defer quiet()()
	loop := &Loop{
		Include: Patterns{"*"},
	}
	wait := make(chan error)
	go func() {
		for {
			ioutil.WriteFile("test", []byte{}, 0644)
			time.Sleep(1 * time.Millisecond)
		}
	}()
	go func() {
		err := loop.Watch()
		wait <- err
		close(wait)
	}()
	var err error
	select {
	case err = <-wait:
	case <-time.After(1 * time.Second):
		t.Fatal("timeout")
	}
	if err != nil {
		t.Fatal(err)
	}
}

func TestEnv(t *testing.T) {
	defer quiet()()
	loop := &Loop{}
	loop.Commands = Commands{
		{
			Env:  map[string]string{"TEST": "x"},
			Exec: "bash",
			Args: []string{"-uc", "test x=$TEST"},
		},
	}
	ok, err := loop.Execute()
	if !ok {
		t.Fatal("not ok")
	}
	if err != nil {
		t.Fatal(err)
	}
}

func quiet() func() {
	null, _ := os.Open(os.DevNull)
	sout := os.Stdout
	serr := os.Stderr
	os.Stdout = null
	os.Stderr = null
	log.SetOutput(null)
	return func() {
		defer null.Close()
		os.Stdout = sout
		os.Stderr = serr
		log.SetOutput(os.Stderr)
	}
}
