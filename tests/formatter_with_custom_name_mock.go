// Code generated by http://github.com/gojuno/minimock (dev). DO NOT EDIT.

package tests

//go:generate minimock -i github.com/gojuno/minimock/v3/tests.Formatter -o formatter_with_custom_name_mock.go -n CustomFormatterNameMock -p tests

import (
	"sync"
	mm_atomic "sync/atomic"
	mm_time "time"

	"github.com/gojuno/minimock/v3"
)

// CustomFormatterNameMock implements Formatter
type CustomFormatterNameMock struct {
	t          minimock.Tester
	finishOnce sync.Once

	funcFormat          func(s1 string, p1 ...interface{}) (s2 string)
	funcFormatOrigin    string
	inspectFuncFormat   func(s1 string, p1 ...interface{})
	afterFormatCounter  uint64
	beforeFormatCounter uint64
	FormatMock          mCustomFormatterNameMockFormat
}

// NewCustomFormatterNameMock returns a mock for Formatter
func NewCustomFormatterNameMock(t minimock.Tester) *CustomFormatterNameMock {
	m := &CustomFormatterNameMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.FormatMock = mCustomFormatterNameMockFormat{mock: m}
	m.FormatMock.callArgs = []*CustomFormatterNameMockFormatParams{}

	t.Cleanup(m.MinimockFinish)

	return m
}

type mCustomFormatterNameMockFormat struct {
	optional           bool
	mock               *CustomFormatterNameMock
	defaultExpectation *CustomFormatterNameMockFormatExpectation
	expectations       []*CustomFormatterNameMockFormatExpectation

	callArgs []*CustomFormatterNameMockFormatParams
	mutex    sync.RWMutex

	expectedInvocations       uint64
	expectedInvocationsOrigin string
}

// CustomFormatterNameMockFormatExpectation specifies expectation struct of the Formatter.Format
type CustomFormatterNameMockFormatExpectation struct {
	mock               *CustomFormatterNameMock
	params             *CustomFormatterNameMockFormatParams
	paramPtrs          *CustomFormatterNameMockFormatParamPtrs
	expectationOrigins CustomFormatterNameMockFormatExpectationOrigins
	results            *CustomFormatterNameMockFormatResults
	returnOrigin       string
	Counter            uint64
}

// CustomFormatterNameMockFormatParams contains parameters of the Formatter.Format
type CustomFormatterNameMockFormatParams struct {
	s1 string
	p1 []interface{}
}

// CustomFormatterNameMockFormatParamPtrs contains pointers to parameters of the Formatter.Format
type CustomFormatterNameMockFormatParamPtrs struct {
	s1 *string
	p1 *[]interface{}
}

// CustomFormatterNameMockFormatResults contains results of the Formatter.Format
type CustomFormatterNameMockFormatResults struct {
	s2 string
}

// CustomFormatterNameMockFormatOrigins contains origins of expectations of the Formatter.Format
type CustomFormatterNameMockFormatExpectationOrigins struct {
	origin   string
	originS1 string
	originP1 string
}

// Marks this method to be optional. The default behavior of any method with Return() is '1 or more', meaning
// the test will fail minimock's automatic final call check if the mocked method was not called at least once.
// Optional() makes method check to work in '0 or more' mode.
// It is NOT RECOMMENDED to use this option unless you really need it, as default behaviour helps to
// catch the problems when the expected method call is totally skipped during test run.
func (mmFormat *mCustomFormatterNameMockFormat) Optional() *mCustomFormatterNameMockFormat {
	mmFormat.optional = true
	return mmFormat
}

// Expect sets up expected params for Formatter.Format
func (mmFormat *mCustomFormatterNameMockFormat) Expect(s1 string, p1 ...interface{}) *mCustomFormatterNameMockFormat {
	if mmFormat.mock.funcFormat != nil {
		mmFormat.mock.t.Fatalf("CustomFormatterNameMock.Format mock is already set by Set")
	}

	if mmFormat.defaultExpectation == nil {
		mmFormat.defaultExpectation = &CustomFormatterNameMockFormatExpectation{}
	}

	if mmFormat.defaultExpectation.paramPtrs != nil {
		mmFormat.mock.t.Fatalf("CustomFormatterNameMock.Format mock is already set by ExpectParams functions")
	}

	mmFormat.defaultExpectation.params = &CustomFormatterNameMockFormatParams{s1, p1}
	mmFormat.defaultExpectation.expectationOrigins.origin = minimock.CallerInfo(1)
	for _, e := range mmFormat.expectations {
		if minimock.Equal(e.params, mmFormat.defaultExpectation.params) {
			mmFormat.mock.t.Fatalf("Expectation set by When has same params: %#v", *mmFormat.defaultExpectation.params)
		}
	}

	return mmFormat
}

