# Loop
 
Minimalistic task runner.

## Installation

### Install

```
go get github.com/grzegorz-zur/loop
```

### Update

```
go get -u github.com/grzegorz-zur/loop
```

### Uninstall

```
go clean -i -r -cache -modcache github.com/grzegorz-zur/loop
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

