###Summary
Minimock parses input Go source file containing interface declaration and generates
implementation of this interface that can be used as a mock.

The main feature of Minimock is that you can reuse generated implementation in different
test cases by attaching interface method implementations on the fly.

Let's say we have following interface declaration:

```go
package sample

type Interface interface {
  GetString() string
}
``` 

For such interface generated implementation will look like:
```go
type InterfaceMock struct {
  t *testing.T
  m *sync.Mutex

  GetStringFunc         func() string
  GetStringCounter      int
}

func (m *InterfaceMock) GetString() string {
  m.m.Lock()
  m.GetStringCounter += 1
  m.m.Unlock()

  if m.GetStringFunc == nil {
    m.t.Fatalf("Unexpected call to InterfaceMock.GetString")
  }

  return m.GetStringFunc()
}

func (m *InterfaceMock) ValidateCallCounters() {
  if m.GetStringFunc != nil && m.GetStringCounter == 0 {
    m.t.Error("Expected call to InterfaceMock.GetString")
  }
}
```

If caller performs a call to method that is not mocked the test case will fail.
If caller does not perform a call to method that is mocked the test case will fail if you call to mock.ValidateCallCounters().
You can also perform more precise checks by using concrete call counters, i.e. mock.GetStringCounter 

Please see more detailed example in examples subpackage.

###Usage of minimock:
```
  -f string
      input file containing interface declaration
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
//go:generate minimock -f ./sample/interface.go -i Interface -o ./sample_interface_mock.go -p examples
```
