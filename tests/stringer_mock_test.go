/*
DO NOT EDIT!
This code was generated automatically using github.com/gojuno/minimock v1.2
Original interface "Stringer" can be found in github.com/gojuno/minimock/tests
*/
package tests

import (
	"sync/atomic"
	"time"
)

//StringerMock implements github.com/gojuno/minimock/tests.Stringer
type StringerMock struct {
	t *TesterMock

	StringFunc func() (r string)

	StringCounter uint64

	StringMock mStringerMockString
}

//NewStringerMock returns a mock for github.com/gojuno/minimock/tests.Stringer
func NewStringerMock(t *TesterMock) *StringerMock {
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
