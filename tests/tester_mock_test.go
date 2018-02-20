package tests

/*
DO NOT EDIT!
This code was generated automatically using github.com/gojuno/minimock v1.9
The original interface "Tester" can be found in github.com/gojuno/minimock
*/
import (
	"sync/atomic"
	"time"

	"github.com/gojuno/minimock"
	testify_assert "github.com/stretchr/testify/assert"
)

//TesterMock implements github.com/gojuno/minimock.Tester
type TesterMock struct {
	t minimock.Tester

	ErrorFunc       func(p ...interface{})
	ErrorCounter    uint64
	ErrorPreCounter uint64
	ErrorMock       mTesterMockError

	ErrorfFunc       func(p string, p1 ...interface{})
	ErrorfCounter    uint64
	ErrorfPreCounter uint64
	ErrorfMock       mTesterMockErrorf

	FatalFunc       func(p ...interface{})
	FatalCounter    uint64
	FatalPreCounter uint64
	FatalMock       mTesterMockFatal

	FatalfFunc       func(p string, p1 ...interface{})
	FatalfCounter    uint64
	FatalfPreCounter uint64
	FatalfMock       mTesterMockFatalf
}

//NewTesterMock returns a mock for github.com/gojuno/minimock.Tester
func NewTesterMock(t minimock.Tester) *TesterMock {
	m := &TesterMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.ErrorMock = mTesterMockError{mock: m}
	m.ErrorfMock = mTesterMockErrorf{mock: m}
	m.FatalMock = mTesterMockFatal{mock: m}
	m.FatalfMock = mTesterMockFatalf{mock: m}

	return m
}

type mTesterMockError struct {
	mock             *TesterMock
	mockExpectations *TesterMockErrorParams
}

//TesterMockErrorParams represents input parameters of the Tester.Error
type TesterMockErrorParams struct {
	p []interface{}
}

//Expect sets up expected params for the Tester.Error
func (m *mTesterMockError) Expect(p ...interface{}) *mTesterMockError {
	m.mockExpectations = &TesterMockErrorParams{p}
	return m
}

//Return sets up a mock for Tester.Error to return Return's arguments
func (m *mTesterMockError) Return() *TesterMock {
	m.mock.ErrorFunc = func(p ...interface{}) {
		return
	}
	return m.mock
}

//Set uses given function f as a mock of Tester.Error method
func (m *mTesterMockError) Set(f func(p ...interface{})) *TesterMock {
	m.mock.ErrorFunc = f
	return m.mock
}

//Error implements github.com/gojuno/minimock.Tester interface
func (m *TesterMock) Error(p ...interface{}) {
	atomic.AddUint64(&m.ErrorPreCounter, 1)
	defer atomic.AddUint64(&m.ErrorCounter, 1)

	if m.ErrorMock.mockExpectations != nil {
		testify_assert.Equal(m.t, *m.ErrorMock.mockExpectations, TesterMockErrorParams{p},
			"Tester.Error got unexpected parameters")

		if m.ErrorFunc == nil {

			m.t.Fatal("No results are set for the TesterMock.Error")

			return
		}
	}

	if m.ErrorFunc == nil {
		m.t.Fatal("Unexpected call to TesterMock.Error")
		return
	}

	m.ErrorFunc(p...)
}

//ErrorMinimockCounter returns a count of TesterMock.ErrorFunc invocations
func (m *TesterMock) ErrorMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.ErrorCounter)
}

//ErrorMinimockPreCounter returns the value of TesterMock.Error invocations
func (m *TesterMock) ErrorMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.ErrorPreCounter)
}

type mTesterMockErrorf struct {
	mock             *TesterMock
	mockExpectations *TesterMockErrorfParams
}

//TesterMockErrorfParams represents input parameters of the Tester.Errorf
type TesterMockErrorfParams struct {
	p  string
	p1 []interface{}
}

