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

### Defaults

```json
{
	"watch" : {
		"directories" : [
			"."
		],
		"patterns" : [
			"*"
		]
	}
}
```

### Example

```json
{
	"watch" : {
		"directories" : [
			".",
			"cmd/app"
		],
		"patterns" : [
			"*.go",
			"go.mod"
		]
	},
	"commands": [
		[ "go", "fmt", "./..." ],
		[ "go", "vet" , "./..."],
		[ "go", "build" , "./cmd/app"]
	],
	"run" : [ "./app", "-addr=:8443" ]
}
```

## Usage

```sh
loop
```
