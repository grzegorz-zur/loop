package main

import (
	"github.com/fsnotify/fsnotify"
	"io/ioutil"
	"path/filepath"
)

// Loop structure.
type Loop struct {
	Include  Patterns
	Exclude  Patterns
	Commands `json:"execute"`
	*Command `json:"run"`
}

// Loop loops endlessly.
func (loop *Loop) Loop() error {
	defer loop.Stop()
	for {
		ok, err := loop.Execute()
		if err != nil {
			return err
		}
		if ok && loop.Command != nil {
			err = loop.Start()
			status(err == nil, loop.String())
			if err != nil {
				return err
			}
		}
		err = loop.Watch()
		if err != nil {
			return err
		}
		if ok && loop.Command != nil {
			_, err = loop.Stop()
			if err != nil {
				return err
			}
		}
	}
}

// Watch waits for file changes.
func (loop *Loop) Watch() error {
	watch, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}
	defer watch.Close()
	err = loop.watch(watch, ".")
	if err != nil {
		return err
	}
	for {
		select {
		case event := <-watch.Events:
			name := filepath.Base(event.Name)
			match, err := loop.match(name)
			if err != nil {
				return err
			}
			if match {
				return nil
			}
		case err := <-watch.Errors:
			return err
		}
	}
}

func (loop *Loop) watch(watch *fsnotify.Watcher, path string) error {
	err := watch.Add(path)
	if err != nil {
		return err
	}
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return err
	}
	for _, file := range files {
		if file.IsDir() {
			name := file.Name()
			match, err := loop.match(name)
			if err != nil {
				return err
			}
			if match {
				subpath := filepath.Join(path, name)
				err = loop.watch(watch, subpath)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func (loop *Loop) match(name string) (match bool, err error) {
	match, err = loop.Include.Match(name)
	if !match {
		return false, err
	}
	match, err = loop.Exclude.Match(name)
	if match {
		return false, err
	}
	return true, nil
}
