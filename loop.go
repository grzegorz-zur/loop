package main

import (
	"fmt"
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
			return fmt.Errorf("error executing loop: %w", err)
		}
		if ok && loop.Command != nil {
			err = loop.Start()
			status(err == nil, loop.String())
			if err != nil {
				return fmt.Errorf("error starting command: %w", err)
			}
		}
		err = loop.Watch()
		if err != nil {
			return fmt.Errorf("error watching changes: %w", err)
		}
		if ok && loop.Command != nil {
			_, err = loop.Stop()
			if err != nil {
				return fmt.Errorf("error stopping command: %w", err)
			}
		}
	}
}

// Watch waits for file changes.
func (loop *Loop) Watch() error {
	watch, err := fsnotify.NewWatcher()
	if err != nil {
		return fmt.Errorf("error creating watcher: %w", err)
	}
	defer watch.Close()
	err = loop.watch(watch, ".")
	if err != nil {
		return fmt.Errorf("error watching: %w", err)
	}
	for {
		select {
		case event := <-watch.Events:
			name := filepath.Base(event.Name)
			match, err := loop.match(name)
			if err != nil {
				return fmt.Errorf("error matching name %s: %w", name, err)
			}
			if match {
				return nil
			}
		case err := <-watch.Errors:
			return fmt.Errorf("error from watch: %w", err)
		}
	}
}

func (loop *Loop) watch(watch *fsnotify.Watcher, path string) error {
	err := watch.Add(path)
	if err != nil {
		return fmt.Errorf("error adding watch path %s: %w", path, err)
	}
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return fmt.Errorf("error reading directory %s: %w", path, err)
	}
	for _, file := range files {
		if file.IsDir() {
			name := file.Name()
			match, err := loop.match(name)
			if err != nil {
				return fmt.Errorf("error matching name %s: %w", name, err)
			}
			if match {
				subpath := filepath.Join(path, name)
				err = loop.watch(watch, subpath)
				if err != nil {
					return fmt.Errorf("error watchng subpath %s: %w", subpath, err)
				}
			}
		}
	}
	return nil
}

func (loop *Loop) match(name string) (match bool, err error) {
	match, err = loop.Include.Match(name)
	if err != nil {
		return false, fmt.Errorf("error matching included name %s: %w", name, err)
	}
	if !match {
		return false, nil
	}
	match, err = loop.Exclude.Match(name)
	if err != nil {
		return false, fmt.Errorf("error matching excluded name %s: %w", name, err)
	}
	if match {
		return false, nil
	}
	return true, nil
}
