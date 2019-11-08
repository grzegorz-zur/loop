# Loop
 
Minimalistic task runner.

## Status

[![CircleCI](https://circleci.com/gh/grzegorz-zur/loop.svg?style=svg)](https://circleci.com/gh/grzegorz-zur/loop)
[![Codacy Badge](https://api.codacy.com/project/badge/Grade/936a0d49701e4517813c9def722b21dd)](https://www.codacy.com/manual/grzegorz.zur/loop?utm_source=github.com&amp;utm_medium=referral&amp;utm_content=grzegorz-zur/loop&amp;utm_campaign=Badge_Grade)
[![Codacy Badge](https://api.codacy.com/project/badge/Coverage/936a0d49701e4517813c9def722b21dd)](https://www.codacy.com/manual/grzegorz.zur/loop?utm_source=github.com&utm_medium=referral&utm_content=grzegorz-zur/loop&utm_campaign=Badge_Coverage)

## Installation

To install or update run the following command.

```
go get -u github.com/grzegorz-zur/loop
```

## Configuration

Create `.loop.json`.

### Defaults

```
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

```
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

```
loop
```

