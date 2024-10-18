// Code generated by http://github.com/gojuno/minimock (dev). DO NOT EDIT.

package tests

//go:generate go run github.com/gojuno/minimock/v3/cmd/minimock -i github.com/gojuno/minimock/v3/tests.reader -o reader_mock.go -n ReaderMock -p tests -gr

import (
	"sync"
	mm_atomic "sync/atomic"
	mm_time "time"

	"github.com/gojuno/minimock/v3"
)

// ReaderMock implements reader
type ReaderMock struct {
	t          minimock.Tester
	finishOnce sync.Once

	funcRead          func(p []byte) (n int, err error)
	funcReadOrigin    string
	inspectFuncRead   func(p []byte)
	afterReadCounter  uint64
	beforeReadCounter uint64
	ReadMock          mReaderMockRead
}

// NewReaderMock returns a mock for reader
func NewReaderMock(t minimock.Tester) *ReaderMock {
	m := &ReaderMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.ReadMock = mReaderMockRead{mock: m}
	m.ReadMock.callArgs = []*ReaderMockReadParams{}

	t.Cleanup(m.MinimockFinish)

	return m
}

type mReaderMockRead struct {
	optional           bool
	mock               *ReaderMock
	defaultExpectation *ReaderMockReadExpectation
	expectations       []*ReaderMockReadExpectation

	callArgs []*ReaderMockReadParams
	mutex    sync.RWMutex

	expectedInvocations       uint64
	expectedInvocationsOrigin string
}

// ReaderMockReadExpectation specifies expectation struct of the reader.Read
type ReaderMockReadExpectation struct {
	mock               *ReaderMock
	params             *ReaderMockReadParams
	paramPtrs          *ReaderMockReadParamPtrs
	expectationOrigins ReaderMockReadExpectationOrigins
	results            *ReaderMockReadResults
	returnOrigin       string
	Counter            uint64
}

// ReaderMockReadParams contains parameters of the reader.Read
type ReaderMockReadParams struct {
	p []byte
}

// ReaderMockReadParamPtrs contains pointers to parameters of the reader.Read
type ReaderMockReadParamPtrs struct {
	p *[]byte
}

// ReaderMockReadResults contains results of the reader.Read
type ReaderMockReadResults struct {
	n   int
	err error
}

// ReaderMockReadOrigins contains origins of expectations of the reader.Read
type ReaderMockReadExpectationOrigins struct {
	origin  string
	originP string
}

// Marks this method to be optional. The default behavior of any method with Return() is '1 or more', meaning
// the test will fail minimock's automatic final call check if the mocked method was not called at least once.
// Optional() makes method check to work in '0 or more' mode.
// It is NOT RECOMMENDED to use this option unless you really need it, as default behaviour helps to
// catch the problems when the expected method call is totally skipped during test run.
func (mmRead *mReaderMockRead) Optional() *mReaderMockRead {
	mmRead.optional = true
	return mmRead
}

// Expect sets up expected params for reader.Read
func (mmRead *mReaderMockRead) Expect(p []byte) *mReaderMockRead {
	if mmRead.mock.funcRead != nil {
		mmRead.mock.t.Fatalf("ReaderMock.Read mock is already set by Set")
	}

	if mmRead.defaultExpectation == nil {
		mmRead.defaultExpectation = &ReaderMockReadExpectation{}
	}

	if mmRead.defaultExpectation.paramPtrs != nil {
		mmRead.mock.t.Fatalf("ReaderMock.Read mock is already set by ExpectParams functions")
	}

	mmRead.defaultExpectation.params = &ReaderMockReadParams{p}
	mmRead.defaultExpectation.expectationOrigins.origin = minimock.CallerInfo(1)
	for _, e := range mmRead.expectations {
		if minimock.Equal(e.params, mmRead.defaultExpectation.params) {
			mmRead.mock.t.Fatalf("Expectation set by When has same params: %#v", *mmRead.defaultExpectation.params)
		}
	}

	return mmRead
}

// ExpectPParam1 sets up expected param p for reader.Read
func (mmRead *mReaderMockRead) ExpectPParam1(p []byte) *mReaderMockRead {
	if mmRead.mock.funcRead != nil {
		mmRead.mock.t.Fatalf("ReaderMock.Read mock is already set by Set")
	}

	if mmRead.defaultExpectation == nil {
		mmRead.defaultExpectation = &ReaderMockReadExpectation{}
	}

	if mmRead.defaultExpectation.params != nil {
		mmRead.mock.t.Fatalf("ReaderMock.Read mock is already set by Expect")
	}

	if mmRead.defaultExpectation.paramPtrs == nil {
		mmRead.defaultExpectation.paramPtrs = &ReaderMockReadParamPtrs{}
	}
	mmRead.defaultExpectation.paramPtrs.p = &p
	mmRead.defaultExpectation.expectationOrigins.originP = minimock.CallerInfo(1)

	return mmRead
}

