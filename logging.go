package logging

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/jonboulle/clockwork"
)

// Global logging properties
var Output io.Writer = os.Stdout     // Output to use for logging
var Pretty = false                   // Pretty print logs
var Clock = clockwork.NewRealClock() // Clock for log timestamps
var DateTimeFormat = time.RFC3339    // Date format for timestamps

// colors for pretty printing
var (
	soft   = color.New(color.FgBlue)
	bright = color.New(color.FgHiWhite)
	red    = color.New(color.FgRed)
)

// global default DefaultLogger
var DefaultLogger = Logger{}

func Print() *log {
	return DefaultLogger.Print()
}

func Debug() *log {
	return DefaultLogger.Debug()
}

// Logger is used to log messages. Multiple loggers could be used within a project
// if fine-grained control over debug levels is desired
type Logger struct {
	EnableDebug bool
}

// Print will always print the resulting message to the output
func (l *Logger) Print() *log {
	return &log{
		logger:  l,
		isDebug: false,
	}
}

// Debug only prints the resulting message to the output if the EnableDebug
// is set to true on the Logger
func (l *Logger) Debug() *log {
	return &log{
		logger:  l,
		isDebug: true,
	}
}

type log struct {
	logger  *Logger
	isDebug bool
	buffer  []ctx
}

type ctx struct {
	key      string
	val      string
	isString bool
}

// Str attaches a string tag to the log
func (l *log) Str(key, val string) *log {
	newCtx := ctx{
		key:      key,
		val:      val,
		isString: true,
	}
	l.buffer = append(l.buffer, newCtx)
	return l
}

// Int attaches an int tag to the log
func (l *log) Int(key string, val int) *log {
	newCtx := ctx{
		key: key,
		val: fmt.Sprintf("%d", val),
	}
	l.buffer = append(l.buffer, newCtx)
	return l
}

// Float attaches a float tag to the log
func (l *log) Float(key string, val float64) *log {
	newCtx := ctx{
		key: key,
		val: strconv.FormatFloat(val, 'f', -1, 64),
	}
	l.buffer = append(l.buffer, newCtx)
	return l
}

// Err attaches an error tag to the log
func (l *log) Err(err error) *log {
	newCtx := ctx{
		key:      "error",
		val:      err.Error(),
		isString: true,
	}
	l.buffer = append(l.buffer, newCtx)
	return l
}

// sub-loggers should really be implemented with a separate chain of context types which are only turned
// into a Logger at the very end:
// logger.With().Str("tag", "value").Logger().Print()...
// func (l *log) Logger() Logger {
// 	return Logger{
// 		EnableDebug: l.logger.EnableDebug,
// 		Output:      l.logger.Output,
// 		Pretty:      l.logger.Pretty,

// 		currentLog: l,
// 	}
// }

// Msg creates a message with the given content and with any tags previously specified
func (l *log) Msg(msg string) {
	l.msgf(msg)
}

// Msgf creates a message with the given format and content and with any tags previously specified
func (l *log) Msgf(format string, args ...interface{}) {
	l.msgf(format, args...)
}

func (l *log) msgf(format string, args ...interface{}) {
	if l.isDebug && !l.logger.EnableDebug {
		return
	}

	msgCtx := ctx{
		key:      "message",
		val:      fmt.Sprintf(format, args...),
		isString: true,
	}
	l.buffer = append(l.buffer, msgCtx)

	if Pretty {
		printLogPretty(l.buffer)
	} else {
		printLog(l.buffer)
	}
}

func printLog(buffer []ctx) {
	now := Clock.Now().Format(DateTimeFormat)
	formattedCtx := []string{fmt.Sprintf(`"time":"%s"`, now)}
	for _, item := range buffer {
		var formattedItem string
		if item.isString {
			formattedItem = fmt.Sprintf(`"%s":"%s"`, item.key, item.val)
		} else {
			formattedItem = fmt.Sprintf(`"%s":%s`, item.key, item.val)
		}
		formattedCtx = append(formattedCtx, formattedItem)
	}

	fmt.Fprintf(Output, "{%s}\n", strings.Join(formattedCtx, ","))
}

func printLogPretty(buffer []ctx) {
	var msg string
	var tags []string
	for _, item := range buffer {
		if item.key == "message" {
			msg = item.val
			continue
		}

		if item.key == "error" {
			tags = append(tags, red.Sprintf("error=%s", item.val))
			continue
		}

		tagKey := soft.Sprintf("%s=", item.key)
		tagVal := bright.Sprint(item.val)
		tags = append(tags, tagKey+tagVal)
	}

	if msg == "" {
		msg = "_"
	}

	now := Clock.Now().Format(DateTimeFormat)

	logComponents := []string{soft.Sprint(now), bright.Sprint(msg)}
	logComponents = append(logComponents, tags...)
	fmt.Fprintln(Output, strings.Join(logComponents, " "))
}
