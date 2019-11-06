package main

import (
	"log"
	"os"
	"testing"
)

func TestExecute(t *testing.T) {
	defer quiet()()
	l := NewLoop()
	l.Commands = [][]string{{"true"}}
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
	l.Commands = [][]string{{"false", "true"}}
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
	l.Commands = [][]string{{"abcdefghijklmnopqrstuwxyz"}}
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
	l.Run = []string{"true"}
	err := l.Start()
	if err != nil {
		t.Fatal(err)
	}
	err = l.Stop()
	if err != nil {
		t.Fatal(err)
	}
}

func TestStartStopDaemon(t *testing.T) {
	defer quiet()()
	l := NewLoop()
	l.Run = []string{"sleep", "1m"}
	err := l.Start()
	if err != nil {
		t.Fatal(err)
	}
	err = l.Stop()
	if err != nil {
		t.Fatal(err)
	}
}

func TestStartStopInvalid(t *testing.T) {
	defer quiet()()
	l := NewLoop()
	l.Run = []string{"abcdefghijklmnopqrstuwxyz"}
	err := l.Start()
	if err == nil {
		t.Fatal("error expected")
	}
	err = l.Stop()
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
