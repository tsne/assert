package assert

import (
	"errors"
	"fmt"
	"testing"
)

type testFataler struct {
	message string
}

func (f *testFataler) fataled() bool {
	return f.message != ""
}

func (f *testFataler) reset() {
	f.message = ""
}

func (f *testFataler) Fatal(args ...interface{}) {
	f.message = fmt.Sprint(args...)
}

func TestTrue(t *testing.T) {
	f := &testFataler{}

	True(f, true)
	if f.fataled() {
		t.Fatalf("[True] unexpected error: %s", f.message)
	}

	f.reset()

	True(f, false)
	if !f.fataled() {
		t.Fatal("[True] expected error, got none")
	}
}

func TestFalse(t *testing.T) {
	f := &testFataler{}

	False(f, false)
	if f.fataled() {
		t.Fatalf("[False] unexpected error: %s", f.message)
	}

	f.reset()

	False(f, true)
	if !f.fataled() {
		t.Fatal("[False] expected error, got none")
	}
}

func TestNil(t *testing.T) {
	f := &testFataler{}

	var v interface{}
	Nil(f, v)
	if f.fataled() {
		t.Fatalf("[Nil] unexpected error: %s", f.message)
	}

	f.reset()

	Nil(f, f)
	if !f.fataled() {
		t.Fatal("[Nil] expected error, got none")
	}

	f.reset()

	var err error
	Nil(f, err)
	if f.fataled() {
		t.Fatalf("[Nil] unexpected error: %s", f.message)
	}

	f.reset()

	var s *struct{}
	Nil(f, s)
	if f.fataled() {
		t.Fatalf("[Nil] unexpected error: %s", f.message)
	}

	f.reset()

	s = &struct{}{}
	Nil(f, s)
	if !f.fataled() {
		t.Fatal("[Nil] expected error, got none")
	}

	f.reset()

	var num int
	Nil(f, num)
	if !f.fataled() {
		t.Fatal("[Nil] expected error, got none")
	}
}

func TestNotNil(t *testing.T) {
	f := &testFataler{}

	var v interface{}
	NotNil(f, v)
	if !f.fataled() {
		t.Fatalf("[NotNil] expected error, got none")
	}

	f.reset()

	NotNil(f, f)
	if f.fataled() {
		t.Fatalf("[NotNil] unexpected error: %s", f.message)
	}

	f.reset()

	var err error
	NotNil(f, err)
	if !f.fataled() {
		t.Fatalf("[Nil] unexpected error: %s", f.message)
	}

	f.reset()

	var s *struct{}
	NotNil(f, s)
	if !f.fataled() {
		t.Fatal("[NotNil] expected error, got none")
	}

	f.reset()

	s = &struct{}{}
	NotNil(f, s)
	if f.fataled() {
		t.Fatalf("[NotNil] unexpected error: %s", f.message)
	}

	f.reset()

	var num int
	NotNil(f, num)
	if f.fataled() {
		t.Fatalf("[NotNil] unexpected error: %s", f.message)
	}
}

func TestErr(t *testing.T) {
	f := &testFataler{}

	var nonerr error
	err1 := errors.New("error one")
	err2 := errors.New("error two")

	Err(f, nonerr, nil)
	if f.fataled() {
		t.Fatalf("[Error] unexpected error: %s", f.message)
	}

	f.reset()

	Err(f, err1, err1)
	if f.fataled() {
		t.Fatalf("[Error] unexpected error: %s", f.message)
	}

	f.reset()

	Err(f, nonerr, err1)
	if !f.fataled() {
		t.Fatal("[Error] expected error, got none")
	}

	f.reset()

	Err(f, err1, err2)
	if !f.fataled() {
		t.Fatal("[Error] expected error, got none")
	}
}

func TestErrMsg(t *testing.T) {
	f := &testFataler{}

	err := errors.New("error message one")

	ErrMsg(f, err, err.Error())
	if f.fataled() {
		t.Fatalf("[ErrorMsg] unexpected error: %s", f.message)
	}

	f.reset()

	ErrMsg(f, err, "wrong error message")
	if !f.fataled() {
		t.Fatal("[ErrorMsg] expected error, got none")
	}

	f.reset()

	ErrMsg(f, nil, err.Error())
	if !f.fataled() {
		t.Fatal("[ErrorMsg] expected error, got none")
	}
}

