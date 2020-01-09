package main

import (
	"github.com/fsnotify/fsnotify"
	"path"
)

// Watch structure.
type Watch struct {
	Dirs     []string
	Patterns []string
}

// Wait waits for file changes.
func (w *Watch) Wait() error {
	fsw, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}
	defer fsw.Close()
	for _, d := range w.Dirs {
		fsw.Add(d)
	}
	for {
		select {
		case e := <-fsw.Events:
			m, err := w.match(e)
			if err != nil {
				return err
			}
			if m {
				return nil
			}
		case err := <-fsw.Errors:
			return err
		}
	}
}

func (w *Watch) match(e fsnotify.Event) (bool, error) {
	_, f := path.Split(e.Name)
	for _, p := range w.Patterns {
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
