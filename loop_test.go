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
	l := NewLoop()
	l.Commands = []*Command{{Exec: "true"}}
	ok, err := l.Execute()
	if !ok {
		t.Error("not ok")
	}
	if err != nil {
		t.Fatal(err)
	}
}

func TestExecuteFail(t *testing.T) {
	defer quiet()()
	l := NewLoop()
	l.Commands = []*Command{{Exec: "false"}, {Exec: "true"}}
	ok, err := l.Execute()
	if ok {
		t.Error("ok")
	}
	if err != nil {
		t.Fatal(err)
	}
}

func TestExecuteInvalid(t *testing.T) {
	defer quiet()()
	l := NewLoop()
	l.Commands = []*Command{{Exec: "abcdefghijklmnopqrstuwxyz"}}
	ok, err := l.Execute()
	if ok {
		t.Error("ok")
	}
	if err == nil {
		t.Error("no error")
	}
}

func TestStartStopTerminated(t *testing.T) {
	defer quiet()()
	l := NewLoop()
	l.Run = &Command{Exec: "true"}
	err := l.Run.Start()
	if err != nil {
		t.Fatal(err)
	}
	_, err = l.Run.Stop()
	if err != nil {
		t.Fatal(err)
	}
}

func TestStartStopDaemon(t *testing.T) {
	defer quiet()()
	l := NewLoop()
	l.Run = &Command{Exec: "sleep", Args: []string{"1m"}}
	err := l.Run.Start()
	if err != nil {
		t.Fatal(err)
	}
	_, err = l.Run.Stop()
	if err != nil {
		t.Fatal(err)
	}
}

func TestStartStopInvalid(t *testing.T) {
	defer quiet()()
	l := NewLoop()
	l.Run = &Command{Exec: "abcdefghijklmnopqrstuwxyz"}
	err := l.Run.Start()
	if err == nil {
		t.Fatal("error expected")
	}
	_, err = l.Run.Stop()
	if err != nil {
		t.Fatal(err)
	}
}

func TestWait(t *testing.T) {
	defer quiet()()
	l := NewLoop()
	w := make(chan error)
	go func() {
		for {
			ioutil.WriteFile("test", []byte{}, 0644)
			time.Sleep(1 * time.Millisecond)
		}
	}()
	go func() {
		w <- l.Watch.Wait()
		close(w)
	}()
	var err error
	select {
	case err = <-w:
	case <-time.After(1 * time.Second):
		t.Fatal("timeout")
	}
	if err != nil {
		t.Fatal(err)
	}
}

func TestEnv(t *testing.T) {
	defer quiet()()
	l := NewLoop()
	l.Commands = []*Command{
		{
			Env:  map[string]string{"TEST": "x"},
			Exec: "bash",
			Args: []string{"-uc", "test x=$TEST"},
		},
	}
	ok, err := l.Execute()
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
