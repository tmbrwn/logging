package main

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/fatih/color"
	"github.com/jonboulle/clockwork"
)

func TestPrintMsg(t *testing.T) {
	output := &logout{}
	Clock = clockAt("2020-02-02T02:02:02Z")
	Pretty = false
	Output = output
	l := Logger{
		EnableDebug: false,
	}

	l.Print().Msg("hello")

	expectedLog := `{"time":"2020-02-02T02:02:02Z","message":"hello"}` + "\n"
	if !equals(output.content, expectedLog) {
		t.Errorf(`expected "%s" but got "%s"`, expectedLog, output.content)
	}
}

func TestDebugMsgNoDebug(t *testing.T) {
	output := &logout{}
	Clock = clockAt("2020-02-02T02:02:02Z")
	Pretty = false
	Output = output
	l := Logger{
		EnableDebug: false,
	}

	l.Debug().Msg("hello")

	if !equals(output.content, "") {
		t.Errorf(`expected empty output but got "%s"`, output.content)
	}
}

func TestDebugMsg(t *testing.T) {
	output := &logout{}
	Clock = clockAt("2020-02-02T02:02:02Z")
	Pretty = false
	Output = output
	l := Logger{
		EnableDebug: true,
	}

	l.Debug().Msg("hello")

	expectedLog := `{"time":"2020-02-02T02:02:02Z","message":"hello"}` + "\n"
	if !equals(output.content, expectedLog) {
		t.Errorf(`expected "%s" but got "%s"`, expectedLog, output.content)
	}
}

func TestPrintTaggedMsg(t *testing.T) {
	output := &logout{}
	Clock = clockAt("2020-02-02T02:02:02Z")
	Pretty = false
	Output = output
	l := Logger{
		EnableDebug: false,
	}

	l.Print().
		Str("service", "hello-service").
		Msg("hello")

	expectedLog := `{"time":"2020-02-02T02:02:02Z","service":"hello-service","message":"hello"}` + "\n"
	if !equals(output.content, expectedLog) {
		t.Errorf(`expected "%s" but got "%s"`, expectedLog, output.content)
	}
}

func TestPrintMultiTaggedMsg(t *testing.T) {
	output := &logout{}
	Clock = clockAt("2020-02-02T02:02:02Z")
	Pretty = false
	Output = output
	l := Logger{
		EnableDebug: false,
	}

	l.Print().
		Str("service", "hello-service").
		Str("service", "same-tag").
		Msg("hello")

	expectedLog := `{"time":"2020-02-02T02:02:02Z","service":"hello-service","service":"same-tag","message":"hello"}` + "\n"
	if !equals(output.content, expectedLog) {
		t.Errorf(`expected "%s" but got "%s"`, expectedLog, output.content)
	}
}

func TestPrintAllTaggedMsg(t *testing.T) {
	output := &logout{}
	Clock = clockAt("2020-02-02T02:02:02Z")
	Pretty = false
	Output = output
	l := Logger{
		EnableDebug: false,
	}

	l.Print().
		Str("service", "hello-service").
		Int("items", 12).
		Float("pi", 3.14).
		Err(errors.New("something happened")).
		Msg("hello")

	expectedLog := `{"time":"2020-02-02T02:02:02Z","service":"hello-service","items":12,"pi":3.14,"error":"something happened","message":"hello"}` + "\n"
	if !equals(output.content, expectedLog) {
		t.Errorf(`expected "%s" but got "%s"`, expectedLog, output.content)
	}
}

func TestPrettyPrintMsg(t *testing.T) {
	output := &logout{}
	Clock = clockAt("2020-02-02T02:02:02Z")
	Pretty = true
	Output = output
	l := Logger{
		EnableDebug: false,
	}

	color.NoColor = true

	l.Print().
		Str("service", "hello-service").
		Err(errors.New("couldn't say hello")).
		Msg("hello")

	expectedLog := "2020-02-02T02:02:02Z hello service=hello-service error=couldn't say hello\n"
	if !equals(output.content, expectedLog) {
		t.Errorf(`expected "%s" but got "%s"`, expectedLog, output.content)
	}
	
	color.NoColor = false
}

func TestPrettyPrintEmptyMessage(t *testing.T) {
	output := &logout{}
	Clock = clockAt("2020-02-02T02:02:02Z")
	Pretty = true
	Output = output
	l := Logger{
		EnableDebug: false,
	}

	color.NoColor = true

	l.Print().
		Str("service", "hello-service").
		Err(errors.New("couldn't say hello")).
		Msg("")

	expectedLog := "2020-02-02T02:02:02Z _ service=hello-service error=couldn't say hello\n"
	if !equals(output.content, expectedLog) {
		t.Errorf(`expected "%s" but got "%s"`, expectedLog, output.content)
	}

	color.NoColor = false
}

func TestPrettyPrintNoTags(t *testing.T) {
	output := &logout{}
	Clock = clockAt("2020-02-02T02:02:02Z")
	Pretty = true
	Output = output
	l := Logger{
		EnableDebug: false,
	}

	color.NoColor = true

	l.Print().
		Msg("hello")

	expectedLog := "2020-02-02T02:02:02Z hello\n"
	if !equals(output.content, expectedLog) {
		t.Errorf(`expected "%s" but got "%s"`, expectedLog, output.content)
	}
	
	color.NoColor = false
}

func clockAt(timeStr string) clockwork.Clock {
	t, err := time.Parse(time.RFC3339, timeStr)
	if err != nil {
		panic(err)
	}

	return clockwork.NewFakeClockAt(t)
}

type logout struct {
	content []byte
}

func (l *logout) Write(p []byte) (int, error) {
	l.content = append(l.content, p...)
	return len(p), nil
}

func equals(a []byte, b string) bool {
	if len(a) != len(b) {
		return false
	}

	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}

	return true
}

func TestEquals(t *testing.T) {
	a := []byte("hello")
	b := "hello"

	if !equals(a, b) {
		t.Error("equals not equal")
	}
}

func TestLogout(t *testing.T) {
	output := &logout{}
	fmt.Fprintf(output, "hello")

	if !equals(output.content, "hello") {
		t.Errorf(`expected "hello" but got "%s"`, output.content)
	}
}