//Expect sets up expected params for the Tester.Errorf
func (m *mTesterMockErrorf) Expect(p string, p1 ...interface{}) *mTesterMockErrorf {
	m.mockExpectations = &TesterMockErrorfParams{p, p1}
	return m
}

//Return sets up a mock for Tester.Errorf to return Return's arguments
func (m *mTesterMockErrorf) Return() *TesterMock {
	m.mock.ErrorfFunc = func(p string, p1 ...interface{}) {
		return
	}
	return m.mock
}

//Set uses given function f as a mock of Tester.Errorf method
func (m *mTesterMockErrorf) Set(f func(p string, p1 ...interface{})) *TesterMock {
	m.mock.ErrorfFunc = f
	return m.mock
}

//Errorf implements github.com/gojuno/minimock.Tester interface
func (m *TesterMock) Errorf(p string, p1 ...interface{}) {
	atomic.AddUint64(&m.ErrorfPreCounter, 1)
	defer atomic.AddUint64(&m.ErrorfCounter, 1)

	if m.ErrorfMock.mockExpectations != nil {
		testify_assert.Equal(m.t, *m.ErrorfMock.mockExpectations, TesterMockErrorfParams{p, p1},
			"Tester.Errorf got unexpected parameters")

		if m.ErrorfFunc == nil {

			m.t.Fatal("No results are set for the TesterMock.Errorf")

			return
		}
	}

	if m.ErrorfFunc == nil {
		m.t.Fatal("Unexpected call to TesterMock.Errorf")
		return
	}

	m.ErrorfFunc(p, p1...)
}

//ErrorfMinimockCounter returns a count of TesterMock.ErrorfFunc invocations
func (m *TesterMock) ErrorfMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.ErrorfCounter)
}

//ErrorfMinimockPreCounter returns the value of TesterMock.Errorf invocations
func (m *TesterMock) ErrorfMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.ErrorfPreCounter)
}

type mTesterMockFatal struct {
	mock             *TesterMock
	mockExpectations *TesterMockFatalParams
}

//TesterMockFatalParams represents input parameters of the Tester.Fatal
type TesterMockFatalParams struct {
	p []interface{}
}

//Expect sets up expected params for the Tester.Fatal
func (m *mTesterMockFatal) Expect(p ...interface{}) *mTesterMockFatal {
	m.mockExpectations = &TesterMockFatalParams{p}
	return m
}

//Return sets up a mock for Tester.Fatal to return Return's arguments
func (m *mTesterMockFatal) Return() *TesterMock {
	m.mock.FatalFunc = func(p ...interface{}) {
		return
	}
	return m.mock
}

//Set uses given function f as a mock of Tester.Fatal method
func (m *mTesterMockFatal) Set(f func(p ...interface{})) *TesterMock {
	m.mock.FatalFunc = f
	return m.mock
}

//Fatal implements github.com/gojuno/minimock.Tester interface
func (m *TesterMock) Fatal(p ...interface{}) {
	atomic.AddUint64(&m.FatalPreCounter, 1)
	defer atomic.AddUint64(&m.FatalCounter, 1)

	if m.FatalMock.mockExpectations != nil {
		testify_assert.Equal(m.t, *m.FatalMock.mockExpectations, TesterMockFatalParams{p},
			"Tester.Fatal got unexpected parameters")

		if m.FatalFunc == nil {

			m.t.Fatal("No results are set for the TesterMock.Fatal")

			return
		}
	}

	if m.FatalFunc == nil {
		m.t.Fatal("Unexpected call to TesterMock.Fatal")
		return
	}

	m.FatalFunc(p...)
}

//FatalMinimockCounter returns a count of TesterMock.FatalFunc invocations
func (m *TesterMock) FatalMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.FatalCounter)
}

//FatalMinimockPreCounter returns the value of TesterMock.Fatal invocations
func (m *TesterMock) FatalMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.FatalPreCounter)
}

type mTesterMockFatalf struct {
	mock             *TesterMock
	mockExpectations *TesterMockFatalfParams
}

