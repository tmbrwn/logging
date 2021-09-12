package logging

import (
	"errors"
	"fmt"
	"regexp"
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

	l.Print("hello")

	expectedLog := `{"time":"2020-02-02T02:02:02Z","message":"hello"}` + "\n"
	if !equals(output.content, expectedLog) {
		t.Errorf(`expected "%s" but got "%s"`, expectedLog, output.content)
	}
}

func TestPrintMsgf(t *testing.T) {
	output := &logout{}
	Clock = clockAt("2020-02-02T02:02:02Z")
	Pretty = false
	Output = output
	l := Logger{
		EnableDebug: false,
	}

	l.Printf("now i have %d fingers", 2)

	expectedLog := `{"time":"2020-02-02T02:02:02Z","message":"now i have 2 fingers"}` + "\n"
	if !equals(output.content, expectedLog) {
		t.Errorf(`expected "%s" but got "%s"`, expectedLog, output.content)
	}
}

func TestGlobalPrintMsg(t *testing.T) {
	output := &logout{}
	Clock = clockAt("2020-02-02T02:02:02Z")
	Pretty = false
	Output = output
	DefaultLogger = Logger{}

	Print("hello")

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

	l.Debug("hello")

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

	l.Debug("hello")

	expectedLog := `{"time":"2020-02-02T02:02:02Z","message":"hello","caller":".*"}` + "\n"
	if !equals(output.content, expectedLog) {
		t.Errorf(`expected "%s" but got "%s"`, expectedLog, output.content)
	}
}

func TestGlobalDebugMsg(t *testing.T) {
	output := &logout{}
	Clock = clockAt("2020-02-02T02:02:02Z")
	Pretty = false
	Output = output
	DefaultLogger = Logger{
		EnableDebug: true,
	}

	Debug("hello")

	expectedLog := `{"time":"2020-02-02T02:02:02Z","message":"hello","caller":".*"}` + "\n"
	if !equals(output.content, expectedLog) {
		t.Errorf(`expected "%s" but got "%s"`, expectedLog, output.content)
	}
}

func TestGlobalDebugMsgNoDebug(t *testing.T) {
	output := &logout{}
	Clock = clockAt("2020-02-02T02:02:02Z")
	Pretty = false
	Output = output
	DefaultLogger = Logger{}

	Debug("hello")

	if !equals(output.content, "") {
		t.Errorf(`expected empty output but got "%s"`, output.content)
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

	l.Tag("service", "hello-service").Print("hello")

	expectedLog := `{"time":"2020-02-02T02:02:02Z","message":"hello","service":"hello-service"}` + "\n"
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

	l.
		Tag("service", "hello-service").
		Tag("service", "same-tag").
		Print("hello")

	expectedLog := `{"time":"2020-02-02T02:02:02Z","message":"hello","service":"hello-service","service":"same-tag"}` + "\n"
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

	l.
		Tag("service", "hello-service").
		Tag("items", 12).
		Tag("pi", 3.14).
		Err(errors.New("something happened")).
		Print("hello")

	expectedLog := `{"time":"2020-02-02T02:02:02Z","message":"hello","service":"hello-service","items":12,"pi":3.14,"error":"something happened"}` + "\n"
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

	l.
		Tag("service", "hello-service").
		Err(errors.New("couldn't say hello")).
		Print("hello")

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

	l.
		Tag("service", "hello-service").
		Err(errors.New("couldn't say hello")).
		Print("")

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

	l.Print("hello")

	expectedLog := "2020-02-02T02:02:02Z hello\n"
	if !equals(output.content, expectedLog) {
		t.Errorf(`expected "%s" but got "%s"`, expectedLog, output.content)
	}

	color.NoColor = false
}

func TestPrettyPrintDebug(t *testing.T) {
	output := &logout{}
	Clock = clockAt("2020-02-02T02:02:02Z")
	Pretty = true
	Output = output
	l := Logger{
		EnableDebug: true,
	}

	color.NoColor = true

	l.Print("hello")

	expectedLog := "2020-02-02T02:02:02Z hello (.*)\n"
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
	regex := fmt.Sprintf("^%s$", b)
	matched, err := regexp.Match(regex, a)
	return matched && err == nil
}

func TestEquals(t *testing.T) {
	t.Run("no regex", func(t *testing.T) {
		a := []byte("hello")
		b := "hello"

		if !equals(a, b) {
			t.Error("equals not equal")
		}
	})

	t.Run("regex", func(t *testing.T) {
		a := []byte("hi, Arthur")
		b := "hi, .*"

		if !equals(a, b) {
			t.Error("equals not equal over regex")
		}
	})
}

func TestLogout(t *testing.T) {
	output := &logout{}
	fmt.Fprintf(output, "hello")

	if !equals(output.content, "hello") {
		t.Errorf(`expected "hello" but got "%s"`, output.content)
	}
}