// ExpectS1Param1 sets up expected param s1 for Formatter.Format
func (mmFormat *mCustomFormatterNameMockFormat) ExpectS1Param1(s1 string) *mCustomFormatterNameMockFormat {
	if mmFormat.mock.funcFormat != nil {
		mmFormat.mock.t.Fatalf("CustomFormatterNameMock.Format mock is already set by Set")
	}

	if mmFormat.defaultExpectation == nil {
		mmFormat.defaultExpectation = &CustomFormatterNameMockFormatExpectation{}
	}

	if mmFormat.defaultExpectation.params != nil {
		mmFormat.mock.t.Fatalf("CustomFormatterNameMock.Format mock is already set by Expect")
	}

	if mmFormat.defaultExpectation.paramPtrs == nil {
		mmFormat.defaultExpectation.paramPtrs = &CustomFormatterNameMockFormatParamPtrs{}
	}
	mmFormat.defaultExpectation.paramPtrs.s1 = &s1
	mmFormat.defaultExpectation.expectationOrigins.originS1 = minimock.CallerInfo(1)

	return mmFormat
}

// ExpectP1Param2 sets up expected param p1 for Formatter.Format
func (mmFormat *mCustomFormatterNameMockFormat) ExpectP1Param2(p1 ...interface{}) *mCustomFormatterNameMockFormat {
	if mmFormat.mock.funcFormat != nil {
		mmFormat.mock.t.Fatalf("CustomFormatterNameMock.Format mock is already set by Set")
	}

	if mmFormat.defaultExpectation == nil {
		mmFormat.defaultExpectation = &CustomFormatterNameMockFormatExpectation{}
	}

	if mmFormat.defaultExpectation.params != nil {
		mmFormat.mock.t.Fatalf("CustomFormatterNameMock.Format mock is already set by Expect")
	}

	if mmFormat.defaultExpectation.paramPtrs == nil {
		mmFormat.defaultExpectation.paramPtrs = &CustomFormatterNameMockFormatParamPtrs{}
	}
	mmFormat.defaultExpectation.paramPtrs.p1 = &p1
	mmFormat.defaultExpectation.expectationOrigins.originP1 = minimock.CallerInfo(1)

	return mmFormat
}

// Inspect accepts an inspector function that has same arguments as the Formatter.Format
func (mmFormat *mCustomFormatterNameMockFormat) Inspect(f func(s1 string, p1 ...interface{})) *mCustomFormatterNameMockFormat {
	if mmFormat.mock.inspectFuncFormat != nil {
		mmFormat.mock.t.Fatalf("Inspect function is already set for CustomFormatterNameMock.Format")
	}

	mmFormat.mock.inspectFuncFormat = f

	return mmFormat
}

// Return sets up results that will be returned by Formatter.Format
func (mmFormat *mCustomFormatterNameMockFormat) Return(s2 string) *CustomFormatterNameMock {
	if mmFormat.mock.funcFormat != nil {
		mmFormat.mock.t.Fatalf("CustomFormatterNameMock.Format mock is already set by Set")
	}

	if mmFormat.defaultExpectation == nil {
		mmFormat.defaultExpectation = &CustomFormatterNameMockFormatExpectation{mock: mmFormat.mock}
	}
	mmFormat.defaultExpectation.results = &CustomFormatterNameMockFormatResults{s2}
	mmFormat.defaultExpectation.returnOrigin = minimock.CallerInfo(1)
	return mmFormat.mock
}

// Set uses given function f to mock the Formatter.Format method
func (mmFormat *mCustomFormatterNameMockFormat) Set(f func(s1 string, p1 ...interface{}) (s2 string)) *CustomFormatterNameMock {
	if mmFormat.defaultExpectation != nil {
		mmFormat.mock.t.Fatalf("Default expectation is already set for the Formatter.Format method")
	}

	if len(mmFormat.expectations) > 0 {
		mmFormat.mock.t.Fatalf("Some expectations are already set for the Formatter.Format method")
	}

	mmFormat.mock.funcFormat = f
	mmFormat.mock.funcFormatOrigin = minimock.CallerInfo(1)
	return mmFormat.mock
}

