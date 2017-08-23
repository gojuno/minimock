/*
DO NOT EDIT!
This code was generated automatically using github.com/gojuno/minimock v1.5
The original interface "Stringer" can be found in github.com/gojuno/minimock/tests
*/
package tests

import (
	"sync/atomic"
	"time"

	"github.com/gojuno/minimock"
)

//StringerMock implements github.com/gojuno/minimock/tests.Stringer
type StringerMock struct {
	t minimock.Tester

	StringFunc    func() (r string)
	StringCounter uint64
	StringMock    mStringerMockString
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

//Set uses given function f as a mock of Stringer.String method
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

//ValidateCallCounters checks that all mocked methods of the iterface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *StringerMock) ValidateCallCounters() {

	if m.StringFunc != nil && atomic.LoadUint64(&m.StringCounter) == 0 {
		m.t.Fatal("Expected call to StringerMock.String")
	}

}

//CheckMocksCalled checks that all mocked methods of the iterface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *StringerMock) CheckMocksCalled() {
	m.Finish()
}

//Finish checks that all mocked methods of the iterface have been called at least once
//Deprecated: please use MinimockFinish or use Finish method of minimock.Controller
func (m *StringerMock) Finish() {
	m.MinimockFinish()
}

//MinimockFinish checks that all mocked methods of the iterface have been called at least once
func (m *StringerMock) MinimockFinish() {

	if m.StringFunc != nil && atomic.LoadUint64(&m.StringCounter) == 0 {
		m.t.Fatal("Expected call to StringerMock.String")
	}

}

//Wait waits for all mocked methods to be called at least once
//Deprecated: please use MinimockWait or use Wait method of minimock.Controller
func (m *StringerMock) Wait(timeout time.Duration) {
	m.MinimockWait(timeout)
}

//MinimockWait waits for all mocked methods to be called at least once
//this method is called by minimock.Controller
func (m *StringerMock) MinimockWait(timeout time.Duration) {
	timeoutCh := time.After(timeout)
	for {
		ok := true
		ok = ok && (m.StringFunc == nil || atomic.LoadUint64(&m.StringCounter) > 0)

		if ok {
			return
		}

		select {
		case <-timeoutCh:

			if m.StringFunc != nil && atomic.LoadUint64(&m.StringCounter) == 0 {
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

	if m.StringFunc != nil && atomic.LoadUint64(&m.StringCounter) == 0 {
		return false
	}

	return true
}
