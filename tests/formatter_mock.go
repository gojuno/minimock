package tests

/*
DO NOT EDIT!
This code was generated automatically using github.com/gojuno/minimock v1.9
The original interface "Formatter" can be found in github.com/gojuno/minimock/tests
*/
import (
	"sync/atomic"
	"time"

	"github.com/gojuno/minimock"
	testify_assert "github.com/stretchr/testify/assert"
)

//FormatterMock implements github.com/gojuno/minimock/tests.Formatter
type FormatterMock struct {
	t minimock.Tester

	FormatFunc       func(p string, p1 ...interface{}) (r string)
	FormatCounter    uint64
	FormatPreCounter uint64
	FormatMock       mFormatterMockFormat
}

//NewFormatterMock returns a mock for github.com/gojuno/minimock/tests.Formatter
func NewFormatterMock(t minimock.Tester) *FormatterMock {
	m := &FormatterMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.FormatMock = mFormatterMockFormat{mock: m}

	return m
}

type mFormatterMockFormat struct {
	mock             *FormatterMock
	mockExpectations *FormatterMockFormatParams
}

//FormatterMockFormatParams represents input parameters of the Formatter.Format
type FormatterMockFormatParams struct {
	p  string
	p1 []interface{}
}

//Expect sets up expected params for the Formatter.Format
func (m *mFormatterMockFormat) Expect(p string, p1 ...interface{}) *mFormatterMockFormat {
	m.mockExpectations = &FormatterMockFormatParams{p, p1}
	return m
}

//Return sets up a mock for Formatter.Format to return Return's arguments
func (m *mFormatterMockFormat) Return(r string) *FormatterMock {
	m.mock.FormatFunc = func(p string, p1 ...interface{}) string {
		return r
	}
	return m.mock
}

//Set uses given function f as a mock of Formatter.Format method
func (m *mFormatterMockFormat) Set(f func(p string, p1 ...interface{}) (r string)) *FormatterMock {
	m.mock.FormatFunc = f
	return m.mock
}

//Format implements github.com/gojuno/minimock/tests.Formatter interface
func (m *FormatterMock) Format(p string, p1 ...interface{}) (r string) {
	atomic.AddUint64(&m.FormatPreCounter, 1)
	defer atomic.AddUint64(&m.FormatCounter, 1)

	if m.FormatMock.mockExpectations != nil {
		testify_assert.Equal(m.t, *m.FormatMock.mockExpectations, FormatterMockFormatParams{p, p1},
			"Formatter.Format got unexpected parameters")

		if m.FormatFunc == nil {

			m.t.Fatal("No results are set for the FormatterMock.Format")

			return
		}
	}

	if m.FormatFunc == nil {
		m.t.Fatal("Unexpected call to FormatterMock.Format")
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

//ValidateCallCounters checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *FormatterMock) ValidateCallCounters() {

	if m.FormatFunc != nil && atomic.LoadUint64(&m.FormatCounter) == 0 {
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

	if m.FormatFunc != nil && atomic.LoadUint64(&m.FormatCounter) == 0 {
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
		ok = ok && (m.FormatFunc == nil || atomic.LoadUint64(&m.FormatCounter) > 0)

		if ok {
			return
		}

		select {
		case <-timeoutCh:

			if m.FormatFunc != nil && atomic.LoadUint64(&m.FormatCounter) == 0 {
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

	if m.FormatFunc != nil && atomic.LoadUint64(&m.FormatCounter) == 0 {
		return false
	}

	return true
}
