package logging

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"runtime"
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

///////////////////////////////////////////////////////////////////////////////
// PACKAGE-LEVEL LOGGING //////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////

// global default logger
var DefaultLogger = Logger{}

func Tag(key string, val interface{}) *log {
	return DefaultLogger.Tag(key, val)
}

func Err(err error) *log {
	return DefaultLogger.Err(err)
}

// Print logs a message to the output
func Print(a ...interface{}) {
	printLog(DefaultLogger.tags, fmt.Sprint(a...), DefaultLogger.EnableDebug)
}

// Printf logs a message to the output
func Printf(format string, a ...interface{}) {
	printLog(DefaultLogger.tags, fmt.Sprintf(format, a...), DefaultLogger.EnableDebug)
}

// Debug only logs a message to the output if the EnableDebug
// is set to true on the Logger
func Debug(a ...interface{}) {
	if DefaultLogger.EnableDebug {
		printLog(DefaultLogger.tags, fmt.Sprint(a...), DefaultLogger.EnableDebug)
	}
}

// Debugf only logs a message to the output if the EnableDebug
// is set to true on the Logger
func Debugf(format string, a ...interface{}) {
	if DefaultLogger.EnableDebug {
		printLog(DefaultLogger.tags, fmt.Sprintf(format, a...), DefaultLogger.EnableDebug)
	}
}

///////////////////////////////////////////////////////////////////////////////
// LOGGER-LEVEL LOGGING ///////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////

// Logger is used to log messages. Multiple loggers could be used within a project
// if fine-grained control over debug levels is desired
type Logger struct {
	EnableDebug bool
	tags        []tag
}

// Tag associates the given key with the given value in
// the resulting message's tags
func (l *Logger) Tag(key string, val interface{}) *log {
	return &log{
		enableDebug: l.EnableDebug,
		tags:        append(l.tags, tag{key, val}),
	}
}

// Err is short for Tag("error", err)
func (l *Logger) Err(err error) *log {
	return &log{
		enableDebug: l.EnableDebug,
		tags:        append(l.tags, tag{"error", err}),
	}
}

// Print logs a message to the output
func (l *Logger) Print(a ...interface{}) {
	printLog(l.tags, fmt.Sprint(a...), l.EnableDebug)
}

// Printf logs a message to the output
func (l *Logger) Printf(format string, a ...interface{}) {
	printLog(l.tags, fmt.Sprintf(format, a...), l.EnableDebug)
}

// Debug only logs a message to the output if the EnableDebug
// is set to true on the Logger
func (l *Logger) Debug(a ...interface{}) {
	if l.EnableDebug {
		printLog(l.tags, fmt.Sprint(a...), l.EnableDebug)
	}
}

// Debugf only logs a message to the output if the EnableDebug
// is set to true on the Logger
func (l *Logger) Debugf(format string, a ...interface{}) {
	if l.EnableDebug {
		printLog(l.tags, fmt.Sprintf(format, a...), l.EnableDebug)
	}
}

type log struct {
	enableDebug bool
	tags        []tag
}

type tag struct {
	key string
	val interface{}
}

// Logger returns a logger with all the tags defined so far
func (l *log) Logger() *Logger {
	return &Logger{
		EnableDebug: l.enableDebug,
		tags:        l.tags,
	}
}

// Tag associates the given key with the given value in
// the resulting message's tags
func (l *log) Tag(key string, val interface{}) *log {
	l.tags = append(l.tags, tag{key, val})
	return l
}

// Err is short for Tag("error", err)
func (l *log) Err(err error) *log {
	l.tags = append(l.tags, tag{"error", err})
	return l
}

// Print logs a message to the output
func (l *log) Print(a ...interface{}) {
	printLog(l.tags, fmt.Sprint(a...), l.enableDebug)
}

// Printf logs a message to the output
func (l *log) Printf(format string, a ...interface{}) {
	printLog(l.tags, fmt.Sprintf(format, a...), l.enableDebug)
}

// Debug only logs a message to the output if the EnableDebug
// is set to true on the Logger
func (l *log) Debug(a ...interface{}) {
	if l.enableDebug {
		printLog(l.tags, fmt.Sprint(a...), l.enableDebug)
	}
}

// Debugf only logs a message to the output if the EnableDebug
// is set to true on the Logger
func (l *log) Debugf(format string, a ...interface{}) {
	if l.enableDebug {
		printLog(l.tags, fmt.Sprintf(format, a...), l.enableDebug)
	}
}

///////////////////////////////////////////////////////////////////////////////
// OUTPUT FORMATTING //////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////

func printLog(tags []tag, msg string, debugEnabled bool) {
	if Pretty {
		printLogPretty(tags, msg, debugEnabled)
	} else {
		printLogJSON(tags, msg, debugEnabled)
	}
}

func printLogJSON(tags []tag, msg string, debugEnabled bool) {
	now := Clock.Now().Format(DateTimeFormat)

	logComponents := []string{
		fmt.Sprintf(`"time":"%s"`, now),
		fmt.Sprintf(`"message":"%s"`, msg),
	}

	if debugEnabled {
		caller := "unknown"
		// skip 3: printLogJSON, printLog, and public caller thereof
		_, fileName, lineNum, ok := runtime.Caller(3)
		if ok {
			caller = fmt.Sprintf("%s:%d", fileName, lineNum)
		}

		logComponents = append(logComponents, fmt.Sprintf(`"caller":"%s"`, caller))
	}

	for _, t := range tags {
		val := `"<?>"`
		switch t.val.(type) {
		case error:
			val = fmt.Sprintf(`"%s"`, t.val)
		default:
			marshaled, err := json.Marshal(t.val)
			if err == nil {
				val = string(marshaled)
			}
		}

		logComponents = append(logComponents, fmt.Sprintf(`"%s":%s`, t.key, val))
	}

	fmt.Fprintf(Output, "{%s}\n", strings.Join(logComponents, ","))
}

func printLogPretty(tags []tag, msg string, debugEnabled bool) {
	if msg == "" {
		msg = "_"
	}

	now := Clock.Now().Format(DateTimeFormat)

	logComponents := []string{
		soft.Sprint(now),
		bright.Sprint(msg),
	}

	if debugEnabled {
		caller := "unknown"
		// skip 3: printLogJSON, printLog, and public caller thereof
		_, fileName, lineNum, ok := runtime.Caller(3)
		if ok {
			caller = fmt.Sprintf("%s:%d", fileName, lineNum)
		}

		logComponents = append(logComponents, soft.Sprintf("(%s)", caller))
	}

	for _, t := range tags {
		var val string
		switch t.val.(type) {
		case error:
			val = red.Sprint(t.val)
		default:
			val = bright.Sprint(t.val)
		}

		key := soft.Sprintf("%s=", t.key)
		logComponents = append(logComponents, key+val)
	}

	fmt.Fprintln(Output, strings.Join(logComponents, " "))
}
