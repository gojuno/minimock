package tests

// DO NOT EDIT!
// The code below was generated with http://github.com/deff7/minimock (dev)

//go:generate minimock -i github.com/deff7/minimock/pkg.Tester -o ./tests/tester_mock_test.go

import (
	"sync/atomic"
	"time"

	minimock "github.com/deff7/minimock/pkg"
)

// TesterMock implements minimock.Tester
type TesterMock struct {
	t minimock.Tester

	funcError          func(p1 ...interface{})
	afterErrorCounter  uint64
	beforeErrorCounter uint64
	ErrorMock          mTesterMockError

	funcErrorf          func(format string, args ...interface{})
	afterErrorfCounter  uint64
	beforeErrorfCounter uint64
	ErrorfMock          mTesterMockErrorf

	funcFailNow          func()
	afterFailNowCounter  uint64
	beforeFailNowCounter uint64
	FailNowMock          mTesterMockFailNow

	funcFatal          func(args ...interface{})
	afterFatalCounter  uint64
	beforeFatalCounter uint64
	FatalMock          mTesterMockFatal

	funcFatalf          func(format string, args ...interface{})
	afterFatalfCounter  uint64
	beforeFatalfCounter uint64
	FatalfMock          mTesterMockFatalf
}

// NewTesterMock returns a mock for minimock.Tester
func NewTesterMock(t minimock.Tester) *TesterMock {
	m := &TesterMock{t: t}
	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}
	m.ErrorMock = mTesterMockError{mock: m}
	m.ErrorfMock = mTesterMockErrorf{mock: m}
	m.FailNowMock = mTesterMockFailNow{mock: m}
	m.FatalMock = mTesterMockFatal{mock: m}
	m.FatalfMock = mTesterMockFatalf{mock: m}

	return m
}

type mTesterMockError struct {
	mock               *TesterMock
	defaultExpectation *TesterMockErrorExpectation
	expectations       []*TesterMockErrorExpectation
}

// TesterMockErrorExpectation specifies expectation struct of the Tester.Error
type TesterMockErrorExpectation struct {
	mock   *TesterMock
	params *TesterMockErrorParams

	Counter uint64
}

// TesterMockErrorParams contains parameters of the Tester.Error
type TesterMockErrorParams struct {
	p1 []interface{}
}

// Expect sets up expected params for Tester.Error
func (m *mTesterMockError) Expect(p1 ...interface{}) *mTesterMockError {
	if m.mock.funcError != nil {
		m.mock.t.Fatalf("TesterMock.Error mock is already set by Set")
	}

	if m.defaultExpectation == nil {
		m.defaultExpectation = &TesterMockErrorExpectation{}
	}

	m.defaultExpectation.params = &TesterMockErrorParams{p1}
	for _, e := range m.expectations {
		if minimock.Equal(e.params, m.defaultExpectation.params) {
			m.mock.t.Fatalf("Expectation set by When has same params: %#v", *m.defaultExpectation.params)
		}
	}

	return m
}

// Return sets up results that will be returned by Tester.Error
func (m *mTesterMockError) Return() *TesterMock {
	if m.mock.funcError != nil {
		m.mock.t.Fatalf("TesterMock.Error mock is already set by Set")
	}

	if m.defaultExpectation == nil {
		m.defaultExpectation = &TesterMockErrorExpectation{mock: m.mock}
	}

	return m.mock
}

//Set uses given function f to mock the Tester.Error method
func (m *mTesterMockError) Set(f func(p1 ...interface{})) *TesterMock {
	if m.defaultExpectation != nil {
		m.mock.t.Fatalf("Default expectation is already set for the Tester.Error method")
	}

	if len(m.expectations) > 0 {
		m.mock.t.Fatalf("Some expectations are already set for the Tester.Error method")
	}

	m.mock.funcError = f
	return m.mock
}

