###Summary
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
}

func NewStringerMock(t *testing.T) *StringerMock {
	return &StringerMock{t: t, m: &sync.RWMutex{}}
}

func (m *StringerMock) String() (r0 string) {
	m.m.Lock()
	m.StringCounter += 1
	m.m.Unlock()

	if m.StringFunc == nil {
		m.t.Fatalf("Unexpected call to StringerMock.String")
	}

	return m.StringFunc()
}

func (m *StringerMock) CheckMocksCalled() {

	if m.StringFunc != nil && m.StringCounter == 0 {
		m.t.Error("Expected call to StringerMock.String")
	}

}
```

If caller performs a call to method that is not mocked the test case will fail.
If caller does not perform a call to method that is mocked the test case will fail if you call to mock.ValidateCallCounters().
You can also perform more precise checks by using concrete call counters, i.e. StringerMock.StringCounter 

Please see more detailed example in examples subpackage.

###Usage of minimock:
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

###Usage of minimock in go:generate instruction:
```go
//go:generate minimock -f fmt -i Stringer -o ./stringer_mock_test.go -p examples
```
