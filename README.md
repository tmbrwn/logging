# logging
A simple structured-logging library with a chaining API

![example](example.png)

The package attempts to provide a logging solution with the following properties:
- Simplest possible API
- Toggleable structured/pretty logging for production/dev work respectively
- Debug-level logging toggleable at an arbitarily fine-grained level

This package takes inspiration from several sources, including other logging packages like [zerolog](https://github.com/rs/zerolog), and from Dave Cheney's discussion of logging [on his blog](https://dave.cheney.net/2015/11/05/lets-talk-about-logging).

_Anyone is welcome to use this package, but it is primarily designed for my own use, so I can't promise that any Issues will be resolved or PRs merged._

## Getting started
```go get github.com/tmbrwn/logging```

## Using the package
Using the package starts with creating a `logging.Logger`. Loggers can be shared between as many or as few parts of an application as desired, and Debug level logging can be configured wherever a Logger is kept.

A package-level default logger has been defined along with package-level logging functions for convenience.

Global settings among Loggers have the following defaults:
```go
var Output io.Writer = os.Stdout     // Output to use for logging
var Pretty = false                   // Pretty print logs
var Clock = clockwork.NewRealClock() // Clock for timestamps
var DateTimeFormat = time.RFC3339    // Date format for timestamps
```

These parameters should be adjusted once on initialization. They have no thread-safe wrappers for adjustment mid-execution.

## Examples
```go
package main

import (
	"errors"

	"github.com/tmbrwn/logging"
)

func main() {
	logging.DateTimeFormat = "3:04pm"

	log := logging.Logger{
		EnableDebug: true,
	}

	log.Print("starting up")
	// {"time":"2:22pm","message":"starting up","caller":"main.go:20"}

	l := log.Tag("request-id", "abc123").Logger()

	l.Debug("received request")
	// {"time":"2:22pm","message":"received request","caller":"main.go:17","request-id":"abc123",}

	err := errors.New("file does not exist")
	l.Err(err).Print("could not read file")
	// {"time":"2:22pm","message":"could not read file","caller":"main.go:26","request-id":"abc123","error":"file does not exist"}
}
```

## API
The following methods can be used at the package level or on an individual `Logger` value.


| Method           | Usage                                                                       |
|------------------|-----------------------------------------------------------------------------|
| `Print`/`Printf` | Messages that a user should be concerned with                               |
| `Debug`/`Debugf` | Messages that are too verbose or too numerous to be useful unless debugging |
| `Tag`            | Give a key-value tag to a message                                           |
| `Err`            | Short for `Tag("error", err)`                                               |
| `Logger`         | Save the tags defined so far in a separate logger value                     |