// Error implements minimock.Tester
func (m *TesterMock) Error(p1 ...interface{}) {
	atomic.AddUint64(&m.beforeErrorCounter, 1)
	defer atomic.AddUint64(&m.afterErrorCounter, 1)

	for _, e := range m.ErrorMock.expectations {
		if minimock.Equal(*e.params, TesterMockErrorParams{p1}) {
			atomic.AddUint64(&e.Counter, 1)
			return
		}
	}

	if m.ErrorMock.defaultExpectation != nil {
		atomic.AddUint64(&m.ErrorMock.defaultExpectation.Counter, 1)
		want := m.ErrorMock.defaultExpectation.params
		got := TesterMockErrorParams{p1}
		if want != nil && !minimock.Equal(*want, got) {
			m.t.Errorf("TesterMock.Error got unexpected parameters, want: %#v, got: %#v\n", *want, got)
		}

		return

	}
	if m.funcError != nil {
		m.funcError(p1...)
		return
	}
	m.t.Fatalf("Unexpected call to TesterMock.Error. %v", p1)

}

// ErrorAfterCounter returns a count of finished TesterMock.Error invocations
func (m *TesterMock) ErrorAfterCounter() uint64 {
	return atomic.LoadUint64(&m.afterErrorCounter)
}

// ErrorBeforeCounter returns a count of TesterMock.Error invocations
func (m *TesterMock) ErrorBeforeCounter() uint64 {
	return atomic.LoadUint64(&m.beforeErrorCounter)
}

// MinimockErrorDone returns true if the count of the Error invocations corresponds
// the number of defined expectations
func (m *TesterMock) MinimockErrorDone() bool {
	for _, e := range m.ErrorMock.expectations {
		if atomic.LoadUint64(&e.Counter) < 1 {
			return false
		}
	}

	// if default expectation was set then invocations count should be greater than zero
	if m.ErrorMock.defaultExpectation != nil && atomic.LoadUint64(&m.afterErrorCounter) < 1 {
		return false
	}
	// if func was set then invocations count should be greater than zero
	if m.funcError != nil && atomic.LoadUint64(&m.afterErrorCounter) < 1 {
		return false
	}
	return true
}

// MinimockErrorInspect logs each unmet expectation
func (m *TesterMock) MinimockErrorInspect() {
	for _, e := range m.ErrorMock.expectations {
		if atomic.LoadUint64(&e.Counter) < 1 {
			m.t.Errorf("Expected call to TesterMock.Error with params: %#v", *e.params)
		}
	}

	// if default expectation was set then invocations count should be greater than zero
	if m.ErrorMock.defaultExpectation != nil && atomic.LoadUint64(&m.afterErrorCounter) < 1 {
		m.t.Errorf("Expected call to TesterMock.Error with params: %#v", *m.ErrorMock.defaultExpectation.params)
	}
	// if func was set then invocations count should be greater than zero
	if m.funcError != nil && atomic.LoadUint64(&m.afterErrorCounter) < 1 {
		m.t.Error("Expected call to TesterMock.Error")
	}
}

type mTesterMockErrorf struct {
	mock               *TesterMock
	defaultExpectation *TesterMockErrorfExpectation
	expectations       []*TesterMockErrorfExpectation
}

// TesterMockErrorfExpectation specifies expectation struct of the Tester.Errorf
type TesterMockErrorfExpectation struct {
	mock   *TesterMock
	params *TesterMockErrorfParams

	Counter uint64
}

// TesterMockErrorfParams contains parameters of the Tester.Errorf
type TesterMockErrorfParams struct {
	format string
	args   []interface{}
}

// Expect sets up expected params for Tester.Errorf
func (m *mTesterMockErrorf) Expect(format string, args ...interface{}) *mTesterMockErrorf {
	if m.mock.funcErrorf != nil {
		m.mock.t.Fatalf("TesterMock.Errorf mock is already set by Set")
	}

	if m.defaultExpectation == nil {
		m.defaultExpectation = &TesterMockErrorfExpectation{}
	}

	m.defaultExpectation.params = &TesterMockErrorfParams{format, args}
	for _, e := range m.expectations {
		if minimock.Equal(e.params, m.defaultExpectation.params) {
			m.mock.t.Fatalf("Expectation set by When has same params: %#v", *m.defaultExpectation.params)
		}
	}

	return m
}

