package tests

// DO NOT EDIT!
// The code below was generated with http://github.com/deff7/minimock (dev)

//go:generate minimock -i github.com/deff7/minimock/tests.Formatter -o ./tests/formatter_mock.go

import (
	"sync/atomic"
	"time"

	minimock "github.com/deff7/minimock/pkg"
)

// FormatterMock implements Formatter
type FormatterMock struct {
	t minimock.Tester

	funcFormat          func(s1 string, p1 ...interface{}) (s2 string)
	afterFormatCounter  uint64
	beforeFormatCounter uint64
	FormatMock          mFormatterMockFormat
}

// NewFormatterMock returns a mock for Formatter
func NewFormatterMock(t minimock.Tester) *FormatterMock {
	m := &FormatterMock{t: t}
	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}
	m.FormatMock = mFormatterMockFormat{mock: m}

	return m
}

type mFormatterMockFormat struct {
	mock               *FormatterMock
	defaultExpectation *FormatterMockFormatExpectation
	expectations       []*FormatterMockFormatExpectation
}

// FormatterMockFormatExpectation specifies expectation struct of the Formatter.Format
type FormatterMockFormatExpectation struct {
	mock    *FormatterMock
	params  *FormatterMockFormatParams
	results *FormatterMockFormatResults
	Counter uint64
}

// FormatterMockFormatParams contains parameters of the Formatter.Format
type FormatterMockFormatParams struct {
	s1 string
	p1 []interface{}
}

// FormatterMockFormatResults contains results of the Formatter.Format
type FormatterMockFormatResults struct {
	s2 string
}

// Expect sets up expected params for Formatter.Format
func (m *mFormatterMockFormat) Expect(s1 string, p1 ...interface{}) *mFormatterMockFormat {
	if m.mock.funcFormat != nil {
		m.mock.t.Fatalf("FormatterMock.Format mock is already set by Set")
	}

	if m.defaultExpectation == nil {
		m.defaultExpectation = &FormatterMockFormatExpectation{}
	}

	m.defaultExpectation.params = &FormatterMockFormatParams{s1, p1}
	for _, e := range m.expectations {
		if minimock.Equal(e.params, m.defaultExpectation.params) {
			m.mock.t.Fatalf("Expectation set by When has same params: %#v", *m.defaultExpectation.params)
		}
	}

	return m
}

// Return sets up results that will be returned by Formatter.Format
func (m *mFormatterMockFormat) Return(s2 string) *FormatterMock {
	if m.mock.funcFormat != nil {
		m.mock.t.Fatalf("FormatterMock.Format mock is already set by Set")
	}

	if m.defaultExpectation == nil {
		m.defaultExpectation = &FormatterMockFormatExpectation{mock: m.mock}
	}
	m.defaultExpectation.results = &FormatterMockFormatResults{s2}
	return m.mock
}

//Set uses given function f to mock the Formatter.Format method
func (m *mFormatterMockFormat) Set(f func(s1 string, p1 ...interface{}) (s2 string)) *FormatterMock {
	if m.defaultExpectation != nil {
		m.mock.t.Fatalf("Default expectation is already set for the Formatter.Format method")
	}

	if len(m.expectations) > 0 {
		m.mock.t.Fatalf("Some expectations are already set for the Formatter.Format method")
	}

	m.mock.funcFormat = f
	return m.mock
}

// When sets expectation for the Formatter.Format which will trigger the result defined by the following
// Then helper
func (m *mFormatterMockFormat) When(s1 string, p1 ...interface{}) *FormatterMockFormatExpectation {
	if m.mock.funcFormat != nil {
		m.mock.t.Fatalf("FormatterMock.Format mock is already set by Set")
	}

	expectation := &FormatterMockFormatExpectation{
		mock:   m.mock,
		params: &FormatterMockFormatParams{s1, p1},
	}
	m.expectations = append(m.expectations, expectation)
	return expectation
}