// When sets expectation for the Formatter.Format which will trigger the result defined by the following
// Then helper
func (mmFormat *mCustomFormatterNameMockFormat) When(s1 string, p1 ...interface{}) *CustomFormatterNameMockFormatExpectation {
	if mmFormat.mock.funcFormat != nil {
		mmFormat.mock.t.Fatalf("CustomFormatterNameMock.Format mock is already set by Set")
	}

	expectation := &CustomFormatterNameMockFormatExpectation{
		mock:               mmFormat.mock,
		params:             &CustomFormatterNameMockFormatParams{s1, p1},
		expectationOrigins: CustomFormatterNameMockFormatExpectationOrigins{origin: minimock.CallerInfo(1)},
	}
	mmFormat.expectations = append(mmFormat.expectations, expectation)
	return expectation
}

// Then sets up Formatter.Format return parameters for the expectation previously defined by the When method
func (e *CustomFormatterNameMockFormatExpectation) Then(s2 string) *CustomFormatterNameMock {
	e.results = &CustomFormatterNameMockFormatResults{s2}
	return e.mock
}

// Times sets number of times Formatter.Format should be invoked
func (mmFormat *mCustomFormatterNameMockFormat) Times(n uint64) *mCustomFormatterNameMockFormat {
	if n == 0 {
		mmFormat.mock.t.Fatalf("Times of CustomFormatterNameMock.Format mock can not be zero")
	}
	mm_atomic.StoreUint64(&mmFormat.expectedInvocations, n)
	mmFormat.expectedInvocationsOrigin = minimock.CallerInfo(1)
	return mmFormat
}

func (mmFormat *mCustomFormatterNameMockFormat) invocationsDone() bool {
	if len(mmFormat.expectations) == 0 && mmFormat.defaultExpectation == nil && mmFormat.mock.funcFormat == nil {
		return true
	}

	totalInvocations := mm_atomic.LoadUint64(&mmFormat.mock.afterFormatCounter)
	expectedInvocations := mm_atomic.LoadUint64(&mmFormat.expectedInvocations)

	return totalInvocations > 0 && (expectedInvocations == 0 || expectedInvocations == totalInvocations)
}

// Format implements Formatter
func (mmFormat *CustomFormatterNameMock) Format(s1 string, p1 ...interface{}) (s2 string) {
	mm_atomic.AddUint64(&mmFormat.beforeFormatCounter, 1)
	defer mm_atomic.AddUint64(&mmFormat.afterFormatCounter, 1)

	mmFormat.t.Helper()

	if mmFormat.inspectFuncFormat != nil {
		mmFormat.inspectFuncFormat(s1, p1...)
	}

	mm_params := CustomFormatterNameMockFormatParams{s1, p1}

	// Record call args
	mmFormat.FormatMock.mutex.Lock()
	mmFormat.FormatMock.callArgs = append(mmFormat.FormatMock.callArgs, &mm_params)
	mmFormat.FormatMock.mutex.Unlock()

	for _, e := range mmFormat.FormatMock.expectations {
		if minimock.Equal(*e.params, mm_params) {
			mm_atomic.AddUint64(&e.Counter, 1)
			return e.results.s2
		}
	}

	if mmFormat.FormatMock.defaultExpectation != nil {
		mm_atomic.AddUint64(&mmFormat.FormatMock.defaultExpectation.Counter, 1)
		mm_want := mmFormat.FormatMock.defaultExpectation.params
		mm_want_ptrs := mmFormat.FormatMock.defaultExpectation.paramPtrs

		mm_got := CustomFormatterNameMockFormatParams{s1, p1}

		if mm_want_ptrs != nil {

			if mm_want_ptrs.s1 != nil && !minimock.Equal(*mm_want_ptrs.s1, mm_got.s1) {
				mmFormat.t.Errorf("CustomFormatterNameMock.Format got unexpected parameter s1, expected at\n%s:\nwant: %#v\n got: %#v%s\n",
					mmFormat.FormatMock.defaultExpectation.expectationOrigins.originS1, *mm_want_ptrs.s1, mm_got.s1, minimock.Diff(*mm_want_ptrs.s1, mm_got.s1))
			}

			if mm_want_ptrs.p1 != nil && !minimock.Equal(*mm_want_ptrs.p1, mm_got.p1) {
				mmFormat.t.Errorf("CustomFormatterNameMock.Format got unexpected parameter p1, expected at\n%s:\nwant: %#v\n got: %#v%s\n",
					mmFormat.FormatMock.defaultExpectation.expectationOrigins.originP1, *mm_want_ptrs.p1, mm_got.p1, minimock.Diff(*mm_want_ptrs.p1, mm_got.p1))
			}

		} else if mm_want != nil && !minimock.Equal(*mm_want, mm_got) {
			mmFormat.t.Errorf("CustomFormatterNameMock.Format got unexpected parameters, expected at\n%s:\nwant: %#v\n got: %#v%s\n",
				mmFormat.FormatMock.defaultExpectation.expectationOrigins.origin, *mm_want, mm_got, minimock.Diff(*mm_want, mm_got))
		}

		mm_results := mmFormat.FormatMock.defaultExpectation.results
		if mm_results == nil {
			mmFormat.t.Fatal("No results are set for the CustomFormatterNameMock.Format")
		}
		return (*mm_results).s2
	}
	if mmFormat.funcFormat != nil {
		return mmFormat.funcFormat(s1, p1...)
	}
	mmFormat.t.Fatalf("Unexpected call to CustomFormatterNameMock.Format. %v %v", s1, p1)
	return
}