// Return sets up results that will be returned by Tester.Errorf
func (m *mTesterMockErrorf) Return() *TesterMock {
	if m.mock.funcErrorf != nil {
		m.mock.t.Fatalf("TesterMock.Errorf mock is already set by Set")
	}

	if m.defaultExpectation == nil {
		m.defaultExpectation = &TesterMockErrorfExpectation{mock: m.mock}
	}

	return m.mock
}

//Set uses given function f to mock the Tester.Errorf method
func (m *mTesterMockErrorf) Set(f func(format string, args ...interface{})) *TesterMock {
	if m.defaultExpectation != nil {
		m.mock.t.Fatalf("Default expectation is already set for the Tester.Errorf method")
	}

	if len(m.expectations) > 0 {
		m.mock.t.Fatalf("Some expectations are already set for the Tester.Errorf method")
	}

	m.mock.funcErrorf = f
	return m.mock
}

// Errorf implements minimock.Tester
func (m *TesterMock) Errorf(format string, args ...interface{}) {
	atomic.AddUint64(&m.beforeErrorfCounter, 1)
	defer atomic.AddUint64(&m.afterErrorfCounter, 1)

	for _, e := range m.ErrorfMock.expectations {
		if minimock.Equal(*e.params, TesterMockErrorfParams{format, args}) {
			atomic.AddUint64(&e.Counter, 1)
			return
		}
	}

	if m.ErrorfMock.defaultExpectation != nil {
		atomic.AddUint64(&m.ErrorfMock.defaultExpectation.Counter, 1)
		want := m.ErrorfMock.defaultExpectation.params
		got := TesterMockErrorfParams{format, args}
		if want != nil && !minimock.Equal(*want, got) {
			m.t.Errorf("TesterMock.Errorf got unexpected parameters, want: %#v, got: %#v\n", *want, got)
		}

		return

	}
	if m.funcErrorf != nil {
		m.funcErrorf(format, args...)
		return
	}
	m.t.Fatalf("Unexpected call to TesterMock.Errorf. %v %v", format, args)

}

// ErrorfAfterCounter returns a count of finished TesterMock.Errorf invocations
func (m *TesterMock) ErrorfAfterCounter() uint64 {
	return atomic.LoadUint64(&m.afterErrorfCounter)
}

// ErrorfBeforeCounter returns a count of TesterMock.Errorf invocations
func (m *TesterMock) ErrorfBeforeCounter() uint64 {
	return atomic.LoadUint64(&m.beforeErrorfCounter)
}

// MinimockErrorfDone returns true if the count of the Errorf invocations corresponds
// the number of defined expectations
func (m *TesterMock) MinimockErrorfDone() bool {
	for _, e := range m.ErrorfMock.expectations {
		if atomic.LoadUint64(&e.Counter) < 1 {
			return false
		}
	}

	// if default expectation was set then invocations count should be greater than zero
	if m.ErrorfMock.defaultExpectation != nil && atomic.LoadUint64(&m.afterErrorfCounter) < 1 {
		return false
	}
	// if func was set then invocations count should be greater than zero
	if m.funcErrorf != nil && atomic.LoadUint64(&m.afterErrorfCounter) < 1 {
		return false
	}
	return true
}

// MinimockErrorfInspect logs each unmet expectation
func (m *TesterMock) MinimockErrorfInspect() {
	for _, e := range m.ErrorfMock.expectations {
		if atomic.LoadUint64(&e.Counter) < 1 {
			m.t.Errorf("Expected call to TesterMock.Errorf with params: %#v", *e.params)
		}
	}

	// if default expectation was set then invocations count should be greater than zero
	if m.ErrorfMock.defaultExpectation != nil && atomic.LoadUint64(&m.afterErrorfCounter) < 1 {
		m.t.Errorf("Expected call to TesterMock.Errorf with params: %#v", *m.ErrorfMock.defaultExpectation.params)
	}
	// if func was set then invocations count should be greater than zero
	if m.funcErrorf != nil && atomic.LoadUint64(&m.afterErrorfCounter) < 1 {
		m.t.Error("Expected call to TesterMock.Errorf")
	}
}

