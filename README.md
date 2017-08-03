![logo](https://rawgit.com/gojuno/minimock/master/logo.svg)
[![GoDoc](https://godoc.org/github.com/gojuno/minimock?status.svg)](http://godoc.org/github.com/gojuno/minimock) [![Build Status](https://travis-ci.org/gojuno/minimock.svg?branch=master)](https://travis-ci.org/gojuno/minimock) [![Go Report Card](https://goreportcard.com/badge/github.com/gojuno/minimock)](https://goreportcard.com/report/github.com/gojuno/minimock) ![cover.run go](https://cover.run/go/github.com/gojuno/minimock/tests.svg)

## Summary 
Minimock parses the input Go source file that contains an interface declaration and generates
implementation of this interface that can be used as a mock.

Main features of minimock:

* It's integrated with the standard Go "testing" package
* It's very convenient to use generated mocks in table tests because it implements builder pattern to set up several mocks
* It provides a useful MinimockWait(time.Duration) helper to test concurrent code
* It generates helpers to check if the mocked methods have been called and keeps your tests clean and up to date
* It generates concurrent-safe mock execution counters that you can use in your mocks to implement sophisticated mocks behaviour

## Usage
Let's say we have the following interface declaration in github.com/gojuno/minimock/tests package:
```go
type Stringer interface {
  fmt.Stringer
}
```

Here is how to generate the mock for this interface:
```
minimock.go -f github.com/gojuno/minimock/tests -i Stringer -o ./tests/stringer_mock_test.go -p tests
```

The result file ./tests/stringer_mock_test.go will be:
```go
//StringerMock implements github.com/gojuno/minimock/tests.Stringer
type StringerMock struct {
  t minimock.Tester

  StringFunc func() (r string)
  StringCounter uint64
  StringMock mStringerMockString
}

//NewStringerMock returns a mock for github.com/gojuno/minimock/tests.Stringer
func NewStringerMock(t minimock.Tester) *StringerMock {
  m := &StringerMock{t: t}

  if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

  m.StringMock = mStringerMockString{mock: m}

  return m
}

type mStringerMockString struct {
  mock *StringerMock
}

//Return sets up a mock for Stringer.String to return Return's arguments
func (m mStringerMockString) Return(r string) *StringerMock {
  m.mock.StringFunc = func() string {
    return r
  }
  return m.mock
}

//Set uses a given function f as a mock of Stringer.String string method
func (m mStringerMockString) Set(f func() (r string)) *StringerMock {
  m.mock.StringFunc = f
  return m.mock
}

//String implements github.com/gojuno/minimock/tests.Stringer interface
func (m *StringerMock) String() (r string) {
  defer atomic.AddUint64(&m.StringCounter, 1)

  if m.StringFunc == nil {
    m.t.Fatal("Unexpected call to StringerMock.String")
    return
  }

  return m.StringFunc()
}

//MinimockFinish checks that all mocked functions of an iterface have been called at least once
func (m *StringerMock) MinimockFinish() {
  if m.StringFunc != nil && m.StringCounter == 0 {
    m.t.Fatal("Expected call to StringerMock.String")
  }
}

//MinimockWait waits for all mocked functions to be called at least once
func (m *StringerMock) MinimockWait(timeout time.Duration) {
  timeoutCh := time.After(timeout)
  for {
    ok := true
    ok = ok && (m.StringFunc == nil || m.StringCounter > 0)

    if ok {
      return
    }

    select {
    case <-timeoutCh:

      if m.StringFunc != nil && m.StringCounter == 0 {
        m.t.Error("Expected call to StringerMock.String")
      }

      m.t.Fatalf("Some mocks were not called on time: %s", timeout)
      return
    default:
      time.Sleep(time.Millisecond)
    }
  }
}

//AllMocksCalled returns true if all mocked methods were called before the execution of AllMocksCalled,
//it can be used with assert/require, i.e. assert.True(mock.AllMocksCalled())
func (m *StringerMock) AllMocksCalled() bool {

  if m.StringFunc != nil && m.StringCounter == 0 {
    return false
  }

  return true
}
```


There are several ways to set up a mock

Setting up a mock using direct assignment:
```go
stringerMock := NewStringerMock(mc)
stringerMock.StringFunc = func() string {
  return "minimock"
}
```

Setting up a mock using builder pattern and Return method:
```go
stringerMock := NewStringerMock(mc).StringMock.Return("minimock")
```

Setting up a mock using builder and Set method:
```go
stringerMock := NewStringerMock(mc).StringMock.Set(func() string {
  return "minimock"
})
```

Builder pattern is convenient when you have to mock more than one method of an interface.
Imagine we have StringerInter interface with two methods:
```go
type StringerInter interface {
  String() string
  Int() int
}
```

Then you can set up a mock using just one assignment:
```go
stringerMock := NewStringerMock(mc).StringMock.Return("minimock").IntMock.Return(5)
```

You can also use invocation counters in your mocks and tests:
```go
stringerMock := NewStringerMock(mc)
stringerMock.StringFunc = func() string {
  return fmt.Sprintf("minimock: %d", stringerMock.StrigCounter)
}
```

## Keep your tests clean
Sometimes we write tons of mocks for our tests but over time the tested code stops using mocked dependencies,
however mocks are still present and being initialized in the test files. So while tested code can shrink, tests are only growing.
To prevent this minimock provides MinimockFinish() method that verifies that all your mocks have been called at least once during the test run.

```go
func TestSomething(t *testing.T) {
  mc := minimock.NewController(t)
  //this will mark your test as failed because there's no stringerMock.String() invocation below
  defer mc.Finish()

  stringerMock := NewStringerMock(mc)
  stringerMock.StringMock.Return("minimock")
}
```

## Testing concurrent code
Testing concurrent code is tough. Fortunately minimock provides you with the helper method that makes testing concurrent code easy.
Here is how it works:

```go
func TestSomething(t *testing.T) {
  mc := minimock.NewController(t)

  //Wait ensures that all mocked methods have been called within given interval
  //if any of the mocked methods have not been called Wait marks test as failed
  defer mc.Wait(time.Second)

  stringerMock := NewStringerMock(mc)
  stringerMock.StringMock.Return("minimock")

  //tested code can run mocked method in a goroutine
  go stirngerMock.String()
}
```

## Minimock command line flags:
```
$ minimock 
Usage of minimock:
  -f string
    	input file or import path of the package that contains interface declaration
  -i string
    	name of the interface to mock
  -o string
    	destination file name to place the generated mock
  -p string
    	destination package name
  -t string
    	mock struct name, default is: <interface name>Mock
  -withTests
```
