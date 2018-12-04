package tests

/*
DO NOT EDIT!
This code was generated automatically using github.com/gojuno/minimock v1.9
The original interface "Formatter" can be found in github.com/kirylandruski/minimock/tests
*/
import (
	"sync/atomic"
	"time"

	"github.com/gojuno/minimock"
	testify_assert "github.com/stretchr/testify/assert"
)

//FormatterMock implements github.com/kirylandruski/minimock/tests.Formatter
type FormatterMock struct {
	t minimock.Tester

	FormatFunc       func(p string, p1 ...interface{}) (r string)
	FormatCounter    uint64
	FormatPreCounter uint64
	FormatMock       mFormatterMockFormat
}

//NewFormatterMock returns a mock for github.com/kirylandruski/minimock/tests.Formatter
func NewFormatterMock(t minimock.Tester) *FormatterMock {
	m := &FormatterMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.FormatMock = mFormatterMockFormat{mock: m}

	return m
}

type mFormatterMockFormat struct {
	mock              *FormatterMock
	mainExpectation   *FormatterMockFormatExpectation
	expectationSeries []*FormatterMockFormatExpectation
}

type FormatterMockFormatExpectation struct {
	input  *FormatterMockFormatInput
	result *FormatterMockFormatResult
}

type FormatterMockFormatInput struct {
	p  string
	p1 []interface{}
}

type FormatterMockFormatResult struct {
	r string
}

//Expect specifies that invocation of Formatter.Format is expected from 1 to Infinity times
func (m *mFormatterMockFormat) Expect(p string, p1 ...interface{}) *mFormatterMockFormat {
	m.mock.FormatFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &FormatterMockFormatExpectation{}
	}
	m.mainExpectation.input = &FormatterMockFormatInput{p, p1}
	return m
}

//Return specifies results of invocation of Formatter.Format
func (m *mFormatterMockFormat) Return(r string) *FormatterMock {
	m.mock.FormatFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &FormatterMockFormatExpectation{}
	}
	m.mainExpectation.result = &FormatterMockFormatResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of Formatter.Format is expected once
func (m *mFormatterMockFormat) ExpectOnce(p string, p1 ...interface{}) *FormatterMockFormatExpectation {
	m.mock.FormatFunc = nil
	m.mainExpectation = nil

	expectation := &FormatterMockFormatExpectation{}
	expectation.input = &FormatterMockFormatInput{p, p1}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *FormatterMockFormatExpectation) Return(r string) {
	e.result = &FormatterMockFormatResult{r}
}

//Set uses given function f as a mock of Formatter.Format method
func (m *mFormatterMockFormat) Set(f func(p string, p1 ...interface{}) (r string)) *FormatterMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.FormatFunc = f
	return m.mock
}

//Format implements github.com/kirylandruski/minimock/tests.Formatter interface
func (m *FormatterMock) Format(p string, p1 ...interface{}) (r string) {
	counter := atomic.AddUint64(&m.FormatPreCounter, 1)
	defer atomic.AddUint64(&m.FormatCounter, 1)

	if len(m.FormatMock.expectationSeries) > 0 {
		if counter > uint64(len(m.FormatMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to FormatterMock.Format. %v %v", p, p1)
			return
		}

		input := m.FormatMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, FormatterMockFormatInput{p, p1}, "Formatter.Format got unexpected parameters")

		result := m.FormatMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the FormatterMock.Format")
			return
		}

		r = result.r

		return
	}

	if m.FormatMock.mainExpectation != nil {

		input := m.FormatMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, FormatterMockFormatInput{p, p1}, "Formatter.Format got unexpected parameters")
		}

		result := m.FormatMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the FormatterMock.Format")
		}

		r = result.r

		return
	}

	if m.FormatFunc == nil {
		m.t.Fatalf("Unexpected call to FormatterMock.Format. %v %v", p, p1)
		return
	}

	return m.FormatFunc(p, p1...)
}

//FormatMinimockCounter returns a count of FormatterMock.FormatFunc invocations
func (m *FormatterMock) FormatMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.FormatCounter)
}

//FormatMinimockPreCounter returns the value of FormatterMock.Format invocations
func (m *FormatterMock) FormatMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.FormatPreCounter)
}

//FormatFinished returns true if mock invocations count is ok
func (m *FormatterMock) FormatFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.FormatMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.FormatCounter) == uint64(len(m.FormatMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.FormatMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.FormatCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.FormatFunc != nil {
		return atomic.LoadUint64(&m.FormatCounter) > 0
	}

	return true
}

//ValidateCallCounters checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *FormatterMock) ValidateCallCounters() {

	if !m.FormatFinished() {
		m.t.Fatal("Expected call to FormatterMock.Format")
	}

}

//CheckMocksCalled checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *FormatterMock) CheckMocksCalled() {
	m.Finish()
}

//Finish checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish or use Finish method of minimock.Controller
func (m *FormatterMock) Finish() {
	m.MinimockFinish()
}

//MinimockFinish checks that all mocked methods of the interface have been called at least once
func (m *FormatterMock) MinimockFinish() {

	if !m.FormatFinished() {
		m.t.Fatal("Expected call to FormatterMock.Format")
	}

}

//Wait waits for all mocked methods to be called at least once
//Deprecated: please use MinimockWait or use Wait method of minimock.Controller
func (m *FormatterMock) Wait(timeout time.Duration) {
	m.MinimockWait(timeout)
}

//MinimockWait waits for all mocked methods to be called at least once
//this method is called by minimock.Controller
func (m *FormatterMock) MinimockWait(timeout time.Duration) {
	timeoutCh := time.After(timeout)
	for {
		ok := true
		ok = ok && m.FormatFinished()

		if ok {
			return
		}

		select {
		case <-timeoutCh:

			if !m.FormatFinished() {
				m.t.Error("Expected call to FormatterMock.Format")
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
func (m *FormatterMock) AllMocksCalled() bool {

	if !m.FormatFinished() {
		return false
	}

	return true
}