type mTesterMockFailNow struct {
	mock               *TesterMock
	defaultExpectation *TesterMockFailNowExpectation
	expectations       []*TesterMockFailNowExpectation
}

// TesterMockFailNowExpectation specifies expectation struct of the Tester.FailNow
type TesterMockFailNowExpectation struct {
	mock *TesterMock

	Counter uint64
}

// Expect sets up expected params for Tester.FailNow
func (m *mTesterMockFailNow) Expect() *mTesterMockFailNow {
	if m.mock.funcFailNow != nil {
		m.mock.t.Fatalf("TesterMock.FailNow mock is already set by Set")
	}

	if m.defaultExpectation == nil {
		m.defaultExpectation = &TesterMockFailNowExpectation{}
	}

	return m
}

// Return sets up results that will be returned by Tester.FailNow
func (m *mTesterMockFailNow) Return() *TesterMock {
	if m.mock.funcFailNow != nil {
		m.mock.t.Fatalf("TesterMock.FailNow mock is already set by Set")
	}

	if m.defaultExpectation == nil {
		m.defaultExpectation = &TesterMockFailNowExpectation{mock: m.mock}
	}

	return m.mock
}

//Set uses given function f to mock the Tester.FailNow method
func (m *mTesterMockFailNow) Set(f func()) *TesterMock {
	if m.defaultExpectation != nil {
		m.mock.t.Fatalf("Default expectation is already set for the Tester.FailNow method")
	}

	if len(m.expectations) > 0 {
		m.mock.t.Fatalf("Some expectations are already set for the Tester.FailNow method")
	}

	m.mock.funcFailNow = f
	return m.mock
}

// FailNow implements minimock.Tester
func (m *TesterMock) FailNow() {
	atomic.AddUint64(&m.beforeFailNowCounter, 1)
	defer atomic.AddUint64(&m.afterFailNowCounter, 1)

	if m.FailNowMock.defaultExpectation != nil {
		atomic.AddUint64(&m.FailNowMock.defaultExpectation.Counter, 1)

		return

	}
	if m.funcFailNow != nil {
		m.funcFailNow()
		return
	}
	m.t.Fatalf("Unexpected call to TesterMock.FailNow.")

}

// FailNowAfterCounter returns a count of finished TesterMock.FailNow invocations
func (m *TesterMock) FailNowAfterCounter() uint64 {
	return atomic.LoadUint64(&m.afterFailNowCounter)
}

// FailNowBeforeCounter returns a count of TesterMock.FailNow invocations
func (m *TesterMock) FailNowBeforeCounter() uint64 {
	return atomic.LoadUint64(&m.beforeFailNowCounter)
}

// MinimockFailNowDone returns true if the count of the FailNow invocations corresponds
// the number of defined expectations
func (m *TesterMock) MinimockFailNowDone() bool {
	for _, e := range m.FailNowMock.expectations {
		if atomic.LoadUint64(&e.Counter) < 1 {
			return false
		}
	}

	// if default expectation was set then invocations count should be greater than zero
	if m.FailNowMock.defaultExpectation != nil && atomic.LoadUint64(&m.afterFailNowCounter) < 1 {
		return false
	}
	// if func was set then invocations count should be greater than zero
	if m.funcFailNow != nil && atomic.LoadUint64(&m.afterFailNowCounter) < 1 {
		return false
	}
	return true
}

// MinimockFailNowInspect logs each unmet expectation
func (m *TesterMock) MinimockFailNowInspect() {
	for _, e := range m.FailNowMock.expectations {
		if atomic.LoadUint64(&e.Counter) < 1 {
			m.t.Error("Expected call to TesterMock.FailNow")
		}
	}

	// if default expectation was set then invocations count should be greater than zero
	if m.FailNowMock.defaultExpectation != nil && atomic.LoadUint64(&m.afterFailNowCounter) < 1 {
		m.t.Error("Expected call to TesterMock.FailNow")
	}
	// if func was set then invocations count should be greater than zero
	if m.funcFailNow != nil && atomic.LoadUint64(&m.afterFailNowCounter) < 1 {
		m.t.Error("Expected call to TesterMock.FailNow")
	}
}

