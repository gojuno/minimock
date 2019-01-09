package tests

/*
DO NOT EDIT!
This code was generated automatically using github.com/gojuno/minimock v1.9.1
The original interface "Tester" can be found in github.com/gojuno/minimock
*/
import (
	"sync/atomic"
	"time"

	"github.com/gojuno/minimock"
	testify_assert "github.com/stretchr/testify/assert"
)

// TesterMock implements github.com/gojuno/minimock.Tester
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

// NewTesterMock returns a mock for github.com/gojuno/minimock.Tester
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
	mock              *TesterMock
	mainExpectation   *TesterMockErrorExpectation
	expectationSeries []*TesterMockErrorExpectation
}

// TesterMockErrorExpectation specifies expectation struct of the Tester.Error
type TesterMockErrorExpectation struct {
	input *TesterMockErrorInput
}

// TesterMockErrorInput represents input parameters of the Tester.Error
type TesterMockErrorInput struct {
	p []interface{}
}

// Expect specifies that invocation of Tester.Error is expected from 1 to Infinity times
func (m *mTesterMockError) Expect(p ...interface{}) *mTesterMockError {
	m.mock.ErrorFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &TesterMockErrorExpectation{}
	}
	m.mainExpectation.input = &TesterMockErrorInput{p}
	return m
}

// Return specifies results of invocation of Tester.Error
func (m *mTesterMockError) Return() *TesterMock {
	m.mock.ErrorFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &TesterMockErrorExpectation{}
	}

	return m.mock
}

// ExpectOnce specifies that invocation of Tester.Error is expected once
func (m *mTesterMockError) ExpectOnce(p ...interface{}) *TesterMockErrorExpectation {
	m.mock.ErrorFunc = nil
	m.mainExpectation = nil

	expectation := &TesterMockErrorExpectation{}
	expectation.input = &TesterMockErrorInput{p}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

// Set uses given function f as a mock of Tester.Error method
func (m *mTesterMockError) Set(f func(p ...interface{})) *TesterMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.ErrorFunc = f
	return m.mock
}

// Error implements github.com/gojuno/minimock.Tester interface
func (m *TesterMock) Error(p ...interface{}) {
	counter := atomic.AddUint64(&m.ErrorPreCounter, 1)
	defer atomic.AddUint64(&m.ErrorCounter, 1)

	if len(m.ErrorMock.expectationSeries) > 0 {
		if counter > uint64(len(m.ErrorMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to TesterMock.Error. %v", p)
			return
		}

		input := m.ErrorMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, TesterMockErrorInput{p}, "Tester.Error got unexpected parameters")

		return
	}

	if m.ErrorMock.mainExpectation != nil {

		input := m.ErrorMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, TesterMockErrorInput{p}, "Tester.Error got unexpected parameters")
		}

		return
	}

	if m.ErrorFunc == nil {
		m.t.Fatalf("Unexpected call to TesterMock.Error. %v", p)
		return
	}

	m.ErrorFunc(p...)
}

// ErrorMinimockCounter returns a count of TesterMock.ErrorFunc invocations
func (m *TesterMock) ErrorMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.ErrorCounter)
}

// ErrorMinimockPreCounter returns the value of TesterMock.Error invocations
func (m *TesterMock) ErrorMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.ErrorPreCounter)
}

// ErrorFinished returns true if mock invocations count is ok
func (m *TesterMock) ErrorFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.ErrorMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.ErrorCounter) == uint64(len(m.ErrorMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.ErrorMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.ErrorCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.ErrorFunc != nil {
		return atomic.LoadUint64(&m.ErrorCounter) > 0
	}

	return true
}

type mTesterMockErrorf struct {
	mock              *TesterMock
	mainExpectation   *TesterMockErrorfExpectation
	expectationSeries []*TesterMockErrorfExpectation
}

// TesterMockErrorfExpectation specifies expectation struct of the Tester.Errorf
type TesterMockErrorfExpectation struct {
	input *TesterMockErrorfInput
}

// TesterMockErrorfInput represents input parameters of the Tester.Errorf
type TesterMockErrorfInput struct {
	p  string
	p1 []interface{}
}

// Expect specifies that invocation of Tester.Errorf is expected from 1 to Infinity times
func (m *mTesterMockErrorf) Expect(p string, p1 ...interface{}) *mTesterMockErrorf {
	m.mock.ErrorfFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &TesterMockErrorfExpectation{}
	}
	m.mainExpectation.input = &TesterMockErrorfInput{p, p1}
	return m
}

