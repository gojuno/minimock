## Summary [![GoDoc](https://godoc.org/github.com/gojuno/minimock?status.svg)](http://godoc.org/github.com/gojuno/minimock) [![Build Status](https://travis-ci.org/gojuno/minimock.svg?branch=master)](https://travis-ci.org/gojuno/minimock) [![Go Report Card](https://goreportcard.com/badge/github.com/gojuno/minimock)](https://goreportcard.com/report/github.com/gojuno/minimock)
Minimock parses input Go source file containing interface declaration and generates
implementation of this interface that can be used as a mock.

The main feature of Minimock is that you can reuse generated implementation in different
test cases by attaching interface method implementations on the fly.

Let's say we have following interface declaration:

```go
package fmt

type Stringer interface {
  String() string
}
``` 

For such interface generated implementation will look like:
```go
type StringerMock struct {
	t *testing.T
	m *sync.RWMutex

	StringFunc func() (r0 string)

	StringCounter int

	StringMock StringerMockString
}

func NewStringerMock(t *testing.T) *StringerMock {
	m := &StringerMock{t: t, m: &sync.RWMutex{}}
	m.StringMock = StringerMockString{mock: m}

	return m
}

type StringerMockString struct {
	mock *StringerMock
}

func (m StringerMockString) Return(r0 string) *StringerMock {
	m.mock.StringFunc = func() string {
		return r0
	}
	return m.mock
}

func (m StringerMockString) Set(f func()) *StringerMock {
	m.mock.StringFunc = f
	return m.mock
}

func (m *StringerMock) String() (r0 string) {
	m.m.Lock()
	m.StringCounter += 1
	m.m.Unlock()

	if m.StringFunc == nil {
		m.t.Errorf("Unexpected call to StringerMock.String")
	}

	return m.StringFunc()
}

func (m *StringerMock) ValidateCallCounters() {
	m.t.Log("ValidateCallCounters is deprecated please use CheckMocksCalled")

	if m.StringFunc != nil && m.StringCounter == 0 {
		m.t.Error("Expected call to StringerMock.String")
	}

}

func (m *StringerMock) CheckMocksCalled() {

	if m.StringFunc != nil && m.StringCounter == 0 {
		m.t.Error("Expected call to StringerMock.String")
	}

}
```

In the test you can use Return helper or you can define StringerFunc for more complex behaviour:
```go

func TestStringerUser(t *testing.T) {
  stringerMock := NewStringerMock(t)
  defer stringerMock.CheckMocksCalled()

  stringerMock.StringMock.Return("Hello, world!")

  //... code that uses stringerMock
}

func TestStringerUserComplex(t *testing.T) {
  stringerMock := NewStringerMock(t)
  defer stringerMock.CheckMocksCalled()

  stringerMock.StringFunc = func() string {
    switch stringerMock.StringCounter {
    case 1:
      return "Hello,"
    case 2:
      return "world"
    default:
      return "!"
    }
  }

  //... code that uses stringerMock
}
```

Alternatively, you can use builder-style mock configuration:
```go
stringerMock := NewStringerMock(t).
  StringMock.Return("Hello, world!").
  IntMock.Return(1).
  MultiplierMock.Set(func(a,b int) int { //example of the mock that checks input params
    assert.Equal(t, 2, a)
    assert.Equal(t, 3, b)
    return a * b
  })

defer stringerMock.CheckMocksCalled()
```

If caller performs a call to method that is not mocked the test case will fail.
If caller does not perform a call to method that is mocked the test case will fail if you call to mock.ValidateCallCounters().
You can also perform more precise checks by using concrete call counters, i.e. StringerMock.StringCounter 

Please see more detailed example in examples subpackage.

## Usage of minimock:
```
  -f string
    input file or import path of the package containing interface declaration
  -i string
    interface name
  -o string
    destination file for interface implementation
  -p string
    destination package name
  -t string
    target struct name, default: <interface name>Mock
```

## Usage of minimock in go:generate instruction:
```go
//go:generate minimock -f fmt -i Stringer -o ./stringer_mock_test.go -p examples
```