type mTesterMockFatal struct {
	mock               *TesterMock
	defaultExpectation *TesterMockFatalExpectation
	expectations       []*TesterMockFatalExpectation
}

// TesterMockFatalExpectation specifies expectation struct of the Tester.Fatal
type TesterMockFatalExpectation struct {
	mock   *TesterMock
	params *TesterMockFatalParams

	Counter uint64
}

// TesterMockFatalParams contains parameters of the Tester.Fatal
type TesterMockFatalParams struct {
	args []interface{}
}

// Expect sets up expected params for Tester.Fatal
func (m *mTesterMockFatal) Expect(args ...interface{}) *mTesterMockFatal {
	if m.mock.funcFatal != nil {
		m.mock.t.Fatalf("TesterMock.Fatal mock is already set by Set")
	}

	if m.defaultExpectation == nil {
		m.defaultExpectation = &TesterMockFatalExpectation{}
	}

	m.defaultExpectation.params = &TesterMockFatalParams{args}
	for _, e := range m.expectations {
		if minimock.Equal(e.params, m.defaultExpectation.params) {
			m.mock.t.Fatalf("Expectation set by When has same params: %#v", *m.defaultExpectation.params)
		}
	}

	return m
}

// Return sets up results that will be returned by Tester.Fatal
func (m *mTesterMockFatal) Return() *TesterMock {
	if m.mock.funcFatal != nil {
		m.mock.t.Fatalf("TesterMock.Fatal mock is already set by Set")
	}

	if m.defaultExpectation == nil {
		m.defaultExpectation = &TesterMockFatalExpectation{mock: m.mock}
	}

	return m.mock
}

//Set uses given function f to mock the Tester.Fatal method
func (m *mTesterMockFatal) Set(f func(args ...interface{})) *TesterMock {
	if m.defaultExpectation != nil {
		m.mock.t.Fatalf("Default expectation is already set for the Tester.Fatal method")
	}

	if len(m.expectations) > 0 {
		m.mock.t.Fatalf("Some expectations are already set for the Tester.Fatal method")
	}

	m.mock.funcFatal = f
	return m.mock
}

// Fatal implements minimock.Tester
func (m *TesterMock) Fatal(args ...interface{}) {
	atomic.AddUint64(&m.beforeFatalCounter, 1)
	defer atomic.AddUint64(&m.afterFatalCounter, 1)

	for _, e := range m.FatalMock.expectations {
		if minimock.Equal(*e.params, TesterMockFatalParams{args}) {
			atomic.AddUint64(&e.Counter, 1)
			return
		}
	}

	if m.FatalMock.defaultExpectation != nil {
		atomic.AddUint64(&m.FatalMock.defaultExpectation.Counter, 1)
		want := m.FatalMock.defaultExpectation.params
		got := TesterMockFatalParams{args}
		if want != nil && !minimock.Equal(*want, got) {
			m.t.Errorf("TesterMock.Fatal got unexpected parameters, want: %#v, got: %#v\n", *want, got)
		}

		return

	}
	if m.funcFatal != nil {
		m.funcFatal(args...)
		return
	}
	m.t.Fatalf("Unexpected call to TesterMock.Fatal. %v", args)

}

// FatalAfterCounter returns a count of finished TesterMock.Fatal invocations
func (m *TesterMock) FatalAfterCounter() uint64 {
	return atomic.LoadUint64(&m.afterFatalCounter)
}

// FatalBeforeCounter returns a count of TesterMock.Fatal invocations
func (m *TesterMock) FatalBeforeCounter() uint64 {
	return atomic.LoadUint64(&m.beforeFatalCounter)
}

