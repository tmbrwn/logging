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

import "github.com/tmbrwn/logging"

func main() {
	logging.DateTimeFormat = "3:45pm"

	log := logging.Logger{
		EnableDebug: true,
	}

	log.Print().
		Str("service", "api").
		Msg("initialized")
    // {"time":"2:22pm","service":"api","message":"initialized"}

	log.Debug().
		Msg("attempting request")
    // {"time":"2:22pm","message":"attempting request"}

	err := errors.New("couldn't read from file")
	log.Print().
		Err(err).
		Msg("read")
    // {"time":"2:22pm","pid":12345,"error":"file does not exist","message":"could not read file"}
}
```

## API
| Logger Method | Usage                                                                       |
|---------------|-----------------------------------------------------------------------------|
| `Print`       | Messages that a user should be concerned with                               |
| `Debug`       | Messages that are too verbose or too numerous to be useful unless debugging |

| Chained Method                   | Description                                                                      |
|----------------------------------|----------------------------------------------------------------------------------|
| `Str(key, val string)`           | Tags log with a `string`                                                         |
| `Int(key string, val int)`       | Tags log with an `int`                                                           |
| `Float(key string, val float64)` | Tags log with a `float64`                                                        |
| `Err(err error)`                 | Adds an "error" tag to the log with the given `error`                            |
| `Msg(msg string)`                | Adds the given (`string`) "message" to the log. Must be called last in the chain |