//TesterMockFatalfParams represents input parameters of the Tester.Fatalf
type TesterMockFatalfParams struct {
	p  string
	p1 []interface{}
}

//Expect sets up expected params for the Tester.Fatalf
func (m *mTesterMockFatalf) Expect(p string, p1 ...interface{}) *mTesterMockFatalf {
	m.mockExpectations = &TesterMockFatalfParams{p, p1}
	return m
}

//Return sets up a mock for Tester.Fatalf to return Return's arguments
func (m *mTesterMockFatalf) Return() *TesterMock {
	m.mock.FatalfFunc = func(p string, p1 ...interface{}) {
		return
	}
	return m.mock
}

//Set uses given function f as a mock of Tester.Fatalf method
func (m *mTesterMockFatalf) Set(f func(p string, p1 ...interface{})) *TesterMock {
	m.mock.FatalfFunc = f
	return m.mock
}

//Fatalf implements github.com/gojuno/minimock.Tester interface
func (m *TesterMock) Fatalf(p string, p1 ...interface{}) {
	atomic.AddUint64(&m.FatalfPreCounter, 1)
	defer atomic.AddUint64(&m.FatalfCounter, 1)

	if m.FatalfMock.mockExpectations != nil {
		testify_assert.Equal(m.t, *m.FatalfMock.mockExpectations, TesterMockFatalfParams{p, p1},
			"Tester.Fatalf got unexpected parameters")

		if m.FatalfFunc == nil {

			m.t.Fatal("No results are set for the TesterMock.Fatalf")

			return
		}
	}

	if m.FatalfFunc == nil {
		m.t.Fatal("Unexpected call to TesterMock.Fatalf")
		return
	}

	m.FatalfFunc(p, p1...)
}

//FatalfMinimockCounter returns a count of TesterMock.FatalfFunc invocations
func (m *TesterMock) FatalfMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.FatalfCounter)
}

//FatalfMinimockPreCounter returns the value of TesterMock.Fatalf invocations
func (m *TesterMock) FatalfMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.FatalfPreCounter)
}

//ValidateCallCounters checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *TesterMock) ValidateCallCounters() {

	if m.ErrorFunc != nil && atomic.LoadUint64(&m.ErrorCounter) == 0 {
		m.t.Fatal("Expected call to TesterMock.Error")
	}

	if m.ErrorfFunc != nil && atomic.LoadUint64(&m.ErrorfCounter) == 0 {
		m.t.Fatal("Expected call to TesterMock.Errorf")
	}

	if m.FatalFunc != nil && atomic.LoadUint64(&m.FatalCounter) == 0 {
		m.t.Fatal("Expected call to TesterMock.Fatal")
	}

	if m.FatalfFunc != nil && atomic.LoadUint64(&m.FatalfCounter) == 0 {
		m.t.Fatal("Expected call to TesterMock.Fatalf")
	}

}

//CheckMocksCalled checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *TesterMock) CheckMocksCalled() {
	m.Finish()
}

//Finish checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish or use Finish method of minimock.Controller
func (m *TesterMock) Finish() {
	m.MinimockFinish()
}

//MinimockFinish checks that all mocked methods of the interface have been called at least once
func (m *TesterMock) MinimockFinish() {

	if m.ErrorFunc != nil && atomic.LoadUint64(&m.ErrorCounter) == 0 {
		m.t.Fatal("Expected call to TesterMock.Error")
	}

	if m.ErrorfFunc != nil && atomic.LoadUint64(&m.ErrorfCounter) == 0 {
		m.t.Fatal("Expected call to TesterMock.Errorf")
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
		ok = ok && (m.ErrorfFunc == nil || atomic.LoadUint64(&m.ErrorfCounter) > 0)
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

			if m.ErrorfFunc != nil && atomic.LoadUint64(&m.ErrorfCounter) == 0 {
				m.t.Error("Expected call to TesterMock.Errorf")
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

	if m.ErrorfFunc != nil && atomic.LoadUint64(&m.ErrorfCounter) == 0 {
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