// MinimockFatalDone returns true if the count of the Fatal invocations corresponds
// the number of defined expectations
func (m *TesterMock) MinimockFatalDone() bool {
	for _, e := range m.FatalMock.expectations {
		if atomic.LoadUint64(&e.Counter) < 1 {
			return false
		}
	}

	// if default expectation was set then invocations count should be greater than zero
	if m.FatalMock.defaultExpectation != nil && atomic.LoadUint64(&m.afterFatalCounter) < 1 {
		return false
	}
	// if func was set then invocations count should be greater than zero
	if m.funcFatal != nil && atomic.LoadUint64(&m.afterFatalCounter) < 1 {
		return false
	}
	return true
}

// MinimockFatalInspect logs each unmet expectation
func (m *TesterMock) MinimockFatalInspect() {
	for _, e := range m.FatalMock.expectations {
		if atomic.LoadUint64(&e.Counter) < 1 {
			m.t.Errorf("Expected call to TesterMock.Fatal with params: %#v", *e.params)
		}
	}

	// if default expectation was set then invocations count should be greater than zero
	if m.FatalMock.defaultExpectation != nil && atomic.LoadUint64(&m.afterFatalCounter) < 1 {
		m.t.Errorf("Expected call to TesterMock.Fatal with params: %#v", *m.FatalMock.defaultExpectation.params)
	}
	// if func was set then invocations count should be greater than zero
	if m.funcFatal != nil && atomic.LoadUint64(&m.afterFatalCounter) < 1 {
		m.t.Error("Expected call to TesterMock.Fatal")
	}
}

type mTesterMockFatalf struct {
	mock               *TesterMock
	defaultExpectation *TesterMockFatalfExpectation
	expectations       []*TesterMockFatalfExpectation
}

// TesterMockFatalfExpectation specifies expectation struct of the Tester.Fatalf
type TesterMockFatalfExpectation struct {
	mock   *TesterMock
	params *TesterMockFatalfParams

	Counter uint64
}

// TesterMockFatalfParams contains parameters of the Tester.Fatalf
type TesterMockFatalfParams struct {
	format string
	args   []interface{}
}

// Expect sets up expected params for Tester.Fatalf
func (m *mTesterMockFatalf) Expect(format string, args ...interface{}) *mTesterMockFatalf {
	if m.mock.funcFatalf != nil {
		m.mock.t.Fatalf("TesterMock.Fatalf mock is already set by Set")
	}

	if m.defaultExpectation == nil {
		m.defaultExpectation = &TesterMockFatalfExpectation{}
	}

	m.defaultExpectation.params = &TesterMockFatalfParams{format, args}
	for _, e := range m.expectations {
		if minimock.Equal(e.params, m.defaultExpectation.params) {
			m.mock.t.Fatalf("Expectation set by When has same params: %#v", *m.defaultExpectation.params)
		}
	}

	return m
}

// Return sets up results that will be returned by Tester.Fatalf
func (m *mTesterMockFatalf) Return() *TesterMock {
	if m.mock.funcFatalf != nil {
		m.mock.t.Fatalf("TesterMock.Fatalf mock is already set by Set")
	}

	if m.defaultExpectation == nil {
		m.defaultExpectation = &TesterMockFatalfExpectation{mock: m.mock}
	}

	return m.mock
}

//Set uses given function f to mock the Tester.Fatalf method
func (m *mTesterMockFatalf) Set(f func(format string, args ...interface{})) *TesterMock {
	if m.defaultExpectation != nil {
		m.mock.t.Fatalf("Default expectation is already set for the Tester.Fatalf method")
	}

	if len(m.expectations) > 0 {
		m.mock.t.Fatalf("Some expectations are already set for the Tester.Fatalf method")
	}

	m.mock.funcFatalf = f
	return m.mock
}

