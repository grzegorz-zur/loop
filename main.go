package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

const (
	config = ".loop.json"
)

const (
	red   = "\033[31m"
	green = "\033[32m"
	reset = "\033[39;49m"
)

func main() {
	loop := &Loop{}
	data, err := ioutil.ReadFile(config)
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

func status(ok bool, text string) {
	color := green
	if !ok {
		color = red
	}
	log.Println(color + text + reset)
}
