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
func (lp *Loop) Loop() error {
	defer lp.Stop()
	for {
		ok, err := lp.Execute()
		if err != nil {
			return err
		}
		if ok && lp.Command != nil {
			err = lp.Start()
			status(err == nil, lp.String())
			if err != nil {
				return err
			}
		}
		err = lp.Watch()
		if err != nil {
			return err
		}
		if ok && lp.Command != nil {
			_, err = lp.Stop()
			if err != nil {
				return err
			}
		}
	}
}

// Watch waits for file changes.
func (lp *Loop) Watch() error {
	watch, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}
	defer watch.Close()
	err = lp.watch(watch, ".")
	if err != nil {
		return err
	}
	for {
		select {
		case event := <-watch.Events:
			name := filepath.Base(event.Name)
			match, err := lp.match(name)
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

func (lp *Loop) watch(watch *fsnotify.Watcher, path string) error {
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
			match, err := lp.match(name)
			if err != nil {
				return err
			}
			if match {
				subpath := filepath.Join(path, name)
				err = lp.watch(watch, subpath)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func (lp *Loop) match(name string) (match bool, err error) {
	match, err = lp.Include.Match(name)
	if !match {
		return false, err
	}
	match, err = lp.Exclude.Match(name)
	if match {
		return false, err
	}
	return true, nil
}
