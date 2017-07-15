/*
Minimock is a command line tool that generates mocks for the Go interfaces

Let's say we have following interface declaration in github.com/gojuno/minimock/tests package:

	type Stringer interface {
  	fmt.Stringer
	}

And we want to generate mock for this interface:

	minimock.go -f github.com/gojuno/minimock/tests -i Stringer -o ./tests/stringer_mock_test.go -p tests

Result file ./tests/stringer_mock_test.go will be:

	//StringerMock implements github.com/gojuno/minimock/tests.Stringer
	type StringerMock struct {
		t *testing.T

		StringFunc func() (r string)

		StringCounter uint64

		StringMock mStringerMockString
	}

	//NewStringerMock returns a mock for github.com/gojuno/minimock/tests.Stringer
	func NewStringerMock(t *testing.T) *StringerMock {
		m := &StringerMock{t: t}
		m.StringMock = mStringerMockString{mock: m}

		return m
	}

	type mStringerMockString struct {
		mock *StringerMock
	}

	//Return set up a mock for Stringer.func() string to return Return's arguments
	func (m mStringerMockString) Return(r string) *StringerMock {
		m.mock.StringFunc = func() string {
			return r
		}
		return m.mock
	}

	//Set uses given function f as a mock of Stringer.func() string method
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

	//DEPRECATED: please use CheckMocksCalled
	func (m *StringerMock) ValidateCallCounters() {

		if m.StringFunc != nil && m.StringCounter == 0 {
			m.t.Fatal("Expected call to StringerMock.String")
		}

	}

	//CheckMocksCalled checks that all mocked functions of an iterface have been called at least once
	func (m *StringerMock) CheckMocksCalled() {

		if m.StringFunc != nil && m.StringCounter == 0 {
			m.t.Fatal("Expected call to StringerMock.String")
		}

	}

	//Wait waits for all mocked functions to be called at least once
	func (m *StringerMock) Wait(timeout time.Duration) {
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

	//AllMocksCalled returns true if all mocked methods were called before the call to AllMocksCalled,
	//it can be used with assert/require, i.e. assert.True(mock.AllMocksCalled())
	func (m *StringerMock) AllMocksCalled() bool {

		if m.StringFunc != nil && m.StringCounter == 0 {
			return false
		}

		return true
	}

Now we can use our StringerMock as a Stringer implementation in your tests:

func TestStringer(t *testing.T) {
	mock := NewStringerMock(t)
	mock.StringMock.Return("")

	assert.Equal(t, "", mock.String())

	//more sophisticated example
	mock.StringFunc = func() string {
		return fmt.Sprintf("%d", mock.StringCounter)
	}

	assert.Equal(t, "1", mock.String())
	assert.Equal(t, "2", mock.String())
}
*/
package main