// Fatalf implements minimock.Tester
func (m *TesterMock) Fatalf(format string, args ...interface{}) {
	atomic.AddUint64(&m.beforeFatalfCounter, 1)
	defer atomic.AddUint64(&m.afterFatalfCounter, 1)

	for _, e := range m.FatalfMock.expectations {
		if minimock.Equal(*e.params, TesterMockFatalfParams{format, args}) {
			atomic.AddUint64(&e.Counter, 1)
			return
		}
	}

	if m.FatalfMock.defaultExpectation != nil {
		atomic.AddUint64(&m.FatalfMock.defaultExpectation.Counter, 1)
		want := m.FatalfMock.defaultExpectation.params
		got := TesterMockFatalfParams{format, args}
		if want != nil && !minimock.Equal(*want, got) {
			m.t.Errorf("TesterMock.Fatalf got unexpected parameters, want: %#v, got: %#v\n", *want, got)
		}

		return

	}
	if m.funcFatalf != nil {
		m.funcFatalf(format, args...)
		return
	}
	m.t.Fatalf("Unexpected call to TesterMock.Fatalf. %v %v", format, args)

}

// FatalfAfterCounter returns a count of finished TesterMock.Fatalf invocations
func (m *TesterMock) FatalfAfterCounter() uint64 {
	return atomic.LoadUint64(&m.afterFatalfCounter)
}

// FatalfBeforeCounter returns a count of TesterMock.Fatalf invocations
func (m *TesterMock) FatalfBeforeCounter() uint64 {
	return atomic.LoadUint64(&m.beforeFatalfCounter)
}

// MinimockFatalfDone returns true if the count of the Fatalf invocations corresponds
// the number of defined expectations
func (m *TesterMock) MinimockFatalfDone() bool {
	for _, e := range m.FatalfMock.expectations {
		if atomic.LoadUint64(&e.Counter) < 1 {
			return false
		}
	}

	// if default expectation was set then invocations count should be greater than zero
	if m.FatalfMock.defaultExpectation != nil && atomic.LoadUint64(&m.afterFatalfCounter) < 1 {
		return false
	}
	// if func was set then invocations count should be greater than zero
	if m.funcFatalf != nil && atomic.LoadUint64(&m.afterFatalfCounter) < 1 {
		return false
	}
	return true
}

// MinimockFatalfInspect logs each unmet expectation
func (m *TesterMock) MinimockFatalfInspect() {
	for _, e := range m.FatalfMock.expectations {
		if atomic.LoadUint64(&e.Counter) < 1 {
			m.t.Errorf("Expected call to TesterMock.Fatalf with params: %#v", *e.params)
		}
	}

	// if default expectation was set then invocations count should be greater than zero
	if m.FatalfMock.defaultExpectation != nil && atomic.LoadUint64(&m.afterFatalfCounter) < 1 {
		m.t.Errorf("Expected call to TesterMock.Fatalf with params: %#v", *m.FatalfMock.defaultExpectation.params)
	}
	// if func was set then invocations count should be greater than zero
	if m.funcFatalf != nil && atomic.LoadUint64(&m.afterFatalfCounter) < 1 {
		m.t.Error("Expected call to TesterMock.Fatalf")
	}
}

// MinimockFinish checks that all mocked methods have been called the expected number of times
func (m *TesterMock) MinimockFinish() {
	if !m.minimockDone() {
		m.MinimockErrorInspect()

		m.MinimockErrorfInspect()

		m.MinimockFailNowInspect()

		m.MinimockFatalInspect()

		m.MinimockFatalfInspect()
		m.t.FailNow()
	}
}

// MinimockWait waits for all mocked methods to be called the expected number of times
func (m *TesterMock) MinimockWait(timeout time.Duration) {
	timeoutCh := time.After(timeout)
	for {
		if m.minimockDone() {
			return
		}
		select {
		case <-timeoutCh:
			m.MinimockFinish()
			return
		case <-time.After(10 * time.Millisecond):
		}
	}
}

func (m *TesterMock) minimockDone() bool {
	done := true
	return done &&
		m.MinimockErrorDone() &&
		m.MinimockErrorfDone() &&
		m.MinimockFailNowDone() &&
		m.MinimockFatalDone() &&
		m.MinimockFatalfDone()
}