// Return specifies results of invocation of Tester.Errorf
func (m *mTesterMockErrorf) Return() *TesterMock {
	m.mock.ErrorfFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &TesterMockErrorfExpectation{}
	}

	return m.mock
}

// ExpectOnce specifies that invocation of Tester.Errorf is expected once
func (m *mTesterMockErrorf) ExpectOnce(p string, p1 ...interface{}) *TesterMockErrorfExpectation {
	m.mock.ErrorfFunc = nil
	m.mainExpectation = nil

	expectation := &TesterMockErrorfExpectation{}
	expectation.input = &TesterMockErrorfInput{p, p1}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

// Set uses given function f as a mock of Tester.Errorf method
func (m *mTesterMockErrorf) Set(f func(p string, p1 ...interface{})) *TesterMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.ErrorfFunc = f
	return m.mock
}

// Errorf implements github.com/gojuno/minimock.Tester interface
func (m *TesterMock) Errorf(p string, p1 ...interface{}) {
	counter := atomic.AddUint64(&m.ErrorfPreCounter, 1)
	defer atomic.AddUint64(&m.ErrorfCounter, 1)

	if len(m.ErrorfMock.expectationSeries) > 0 {
		if counter > uint64(len(m.ErrorfMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to TesterMock.Errorf. %v %v", p, p1)
			return
		}

		input := m.ErrorfMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, TesterMockErrorfInput{p, p1}, "Tester.Errorf got unexpected parameters")

		return
	}

	if m.ErrorfMock.mainExpectation != nil {

		input := m.ErrorfMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, TesterMockErrorfInput{p, p1}, "Tester.Errorf got unexpected parameters")
		}

		return
	}

	if m.ErrorfFunc == nil {
		m.t.Fatalf("Unexpected call to TesterMock.Errorf. %v %v", p, p1)
		return
	}

	m.ErrorfFunc(p, p1...)
}

// ErrorfMinimockCounter returns a count of TesterMock.ErrorfFunc invocations
func (m *TesterMock) ErrorfMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.ErrorfCounter)
}

// ErrorfMinimockPreCounter returns the value of TesterMock.Errorf invocations
func (m *TesterMock) ErrorfMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.ErrorfPreCounter)
}

// ErrorfFinished returns true if mock invocations count is ok
func (m *TesterMock) ErrorfFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.ErrorfMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.ErrorfCounter) == uint64(len(m.ErrorfMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.ErrorfMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.ErrorfCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.ErrorfFunc != nil {
		return atomic.LoadUint64(&m.ErrorfCounter) > 0
	}

	return true
}

type mTesterMockFatal struct {
	mock              *TesterMock
	mainExpectation   *TesterMockFatalExpectation
	expectationSeries []*TesterMockFatalExpectation
}

// TesterMockFatalExpectation specifies expectation struct of the Tester.Fatal
type TesterMockFatalExpectation struct {
	input *TesterMockFatalInput
}

// TesterMockFatalInput represents input parameters of the Tester.Fatal
type TesterMockFatalInput struct {
	p []interface{}
}

// Expect specifies that invocation of Tester.Fatal is expected from 1 to Infinity times
func (m *mTesterMockFatal) Expect(p ...interface{}) *mTesterMockFatal {
	m.mock.FatalFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &TesterMockFatalExpectation{}
	}
	m.mainExpectation.input = &TesterMockFatalInput{p}
	return m
}

// Return specifies results of invocation of Tester.Fatal
func (m *mTesterMockFatal) Return() *TesterMock {
	m.mock.FatalFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &TesterMockFatalExpectation{}
	}

	return m.mock
}

// ExpectOnce specifies that invocation of Tester.Fatal is expected once
func (m *mTesterMockFatal) ExpectOnce(p ...interface{}) *TesterMockFatalExpectation {
	m.mock.FatalFunc = nil
	m.mainExpectation = nil

	expectation := &TesterMockFatalExpectation{}
	expectation.input = &TesterMockFatalInput{p}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

// Set uses given function f as a mock of Tester.Fatal method
func (m *mTesterMockFatal) Set(f func(p ...interface{})) *TesterMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.FatalFunc = f
	return m.mock
}