// Inspect accepts an inspector function that has same arguments as the reader.Read
func (mmRead *mReaderMockRead) Inspect(f func(p []byte)) *mReaderMockRead {
	if mmRead.mock.inspectFuncRead != nil {
		mmRead.mock.t.Fatalf("Inspect function is already set for ReaderMock.Read")
	}

	mmRead.mock.inspectFuncRead = f

	return mmRead
}

// Return sets up results that will be returned by reader.Read
func (mmRead *mReaderMockRead) Return(n int, err error) *ReaderMock {
	if mmRead.mock.funcRead != nil {
		mmRead.mock.t.Fatalf("ReaderMock.Read mock is already set by Set")
	}

	if mmRead.defaultExpectation == nil {
		mmRead.defaultExpectation = &ReaderMockReadExpectation{mock: mmRead.mock}
	}
	mmRead.defaultExpectation.results = &ReaderMockReadResults{n, err}
	mmRead.defaultExpectation.returnOrigin = minimock.CallerInfo(1)
	return mmRead.mock
}

// Set uses given function f to mock the reader.Read method
func (mmRead *mReaderMockRead) Set(f func(p []byte) (n int, err error)) *ReaderMock {
	if mmRead.defaultExpectation != nil {
		mmRead.mock.t.Fatalf("Default expectation is already set for the reader.Read method")
	}

	if len(mmRead.expectations) > 0 {
		mmRead.mock.t.Fatalf("Some expectations are already set for the reader.Read method")
	}

	mmRead.mock.funcRead = f
	mmRead.mock.funcReadOrigin = minimock.CallerInfo(1)
	return mmRead.mock
}

// When sets expectation for the reader.Read which will trigger the result defined by the following
// Then helper
func (mmRead *mReaderMockRead) When(p []byte) *ReaderMockReadExpectation {
	if mmRead.mock.funcRead != nil {
		mmRead.mock.t.Fatalf("ReaderMock.Read mock is already set by Set")
	}

	expectation := &ReaderMockReadExpectation{
		mock:               mmRead.mock,
		params:             &ReaderMockReadParams{p},
		expectationOrigins: ReaderMockReadExpectationOrigins{origin: minimock.CallerInfo(1)},
	}
	mmRead.expectations = append(mmRead.expectations, expectation)
	return expectation
}

// Then sets up reader.Read return parameters for the expectation previously defined by the When method
func (e *ReaderMockReadExpectation) Then(n int, err error) *ReaderMock {
	e.results = &ReaderMockReadResults{n, err}
	return e.mock
}

// Times sets number of times reader.Read should be invoked
func (mmRead *mReaderMockRead) Times(n uint64) *mReaderMockRead {
	if n == 0 {
		mmRead.mock.t.Fatalf("Times of ReaderMock.Read mock can not be zero")
	}
	mm_atomic.StoreUint64(&mmRead.expectedInvocations, n)
	mmRead.expectedInvocationsOrigin = minimock.CallerInfo(1)
	return mmRead
}

func (mmRead *mReaderMockRead) invocationsDone() bool {
	if len(mmRead.expectations) == 0 && mmRead.defaultExpectation == nil && mmRead.mock.funcRead == nil {
		return true
	}

	totalInvocations := mm_atomic.LoadUint64(&mmRead.mock.afterReadCounter)
	expectedInvocations := mm_atomic.LoadUint64(&mmRead.expectedInvocations)

	return totalInvocations > 0 && (expectedInvocations == 0 || expectedInvocations == totalInvocations)
}

