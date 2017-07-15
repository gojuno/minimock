/*
DO NOT EDIT!
This code was generated automatically using github.com/gojuno/minimock v1.2
Original interface "Tester" can be found in github.com/gojuno/minimock/tests
*/
package tests

import (
	"sync/atomic"
	"testing"
	"time"
)

//TesterMock implements github.com/gojuno/minimock/tests.Tester
type TesterMock struct {
	t *testing.T

	ErrorFunc  func(p ...interface{})
	FatalFunc  func(p ...interface{})
	FatalfFunc func(p string, p1 ...interface{})

	ErrorCounter  uint64
	FatalCounter  uint64
	FatalfCounter uint64

	ErrorMock  mTesterMockError
	FatalMock  mTesterMockFatal
	FatalfMock mTesterMockFatalf
}

//NewTesterMock returns a mock for github.com/gojuno/minimock/tests.Tester
func NewTesterMock(t *testing.T) *TesterMock {
	m := &TesterMock{t: t}
	m.ErrorMock = mTesterMockError{mock: m}
	m.FatalMock = mTesterMockFatal{mock: m}
	m.FatalfMock = mTesterMockFatalf{mock: m}

	return m
}

type mTesterMockError struct {
	mock *TesterMock
}

//Return set up a mock for Tester.func(...interface{}) to return Return's arguments
func (m mTesterMockError) Return() *TesterMock {
	m.mock.ErrorFunc = func(p ...interface{}) {
		return
	}
	return m.mock
}

//Set uses given function f as a mock of Tester.func(...interface{}) method
func (m mTesterMockError) Set(f func(p ...interface{})) *TesterMock {
	m.mock.ErrorFunc = f
	return m.mock
}

//Error implements github.com/gojuno/minimock/tests.Tester interface
func (m *TesterMock) Error(p ...interface{}) {
	defer atomic.AddUint64(&m.ErrorCounter, 1)

	if m.ErrorFunc == nil {
		m.t.Fatal("Unexpected call to TesterMock.Error")
		return
	}

	m.ErrorFunc(p...)
}

type mTesterMockFatal struct {
	mock *TesterMock
}

//Return set up a mock for Tester.func(args ...interface{}) to return Return's arguments
func (m mTesterMockFatal) Return() *TesterMock {
	m.mock.FatalFunc = func(p ...interface{}) {
		return
	}
	return m.mock
}

//Set uses given function f as a mock of Tester.func(args ...interface{}) method
func (m mTesterMockFatal) Set(f func(p ...interface{})) *TesterMock {
	m.mock.FatalFunc = f
	return m.mock
}

//Fatal implements github.com/gojuno/minimock/tests.Tester interface
func (m *TesterMock) Fatal(p ...interface{}) {
	defer atomic.AddUint64(&m.FatalCounter, 1)

	if m.FatalFunc == nil {
		m.t.Fatal("Unexpected call to TesterMock.Fatal")
		return
	}

	m.FatalFunc(p...)
}

type mTesterMockFatalf struct {
	mock *TesterMock
}

//Return set up a mock for Tester.func(format string, args ...interface{}) to return Return's arguments
func (m mTesterMockFatalf) Return() *TesterMock {
	m.mock.FatalfFunc = func(p string, p1 ...interface{}) {
		return
	}
	return m.mock
}

//Set uses given function f as a mock of Tester.func(format string, args ...interface{}) method
func (m mTesterMockFatalf) Set(f func(p string, p1 ...interface{})) *TesterMock {
	m.mock.FatalfFunc = f
	return m.mock
}

//Fatalf implements github.com/gojuno/minimock/tests.Tester interface
func (m *TesterMock) Fatalf(p string, p1 ...interface{}) {
	defer atomic.AddUint64(&m.FatalfCounter, 1)

	if m.FatalfFunc == nil {
		m.t.Fatal("Unexpected call to TesterMock.Fatalf")
		return
	}

	m.FatalfFunc(p, p1...)
}

//DEPRECATED: please use CheckMocksCalled
func (m *TesterMock) ValidateCallCounters() {

	if m.ErrorFunc != nil && m.ErrorCounter == 0 {
		m.t.Fatal("Expected call to TesterMock.Error")
	}

	if m.FatalFunc != nil && m.FatalCounter == 0 {
		m.t.Fatal("Expected call to TesterMock.Fatal")
	}

	if m.FatalfFunc != nil && m.FatalfCounter == 0 {
		m.t.Fatal("Expected call to TesterMock.Fatalf")
	}

}

//CheckMocksCalled checks that all mocked functions of an iterface have been called at least once
func (m *TesterMock) CheckMocksCalled() {

	if m.ErrorFunc != nil && m.ErrorCounter == 0 {
		m.t.Fatal("Expected call to TesterMock.Error")
	}

	if m.FatalFunc != nil && m.FatalCounter == 0 {
		m.t.Fatal("Expected call to TesterMock.Fatal")
	}

	if m.FatalfFunc != nil && m.FatalfCounter == 0 {
		m.t.Fatal("Expected call to TesterMock.Fatalf")
	}

}

//Wait waits for all mocked functions to be called at least once
func (m *TesterMock) Wait(timeout time.Duration) {
	timeoutCh := time.After(timeout)
	for {
		ok := true
		ok = ok && (m.ErrorFunc == nil || m.ErrorCounter > 0)
		ok = ok && (m.FatalFunc == nil || m.FatalCounter > 0)
		ok = ok && (m.FatalfFunc == nil || m.FatalfCounter > 0)

		if ok {
			return
		}

		select {
		case <-timeoutCh:

			if m.ErrorFunc != nil && m.ErrorCounter == 0 {
				m.t.Error("Expected call to TesterMock.Error")
			}

			if m.FatalFunc != nil && m.FatalCounter == 0 {
				m.t.Error("Expected call to TesterMock.Fatal")
			}

			if m.FatalfFunc != nil && m.FatalfCounter == 0 {
				m.t.Error("Expected call to TesterMock.Fatalf")
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
func (m *TesterMock) AllMocksCalled() bool {

	if m.ErrorFunc != nil && m.ErrorCounter == 0 {
		return false
	}

	if m.FatalFunc != nil && m.FatalCounter == 0 {
		return false
	}

	if m.FatalfFunc != nil && m.FatalfCounter == 0 {
		return false
	}

	return true
}