// Fatal implements github.com/gojuno/minimock.Tester interface
func (m *TesterMock) Fatal(p ...interface{}) {
	counter := atomic.AddUint64(&m.FatalPreCounter, 1)
	defer atomic.AddUint64(&m.FatalCounter, 1)

	if len(m.FatalMock.expectationSeries) > 0 {
		if counter > uint64(len(m.FatalMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to TesterMock.Fatal. %v", p)
			return
		}

		input := m.FatalMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, TesterMockFatalInput{p}, "Tester.Fatal got unexpected parameters")

		return
	}

	if m.FatalMock.mainExpectation != nil {

		input := m.FatalMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, TesterMockFatalInput{p}, "Tester.Fatal got unexpected parameters")
		}

		return
	}

	if m.FatalFunc == nil {
		m.t.Fatalf("Unexpected call to TesterMock.Fatal. %v", p)
		return
	}

	m.FatalFunc(p...)
}

// FatalMinimockCounter returns a count of TesterMock.FatalFunc invocations
func (m *TesterMock) FatalMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.FatalCounter)
}

// FatalMinimockPreCounter returns the value of TesterMock.Fatal invocations
func (m *TesterMock) FatalMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.FatalPreCounter)
}

// FatalFinished returns true if mock invocations count is ok
func (m *TesterMock) FatalFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.FatalMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.FatalCounter) == uint64(len(m.FatalMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.FatalMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.FatalCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.FatalFunc != nil {
		return atomic.LoadUint64(&m.FatalCounter) > 0
	}

	return true
}

type mTesterMockFatalf struct {
	mock              *TesterMock
	mainExpectation   *TesterMockFatalfExpectation
	expectationSeries []*TesterMockFatalfExpectation
}

// TesterMockFatalfExpectation specifies expectation struct of the Tester.Fatalf
type TesterMockFatalfExpectation struct {
	input *TesterMockFatalfInput
}

// TesterMockFatalfInput represents input parameters of the Tester.Fatalf
type TesterMockFatalfInput struct {
	p  string
	p1 []interface{}
}

// Expect specifies that invocation of Tester.Fatalf is expected from 1 to Infinity times
func (m *mTesterMockFatalf) Expect(p string, p1 ...interface{}) *mTesterMockFatalf {
	m.mock.FatalfFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &TesterMockFatalfExpectation{}
	}
	m.mainExpectation.input = &TesterMockFatalfInput{p, p1}
	return m
}

// Return specifies results of invocation of Tester.Fatalf
func (m *mTesterMockFatalf) Return() *TesterMock {
	m.mock.FatalfFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &TesterMockFatalfExpectation{}
	}

	return m.mock
}

// ExpectOnce specifies that invocation of Tester.Fatalf is expected once
func (m *mTesterMockFatalf) ExpectOnce(p string, p1 ...interface{}) *TesterMockFatalfExpectation {
	m.mock.FatalfFunc = nil
	m.mainExpectation = nil

	expectation := &TesterMockFatalfExpectation{}
	expectation.input = &TesterMockFatalfInput{p, p1}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

// Set uses given function f as a mock of Tester.Fatalf method
func (m *mTesterMockFatalf) Set(f func(p string, p1 ...interface{})) *TesterMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.FatalfFunc = f
	return m.mock
}

// Fatalf implements github.com/gojuno/minimock.Tester interface
func (m *TesterMock) Fatalf(p string, p1 ...interface{}) {
	counter := atomic.AddUint64(&m.FatalfPreCounter, 1)
	defer atomic.AddUint64(&m.FatalfCounter, 1)

	if len(m.FatalfMock.expectationSeries) > 0 {
		if counter > uint64(len(m.FatalfMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to TesterMock.Fatalf. %v %v", p, p1)
			return
		}

		input := m.FatalfMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, TesterMockFatalfInput{p, p1}, "Tester.Fatalf got unexpected parameters")

		return
	}

	if m.FatalfMock.mainExpectation != nil {

		input := m.FatalfMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, TesterMockFatalfInput{p, p1}, "Tester.Fatalf got unexpected parameters")
		}

		return
	}

	if m.FatalfFunc == nil {
		m.t.Fatalf("Unexpected call to TesterMock.Fatalf. %v %v", p, p1)
		return
	}

	m.FatalfFunc(p, p1...)
}