// Read implements reader
func (mmRead *ReaderMock) Read(p []byte) (n int, err error) {
	mm_atomic.AddUint64(&mmRead.beforeReadCounter, 1)
	defer mm_atomic.AddUint64(&mmRead.afterReadCounter, 1)

	mmRead.t.Helper()

	if mmRead.inspectFuncRead != nil {
		mmRead.inspectFuncRead(p)
	}

	mm_params := ReaderMockReadParams{p}

	// Record call args
	mmRead.ReadMock.mutex.Lock()
	mmRead.ReadMock.callArgs = append(mmRead.ReadMock.callArgs, &mm_params)
	mmRead.ReadMock.mutex.Unlock()

	for _, e := range mmRead.ReadMock.expectations {
		if minimock.Equal(*e.params, mm_params) {
			mm_atomic.AddUint64(&e.Counter, 1)
			return e.results.n, e.results.err
		}
	}

	if mmRead.ReadMock.defaultExpectation != nil {
		mm_atomic.AddUint64(&mmRead.ReadMock.defaultExpectation.Counter, 1)
		mm_want := mmRead.ReadMock.defaultExpectation.params
		mm_want_ptrs := mmRead.ReadMock.defaultExpectation.paramPtrs

		mm_got := ReaderMockReadParams{p}

		if mm_want_ptrs != nil {

			if mm_want_ptrs.p != nil && !minimock.Equal(*mm_want_ptrs.p, mm_got.p) {
				mmRead.t.Errorf("ReaderMock.Read got unexpected parameter p, expected at\n%s:\nwant: %#v\n got: %#v%s\n",
					mmRead.ReadMock.defaultExpectation.expectationOrigins.originP, *mm_want_ptrs.p, mm_got.p, minimock.Diff(*mm_want_ptrs.p, mm_got.p))
			}

		} else if mm_want != nil && !minimock.Equal(*mm_want, mm_got) {
			mmRead.t.Errorf("ReaderMock.Read got unexpected parameters, expected at\n%s:\nwant: %#v\n got: %#v%s\n",
				mmRead.ReadMock.defaultExpectation.expectationOrigins.origin, *mm_want, mm_got, minimock.Diff(*mm_want, mm_got))
		}

		mm_results := mmRead.ReadMock.defaultExpectation.results
		if mm_results == nil {
			mmRead.t.Fatal("No results are set for the ReaderMock.Read")
		}
		return (*mm_results).n, (*mm_results).err
	}
	if mmRead.funcRead != nil {
		return mmRead.funcRead(p)
	}
	mmRead.t.Fatalf("Unexpected call to ReaderMock.Read. %v", p)
	return
}

// ReadAfterCounter returns a count of finished ReaderMock.Read invocations
func (mmRead *ReaderMock) ReadAfterCounter() uint64 {
	return mm_atomic.LoadUint64(&mmRead.afterReadCounter)
}

// ReadBeforeCounter returns a count of ReaderMock.Read invocations
func (mmRead *ReaderMock) ReadBeforeCounter() uint64 {
	return mm_atomic.LoadUint64(&mmRead.beforeReadCounter)
}

// Calls returns a list of arguments used in each call to ReaderMock.Read.
// The list is in the same order as the calls were made (i.e. recent calls have a higher index)
func (mmRead *mReaderMockRead) Calls() []*ReaderMockReadParams {
	mmRead.mutex.RLock()

	argCopy := make([]*ReaderMockReadParams, len(mmRead.callArgs))
	copy(argCopy, mmRead.callArgs)

	mmRead.mutex.RUnlock()

	return argCopy
}

// MinimockReadDone returns true if the count of the Read invocations corresponds
// the number of defined expectations
func (m *ReaderMock) MinimockReadDone() bool {
	if m.ReadMock.optional {
		// Optional methods provide '0 or more' call count restriction.
		return true
	}

	for _, e := range m.ReadMock.expectations {
		if mm_atomic.LoadUint64(&e.Counter) < 1 {
			return false
		}
	}

	return m.ReadMock.invocationsDone()
}

// MinimockReadInspect logs each unmet expectation
func (m *ReaderMock) MinimockReadInspect() {
	for _, e := range m.ReadMock.expectations {
		if mm_atomic.LoadUint64(&e.Counter) < 1 {
			m.t.Errorf("Expected call to ReaderMock.Read at\n%s with params: %#v", e.expectationOrigins.origin, *e.params)
		}
	}

	afterReadCounter := mm_atomic.LoadUint64(&m.afterReadCounter)
	// if default expectation was set then invocations count should be greater than zero
	if m.ReadMock.defaultExpectation != nil && afterReadCounter < 1 {
		if m.ReadMock.defaultExpectation.params == nil {
			m.t.Errorf("Expected call to ReaderMock.Read at\n%s", m.ReadMock.defaultExpectation.returnOrigin)
		} else {
			m.t.Errorf("Expected call to ReaderMock.Read at\n%s with params: %#v", m.ReadMock.defaultExpectation.expectationOrigins.origin, *m.ReadMock.defaultExpectation.params)
		}
	}
	// if func was set then invocations count should be greater than zero
	if m.funcRead != nil && afterReadCounter < 1 {
		m.t.Errorf("Expected call to ReaderMock.Read at\n%s", m.funcReadOrigin)
	}

	if !m.ReadMock.invocationsDone() && afterReadCounter > 0 {
		m.t.Errorf("Expected %d calls to ReaderMock.Read at\n%s but found %d calls",
			mm_atomic.LoadUint64(&m.ReadMock.expectedInvocations), m.ReadMock.expectedInvocationsOrigin, afterReadCounter)
	}
}

// MinimockFinish checks that all mocked methods have been called the expected number of times
func (m *ReaderMock) MinimockFinish() {
	m.finishOnce.Do(func() {
		if !m.minimockDone() {
			m.MinimockReadInspect()
		}
	})
}

// MinimockWait waits for all mocked methods to be called the expected number of times
func (m *ReaderMock) MinimockWait(timeout mm_time.Duration) {
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

func (m *ReaderMock) minimockDone() bool {
	done := true
	return done &&
		m.MinimockReadDone()
}