// FormatAfterCounter returns a count of finished CustomFormatterNameMock.Format invocations
func (mmFormat *CustomFormatterNameMock) FormatAfterCounter() uint64 {
	return mm_atomic.LoadUint64(&mmFormat.afterFormatCounter)
}

// FormatBeforeCounter returns a count of CustomFormatterNameMock.Format invocations
func (mmFormat *CustomFormatterNameMock) FormatBeforeCounter() uint64 {
	return mm_atomic.LoadUint64(&mmFormat.beforeFormatCounter)
}

// Calls returns a list of arguments used in each call to CustomFormatterNameMock.Format.
// The list is in the same order as the calls were made (i.e. recent calls have a higher index)
func (mmFormat *mCustomFormatterNameMockFormat) Calls() []*CustomFormatterNameMockFormatParams {
	mmFormat.mutex.RLock()

	argCopy := make([]*CustomFormatterNameMockFormatParams, len(mmFormat.callArgs))
	copy(argCopy, mmFormat.callArgs)

	mmFormat.mutex.RUnlock()

	return argCopy
}

// MinimockFormatDone returns true if the count of the Format invocations corresponds
// the number of defined expectations
func (m *CustomFormatterNameMock) MinimockFormatDone() bool {
	if m.FormatMock.optional {
		// Optional methods provide '0 or more' call count restriction.
		return true
	}

	for _, e := range m.FormatMock.expectations {
		if mm_atomic.LoadUint64(&e.Counter) < 1 {
			return false
		}
	}

	return m.FormatMock.invocationsDone()
}

// MinimockFormatInspect logs each unmet expectation
func (m *CustomFormatterNameMock) MinimockFormatInspect() {
	for _, e := range m.FormatMock.expectations {
		if mm_atomic.LoadUint64(&e.Counter) < 1 {
			m.t.Errorf("Expected call to CustomFormatterNameMock.Format at\n%s with params: %#v", e.expectationOrigins.origin, *e.params)
		}
	}

	afterFormatCounter := mm_atomic.LoadUint64(&m.afterFormatCounter)
	// if default expectation was set then invocations count should be greater than zero
	if m.FormatMock.defaultExpectation != nil && afterFormatCounter < 1 {
		if m.FormatMock.defaultExpectation.params == nil {
			m.t.Errorf("Expected call to CustomFormatterNameMock.Format at\n%s", m.FormatMock.defaultExpectation.returnOrigin)
		} else {
			m.t.Errorf("Expected call to CustomFormatterNameMock.Format at\n%s with params: %#v", m.FormatMock.defaultExpectation.expectationOrigins.origin, *m.FormatMock.defaultExpectation.params)
		}
	}
	// if func was set then invocations count should be greater than zero
	if m.funcFormat != nil && afterFormatCounter < 1 {
		m.t.Errorf("Expected call to CustomFormatterNameMock.Format at\n%s", m.funcFormatOrigin)
	}

	if !m.FormatMock.invocationsDone() && afterFormatCounter > 0 {
		m.t.Errorf("Expected %d calls to CustomFormatterNameMock.Format at\n%s but found %d calls",
			mm_atomic.LoadUint64(&m.FormatMock.expectedInvocations), m.FormatMock.expectedInvocationsOrigin, afterFormatCounter)
	}
}

// MinimockFinish checks that all mocked methods have been called the expected number of times
func (m *CustomFormatterNameMock) MinimockFinish() {
	m.finishOnce.Do(func() {
		if !m.minimockDone() {
			m.MinimockFormatInspect()
		}
	})
}

// MinimockWait waits for all mocked methods to be called the expected number of times
func (m *CustomFormatterNameMock) MinimockWait(timeout mm_time.Duration) {
	timeoutCh := mm_time.After(timeout)
	for {
		if m.minimockDone() {
			return
		}
		select {
		case <-timeoutCh:
			m.MinimockFinish()
			return
		case <-mm_time.After(10 * mm_time.Millisecond):
		}
	}
}

func (m *CustomFormatterNameMock) minimockDone() bool {
	done := true
	return done &&
		m.MinimockFormatDone()
}
