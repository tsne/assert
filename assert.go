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

// Fataler represents an interface to trigger Fatal if an assertion does
// not hold.
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

// True asserts that the given condition is fulfilled.
func True(f Fataler, condition bool) {
	if !condition {
		fatal(f, 1, "expected true, got false")
	}
}

// False asserts that the given condition is not fulfilled.
func False(f Fataler, condition bool) {
	if condition {
		fatal(f, 1, "expected false, got true")
	}
}

// Nil asserts that the given value is nil.
func Nil(f Fataler, value interface{}) {
	if !isNil(value) {
		fatal(f, 1, "expected nil, got %+v", value)
	}
}

// NotNil asserts that the given value is not nil.
func NotNil(f Fataler, value interface{}) {
	if isNil(value) {
		fatal(f, 1, "expected a value, got nil")
	}
}

// Err asserts that the actual error equals the expected error.
func Err(f Fataler, actual, expected error) {
	if actual != expected {
		fatal(f, 1, "unexpected error\nexpected: %v\ngot: %v", expected, actual)
	}
}

// ErrMsg asserts that the message of the given error equals the given message.
func ErrMsg(f Fataler, err error, msg string) {
	if err == nil || err.Error() != msg {
		actual := ""
		if err != nil {
			actual = err.Error()
		}
		fatal(f, 1, "unexpected error message\nexpected: %s\ngot: %s", msg, actual)
	}
}

// ErrMsgMatch asserts that the message of the given error matches the given
// regular expression.
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

// Equal asserts that the actual value equals the expected value. This function
// also compares the elements of arrays, slices, maps, and struct fields.
func Equal(f Fataler, actual, expected interface{}) {
	if !reflect.DeepEqual(actual, expected) {
		fatal(f, 1, "unexpected value\n expected: %+v\ngot: %+v", expected, actual)
	}
}

// NotEqual asserts that the actual value is not equal to the expected value.
// This functions also compares the elements of arrays, slices, maps, and struct fields.
func NotEqual(f Fataler, actual, expected interface{}) {
	if reflect.DeepEqual(actual, expected) {
		fatal(f, 1, "unexpected different values, got the same")
	}
}

// Panic asserts that the function fn called with the arguments args panics.
func Panic(f Fataler, fn interface{}, args ...interface{}) {
	if panicked, _ := recoverPanic(f, fn, args...); !panicked {
		fatal(f, 1, "expected a panic, got none")
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

func recoverPanic(f Fataler, fn interface{}, args ...interface{}) (bool, string) {
	function := reflect.ValueOf(fn)
	if function.Kind() != reflect.Func {
		fatal(f, 2, "expected function, got %s", function.Kind())
	}

	arguments := make([]reflect.Value, 0, len(args))
	for _, a := range args {
		arguments = append(arguments, reflect.ValueOf(a))
	}

	var (
		panicked bool
		message  string
	)
	func() {
		defer func() {
			if r := recover(); r != nil {
				panicked = true
				message = fmt.Sprintf("%v", r)
			}
		}()
		function.Call(arguments)
	}()
	return panicked, message
}
