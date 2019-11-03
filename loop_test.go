package main

import (
	"testing"
)

func TestStartStopTerminated(t *testing.T) {
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

func TestStartStowWrongCommand(t *testing.T) {
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