// FatalfMinimockCounter returns a count of TesterMock.FatalfFunc invocations
func (m *TesterMock) FatalfMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.FatalfCounter)
}

// FatalfMinimockPreCounter returns the value of TesterMock.Fatalf invocations
func (m *TesterMock) FatalfMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.FatalfPreCounter)
}

// FatalfFinished returns true if mock invocations count is ok
func (m *TesterMock) FatalfFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.FatalfMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.FatalfCounter) == uint64(len(m.FatalfMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.FatalfMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.FatalfCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.FatalfFunc != nil {
		return atomic.LoadUint64(&m.FatalfCounter) > 0
	}

	return true
}

// ValidateCallCounters checks that all mocked methods of the interface have been called at least once
// Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *TesterMock) ValidateCallCounters() {

	if !m.ErrorFinished() {
		m.t.Fatal("Expected call to TesterMock.Error")
	}

	if !m.ErrorfFinished() {
		m.t.Fatal("Expected call to TesterMock.Errorf")
	}

	if !m.FatalFinished() {
		m.t.Fatal("Expected call to TesterMock.Fatal")
	}

	if !m.FatalfFinished() {
		m.t.Fatal("Expected call to TesterMock.Fatalf")
	}

}

// CheckMocksCalled checks that all mocked methods of the interface have been called at least once
// Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *TesterMock) CheckMocksCalled() {
	m.Finish()
}

// Finish checks that all mocked methods of the interface have been called at least once
// Deprecated: please use MinimockFinish or use Finish method of minimock.Controller
func (m *TesterMock) Finish() {
	m.MinimockFinish()
}

// MinimockFinish checks that all mocked methods of the interface have been called at least once
func (m *TesterMock) MinimockFinish() {

	if !m.ErrorFinished() {
		m.t.Fatal("Expected call to TesterMock.Error")
	}

	if !m.ErrorfFinished() {
		m.t.Fatal("Expected call to TesterMock.Errorf")
	}

	if !m.FatalFinished() {
		m.t.Fatal("Expected call to TesterMock.Fatal")
	}

	if !m.FatalfFinished() {
		m.t.Fatal("Expected call to TesterMock.Fatalf")
	}

}

// Wait waits for all mocked methods to be called at least once
// Deprecated: please use MinimockWait or use Wait method of minimock.Controller
func (m *TesterMock) Wait(timeout time.Duration) {
	m.MinimockWait(timeout)
}

// MinimockWait waits for all mocked methods to be called at least once
// this method is called by minimock.Controller
func (m *TesterMock) MinimockWait(timeout time.Duration) {
	timeoutCh := time.After(timeout)
	for {
		ok := true
		ok = ok && m.ErrorFinished()
		ok = ok && m.ErrorfFinished()
		ok = ok && m.FatalFinished()
		ok = ok && m.FatalfFinished()

		if ok {
			return
		}

		select {
		case <-timeoutCh:

			if !m.ErrorFinished() {
				m.t.Error("Expected call to TesterMock.Error")
			}

			if !m.ErrorfFinished() {
				m.t.Error("Expected call to TesterMock.Errorf")
			}

			if !m.FatalFinished() {
				m.t.Error("Expected call to TesterMock.Fatal")
			}

			if !m.FatalfFinished() {
				m.t.Error("Expected call to TesterMock.Fatalf")
			}

			m.t.Fatalf("Some mocks were not called on time: %s", timeout)
			return
		default:
			time.Sleep(time.Millisecond)
		}
	}
}

// AllMocksCalled returns true if all mocked methods were called before the execution of AllMocksCalled,
// it can be used with assert/require, i.e. assert.True(mock.AllMocksCalled())
func (m *TesterMock) AllMocksCalled() bool {

	if !m.ErrorFinished() {
		return false
	}

	if !m.ErrorfFinished() {
		return false
	}

	if !m.FatalFinished() {
		return false
	}

	if !m.FatalfFinished() {
		return false
	}

	return true
}
