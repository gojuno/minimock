/*
DO NOT EDIT!
This code was generated automatically using github.com/gojuno/minimock v1.5
The original interface "Tester" can be found in github.com/gojuno/minimock
*/
package tests

import (
	"sync/atomic"
	"time"

	"github.com/gojuno/minimock"
)

//TesterMock implements github.com/gojuno/minimock.Tester
type TesterMock struct {
	t minimock.Tester

	ErrorFunc    func(p ...interface{})
	ErrorCounter uint64
	ErrorMock    mTesterMockError

	FatalFunc    func(p ...interface{})
	FatalCounter uint64
	FatalMock    mTesterMockFatal

	FatalfFunc    func(p string, p1 ...interface{})
	FatalfCounter uint64
	FatalfMock    mTesterMockFatalf
}

//NewTesterMock returns a mock for github.com/gojuno/minimock.Tester
func NewTesterMock(t minimock.Tester) *TesterMock {
	m := &TesterMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.ErrorMock = mTesterMockError{mock: m}
	m.FatalMock = mTesterMockFatal{mock: m}
	m.FatalfMock = mTesterMockFatalf{mock: m}

	return m
}

type mTesterMockError struct {
	mock *TesterMock
}

//Return sets up a mock for Tester.Error to return Return's arguments
func (m mTesterMockError) Return() *TesterMock {
	m.mock.ErrorFunc = func(p ...interface{}) {
		return
	}
	return m.mock
}

//Set uses given function f as a mock of Tester.Error method
func (m mTesterMockError) Set(f func(p ...interface{})) *TesterMock {
	m.mock.ErrorFunc = f
	return m.mock
}

//Error implements github.com/gojuno/minimock.Tester interface
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

//Return sets up a mock for Tester.Fatal to return Return's arguments
func (m mTesterMockFatal) Return() *TesterMock {
	m.mock.FatalFunc = func(p ...interface{}) {
		return
	}
	return m.mock
}

//Set uses given function f as a mock of Tester.Fatal method
func (m mTesterMockFatal) Set(f func(p ...interface{})) *TesterMock {
	m.mock.FatalFunc = f
	return m.mock
}

//Fatal implements github.com/gojuno/minimock.Tester interface
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

//Return sets up a mock for Tester.Fatalf to return Return's arguments
func (m mTesterMockFatalf) Return() *TesterMock {
	m.mock.FatalfFunc = func(p string, p1 ...interface{}) {
		return
	}
	return m.mock
}

//Set uses given function f as a mock of Tester.Fatalf method
func (m mTesterMockFatalf) Set(f func(p string, p1 ...interface{})) *TesterMock {
	m.mock.FatalfFunc = f
	return m.mock
}

//Fatalf implements github.com/gojuno/minimock.Tester interface
func (m *TesterMock) Fatalf(p string, p1 ...interface{}) {
	defer atomic.AddUint64(&m.FatalfCounter, 1)

	if m.FatalfFunc == nil {
		m.t.Fatal("Unexpected call to TesterMock.Fatalf")
		return
	}

	m.FatalfFunc(p, p1...)
}

//ValidateCallCounters checks that all mocked methods of the iterface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *TesterMock) ValidateCallCounters() {

	if m.ErrorFunc != nil && atomic.LoadUint64(&m.ErrorCounter) == 0 {
		m.t.Fatal("Expected call to TesterMock.Error")
	}

	if m.FatalFunc != nil && atomic.LoadUint64(&m.FatalCounter) == 0 {
		m.t.Fatal("Expected call to TesterMock.Fatal")
	}

	if m.FatalfFunc != nil && atomic.LoadUint64(&m.FatalfCounter) == 0 {
		m.t.Fatal("Expected call to TesterMock.Fatalf")
	}

}

//CheckMocksCalled checks that all mocked methods of the iterface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *TesterMock) CheckMocksCalled() {
	m.Finish()
}

//Finish checks that all mocked methods of the iterface have been called at least once
//Deprecated: please use MinimockFinish or use Finish method of minimock.Controller
func (m *TesterMock) Finish() {
	m.MinimockFinish()
}

//MinimockFinish checks that all mocked methods of the iterface have been called at least once
func (m *TesterMock) MinimockFinish() {

	if m.ErrorFunc != nil && atomic.LoadUint64(&m.ErrorCounter) == 0 {
		m.t.Fatal("Expected call to TesterMock.Error")
	}

	if m.FatalFunc != nil && atomic.LoadUint64(&m.FatalCounter) == 0 {
		m.t.Fatal("Expected call to TesterMock.Fatal")
	}

	if m.FatalfFunc != nil && atomic.LoadUint64(&m.FatalfCounter) == 0 {
		m.t.Fatal("Expected call to TesterMock.Fatalf")
	}

}

//Wait waits for all mocked methods to be called at least once
//Deprecated: please use MinimockWait or use Wait method of minimock.Controller
func (m *TesterMock) Wait(timeout time.Duration) {
	m.MinimockWait(timeout)
}

//MinimockWait waits for all mocked methods to be called at least once
//this method is called by minimock.Controller
func (m *TesterMock) MinimockWait(timeout time.Duration) {
	timeoutCh := time.After(timeout)
	for {
		ok := true
		ok = ok && (m.ErrorFunc == nil || atomic.LoadUint64(&m.ErrorCounter) > 0)
		ok = ok && (m.FatalFunc == nil || atomic.LoadUint64(&m.FatalCounter) > 0)
		ok = ok && (m.FatalfFunc == nil || atomic.LoadUint64(&m.FatalfCounter) > 0)

		if ok {
			return
		}

		select {
		case <-timeoutCh:

			if m.ErrorFunc != nil && atomic.LoadUint64(&m.ErrorCounter) == 0 {
				m.t.Error("Expected call to TesterMock.Error")
			}

			if m.FatalFunc != nil && atomic.LoadUint64(&m.FatalCounter) == 0 {
				m.t.Error("Expected call to TesterMock.Fatal")
			}

			if m.FatalfFunc != nil && atomic.LoadUint64(&m.FatalfCounter) == 0 {
				m.t.Error("Expected call to TesterMock.Fatalf")
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
func (m *TesterMock) AllMocksCalled() bool {

	if m.ErrorFunc != nil && atomic.LoadUint64(&m.ErrorCounter) == 0 {
		return false
	}

	if m.FatalFunc != nil && atomic.LoadUint64(&m.FatalCounter) == 0 {
		return false
	}

	if m.FatalfFunc != nil && atomic.LoadUint64(&m.FatalfCounter) == 0 {
		return false
	}

	return true
}