// Then sets up Formatter.Format return parameters for the expectation previously defined by the When method
func (e *FormatterMockFormatExpectation) Then(s2 string) *FormatterMock {
	e.results = &FormatterMockFormatResults{s2}
	return e.mock
}

// Format implements Formatter
func (m *FormatterMock) Format(s1 string, p1 ...interface{}) (s2 string) {
	atomic.AddUint64(&m.beforeFormatCounter, 1)
	defer atomic.AddUint64(&m.afterFormatCounter, 1)

	for _, e := range m.FormatMock.expectations {
		if minimock.Equal(*e.params, FormatterMockFormatParams{s1, p1}) {
			atomic.AddUint64(&e.Counter, 1)
			return e.results.s2
		}
	}

	if m.FormatMock.defaultExpectation != nil {
		atomic.AddUint64(&m.FormatMock.defaultExpectation.Counter, 1)
		want := m.FormatMock.defaultExpectation.params
		got := FormatterMockFormatParams{s1, p1}
		if want != nil && !minimock.Equal(*want, got) {
			m.t.Errorf("FormatterMock.Format got unexpected parameters, want: %#v, got: %#v\n", *want, got)
		}

		results := m.FormatMock.defaultExpectation.results
		if results == nil {
			m.t.Fatal("No results are set for the FormatterMock.Format")
		}
		return (*results).s2
	}
	if m.funcFormat != nil {
		return m.funcFormat(s1, p1...)
	}
	m.t.Fatalf("Unexpected call to FormatterMock.Format. %v %v", s1, p1)
	return
}

// FormatAfterCounter returns a count of finished FormatterMock.Format invocations
func (m *FormatterMock) FormatAfterCounter() uint64 {
	return atomic.LoadUint64(&m.afterFormatCounter)
}

// FormatBeforeCounter returns a count of FormatterMock.Format invocations
func (m *FormatterMock) FormatBeforeCounter() uint64 {
	return atomic.LoadUint64(&m.beforeFormatCounter)
}

// MinimockFormatDone returns true if the count of the Format invocations corresponds
// the number of defined expectations
func (m *FormatterMock) MinimockFormatDone() bool {
	for _, e := range m.FormatMock.expectations {
		if atomic.LoadUint64(&e.Counter) < 1 {
			return false
		}
	}

	// if default expectation was set then invocations count should be greater than zero
	if m.FormatMock.defaultExpectation != nil && atomic.LoadUint64(&m.afterFormatCounter) < 1 {
		return false
	}
	// if func was set then invocations count should be greater than zero
	if m.funcFormat != nil && atomic.LoadUint64(&m.afterFormatCounter) < 1 {
		return false
	}
	return true
}

// MinimockFormatInspect logs each unmet expectation
func (m *FormatterMock) MinimockFormatInspect() {
	for _, e := range m.FormatMock.expectations {
		if atomic.LoadUint64(&e.Counter) < 1 {
			m.t.Errorf("Expected call to FormatterMock.Format with params: %#v", *e.params)
		}
	}

	// if default expectation was set then invocations count should be greater than zero
	if m.FormatMock.defaultExpectation != nil && atomic.LoadUint64(&m.afterFormatCounter) < 1 {
		m.t.Errorf("Expected call to FormatterMock.Format with params: %#v", *m.FormatMock.defaultExpectation.params)
	}
	// if func was set then invocations count should be greater than zero
	if m.funcFormat != nil && atomic.LoadUint64(&m.afterFormatCounter) < 1 {
		m.t.Error("Expected call to FormatterMock.Format")
	}
}

// MinimockFinish checks that all mocked methods have been called the expected number of times
func (m *FormatterMock) MinimockFinish() {
	if !m.minimockDone() {
		m.MinimockFormatInspect()
		m.t.FailNow()
	}
}

// MinimockWait waits for all mocked methods to be called the expected number of times
func (m *FormatterMock) MinimockWait(timeout time.Duration) {
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

func (m *FormatterMock) minimockDone() bool {
	done := true
	return done &&
		m.MinimockFormatDone()
}
