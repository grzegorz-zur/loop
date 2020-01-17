# Loop
 
Minimalistic task runner.

## Status

[![Actions](https://github.com/grzegorz-zur/loop/workflows/Test/badge.svg)](https://github.com/grzegorz-zur/loop/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/grzegorz-zur/loop)](https://goreportcard.com/report/github.com/grzegorz-zur/loop)
[![codecov](https://codecov.io/gh/grzegorz-zur/loop/branch/master/graph/badge.svg)](https://codecov.io/gh/grzegorz-zur/loop)

## Installation

To install or update run the following command.

```sh
go get -u github.com/grzegorz-zur/loop
```

## Configuration

Create `.loop.json`.

### Example

```json
{
	"include" : [
		"*.go",
		"go.mod"
	],
	"exclude" : [
		".*"
	],
	"execute": [
		{
			"exec": "gofmt",
			"args": [ "-s", "-w", "./..." ]
		},
		{
			"exec": "golint",
			"args": [ "-set_exit_status", "./..." ]
		},
		{
			"exec": "go",
			"args": [ "vet", "./..." ]
		},
		{
			"exec": "go",
			"args": [ "build", "./..." ]
		},
		{
			"exec": "go",
			"args": [ "test", "-timeout", "5s", "./..." ]
		},
		{
			"exec": "go",
			"args": [ "test", "-timeout", "5s", "-race", "-cover", "./..." ],
			"env": {
				"CGO_ENABLED": "1"
			}
		}
	],
	"run" : {
		"exec": "./app",
		"args": [ "-addr=:8443" ]
	}
}
```

## Usage

```sh
loop
```
