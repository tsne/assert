package assert

import (
	"bytes"
	"fmt"
	"path/filepath"
	"reflect"
	"regexp"
	"runtime"
	"strings"
	"testing"
)

type Fataler interface {
	Fatal(args ...interface{})
}

// If we are in a testing scenario we will need the line number
// of the Fatal call to remove the "file:line" prefix. To make sure
// this line number has two digits only this function is at the top
// of the file.
func fatal(f Fataler, skip int, message string, args ...interface{}) {
	buf := bytes.NewBuffer(nil)
	if _, isT := f.(*testing.T); isT {
		// We have to remove the "<filename>:<linenumber>: " prefix.
		// So we take the length of the file plus two line number digits
		// plus two colons plus a space separator.
		_, thisFile, _, _ := runtime.Caller(0)
		thisFile = filepath.Base(thisFile)
		fmt.Fprint(buf, strings.Repeat("\b", len(thisFile)+5))
	}

	_, file, line, _ := runtime.Caller(skip + 1)
	fmt.Fprintf(buf, "%s:%d: ", filepath.Base(file), line)
	fmt.Fprintf(buf, message, args...)
	f.Fatal(buf.String())
}

func True(f Fataler, condition bool) {
	if !condition {
		fatal(f, 1, "expected true, got false")
	}
}

func False(f Fataler, condition bool) {
	if condition {
		fatal(f, 1, "expected false, got true")
	}
}

func Nil(f Fataler, value interface{}) {
	if !isNil(value) {
		fatal(f, 1, "expected nil, got %+v", value)
	}
}

func NotNil(f Fataler, value interface{}) {
	if isNil(value) {
		fatal(f, 1, "expected a value, got nil")
	}
}

func Err(f Fataler, actual, expected error) {
	if actual != expected {
		fatal(f, 1, "unexpected error\nexpected: %v\ngot: %v", expected, actual)
	}
}

func ErrMsg(f Fataler, err error, msg string) {
	if err == nil || err.Error() != msg {
		actual := ""
		if err != nil {
			actual = err.Error()
		}
		fatal(f, 1, "unexpected error message\nexpected: %s\ngot: %s", msg, actual)
	}
}

func ErrMsgMatch(f Fataler, err error, msgRegexp string) {
	matches := false
	actual := ""
	if err != nil {
		actual = err.Error()
		rx := regexp.MustCompile(msgRegexp)
		matches = rx.MatchString(actual)
	}

	if !matches {
		fatal(f, 1, "unexpected error pattern\nexpected: %s\ngot: %s", msgRegexp, actual)
	}
}

func Equal(f Fataler, actual, expected interface{}) {
	if !reflect.DeepEqual(actual, expected) {
		fatal(f, 1, "unexpected value\n expected: %+v\ngot: %+v", expected, actual)
	}
}

func NotEqual(f Fataler, actual, expected interface{}) {
	if reflect.DeepEqual(actual, expected) {
		fatal(f, 1, "unexpected different values, got the same")
	}
}

func isNil(v interface{}) bool {
	if v == nil {
		return true
	}

	val := reflect.ValueOf(v)
	switch val.Kind() {
	case reflect.Invalid:
		return true
	case reflect.Chan, reflect.Func, reflect.Interface, reflect.Ptr, reflect.Slice, reflect.Map:
		return val.IsNil()
	default:
		return false
	}
}