func TestMsgMatch(t *testing.T) {
	f := &testFataler{}

	err := errors.New("an error message with a 74 in it")
	ErrMsgMatch(f, err, `an error message with a \d{2} in it`)
	if f.fataled() {
		t.Fatalf("[ErrorMsgMatch] unexpected error: %s", f.message)
	}

	f.reset()

	ErrMsgMatch(f, err, `an error message with a \d{3} in it`)
	if !f.fataled() {
		t.Fatal("[ErrorMsg] expected error, got none")
	}

	f.reset()

	ErrMsgMatch(f, nil, `an error message with a \d{3} in it`)
	if !f.fataled() {
		t.Fatal("[ErrorMsg] expected error, got none")
	}
}

func TestEqual(t *testing.T) {
	f := &testFataler{}

	i := 7
	Equal(f, i, 7)
	if f.fataled() {
		t.Fatalf("[Equal] unexpected error: %s", f.message)
	}

	f.reset()

	Equal(f, i, 13)
	if !f.fataled() {
		t.Fatal("[Equal] expected error, got none")
	}

	f.reset()

	d := 1.5
	Equal(f, d, 1.5)
	if f.fataled() {
		t.Fatalf("[Equal] unexpected error: %s", f.message)
	}

	f.reset()

	Equal(f, d, i)
	if !f.fataled() {
		t.Fatal("[Equal] expected error, got none")
	}

	f.reset()

	p := &struct{}{}
	Equal(f, p, nil)
	if !f.fataled() {
		t.Fatal("[Equal] expected error, got none")
	}

	f.reset()

	Equal(f, p, p)
	if f.fataled() {
		t.Fatalf("[Equal] unexpected error: %s", f.message)
	}

	f.reset()

	s := &struct {
		field1 int
		field2 string
	}{
		7,
		"seven",
	}
	Equal(f, s, p)
	if !f.fataled() {
		t.Fatal("[Equal] expected error, got none")
	}

	f.reset()

	Equal(f, s, &struct {
		field1 int
		field2 string
	}{
		7,
		"seven",
	})
	if f.fataled() {
		t.Fatalf("[Equal] unexpected error: %s", f.message)
	}

	f.reset()

	Equal(f, s, &struct {
		field1 int
		field2 string
	}{
		7,
		"Seven",
	})
	if !f.fataled() {
		t.Fatal("[Equal] expected error, got none")
	}
}

func TestNotEqual(t *testing.T) {
	f := &testFataler{}

	i := 7
	NotEqual(f, i, 13)
	if f.fataled() {
		t.Fatalf("[NotEqual] unexpected error: %s", f.message)
	}

	f.reset()

	NotEqual(f, i, 7)
	if !f.fataled() {
		t.Fatal("[NotEqual] expected error, got none")
	}

	f.reset()

	d := 1.5
	NotEqual(f, d, i)
	if f.fataled() {
		t.Fatalf("[NotEqual] unexpected error: %s", f.message)
	}

	f.reset()

	NotEqual(f, d, 1.5)
	if !f.fataled() {
		t.Fatal("[NotEqual] expected error, got none")
	}

	f.reset()

	p := &struct{}{}
	NotEqual(f, p, p)
	if !f.fataled() {
		t.Fatal("[NotEqual] expected error, got none")
	}

	f.reset()

	NotEqual(f, p, nil)
	if f.fataled() {
		t.Fatalf("[NotEqual] unexpected error: %s", f.message)
	}

	f.reset()

	s := &struct {
		field1 int
		field2 string
	}{
		7,
		"seven",
	}
	NotEqual(f, s, p)
	if f.fataled() {
		t.Fatalf("[NotEqual] unexpected error: %s", f.message)
	}

	f.reset()

	NotEqual(f, s, &struct {
		field1 int
		field2 string
	}{
		7,
		"Seven",
	})
	if f.fataled() {
		t.Fatalf("[NotEqual] unexpected error: %s", f.message)
	}

	f.reset()

	NotEqual(f, s, &struct {
		field1 int
		field2 string
	}{
		7,
		"seven",
	})
	if !f.fataled() {
		t.Fatal("[NotEqual] expected error, got none")
	}
}
