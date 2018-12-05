package tests

import (
	"reflect"
	"testing"
	"time"

	"github.com/gojuno/minimock"
	"github.com/stretchr/testify/assert"
)

func TestFormatterMock_ImplementsStringer(t *testing.T) {
	v := NewFormatterMock(NewTesterMock(t))
	assert.True(t, reflect.TypeOf(v).Implements(reflect.TypeOf((*Formatter)(nil)).Elem()))
}

func TestFormatterMock_UnmockedCallFailsTest(t *testing.T) {
	var mockCalled bool
	tester := NewTesterMock(t)
	tester.FatalfFunc = func(s string, args ...interface{}) {
		assert.Equal(t, "Unexpected call to FormatterMock.Format. %v %v", s)
		assert.Equal(t, "this call fails because Format method isn't mocked", args[0])

		mockCalled = true
	}

	defer tester.MinimockFinish()

	formatterMock := NewFormatterMock(tester)
	dummyFormatter{formatterMock}.Format("this call fails because Format method isn't mocked")
	assert.True(t, mockCalled)
}

func TestFormatterMock_MockedCallSucceeds(t *testing.T) {
	tester := NewTesterMock(t)

	formatterMock := NewFormatterMock(tester)
	formatterMock.FormatFunc = func(format string, args ...interface{}) string {
		return "mock is successfully called"
	}
	defer tester.MinimockFinish()

	df := dummyFormatter{formatterMock}
	assert.Equal(t, "mock is successfully called", df.Format(""))
}

func TestFormatterMock_Wait(t *testing.T) {
	tester := NewTesterMock(t)

	formatterMock := NewFormatterMock(tester)
	formatterMock.FormatFunc = func(format string, args ...interface{}) string {
		return "mock is successfully called from the goroutine"
	}

	go func() {
		df := dummyFormatter{formatterMock}
		assert.Equal(t, "mock is successfully called from the goroutine", df.Format(""))
	}()

	formatterMock.MinimockWait(time.Second)
}

func TestFormatterMock_Expect(t *testing.T) {
	tester := NewTesterMock(t)

	formatterMock := NewFormatterMock(tester).FormatMock.Expect("Hello", "world", "!").Return("")

	df := dummyFormatter{formatterMock}
	df.Format("Hello", "world", "!")
}

func TestFormatterMock_Return(t *testing.T) {
	tester := NewTesterMock(t)

	formatterMock := NewFormatterMock(tester).FormatMock.Return("Hello world!")
	df := dummyFormatter{formatterMock}
	assert.Equal(t, "Hello world!", df.Format(""))
}

func TestFormatterMock_Set(t *testing.T) {
	tester := NewTesterMock(t)

	formatterMock := NewFormatterMock(tester).FormatMock.Set(func(string, ...interface{}) string {
		return "set"
	})

	df := dummyFormatter{formatterMock}
	assert.Equal(t, "set", df.Format(""))
}

func TestFormatterMock_AllMocksCalled(t *testing.T) {
	tester := NewTesterMock(t)

	formatterMock := NewFormatterMock(tester).FormatMock.Return("all mocks called")
	assert.False(t, formatterMock.AllMocksCalled())

	assert.Equal(t, "all mocks called", formatterMock.Format(""))
	assert.True(t, formatterMock.AllMocksCalled())
}

func TestFormatterMock_Finish(t *testing.T) {
	var mockCalled bool

	tester := NewTesterMock(t)
	tester.FatalMock.Set(func(args ...interface{}) {
		assert.Len(t, args, 1)
		assert.Equal(t, "Expected call to FormatterMock.Format", args[0])
		mockCalled = true
	})

	formatterMock := NewFormatterMock(tester).FormatMock.Return("")
	formatterMock.Finish()
	assert.True(t, mockCalled)
}

func TestMFormatterMockFormat_SeriesExpect(t *testing.T) {
	tester := NewTesterMock(t)

	formatter := NewFormatterMock(tester)
	defer formatter.Finish()

	formatter.FormatMock.SeriesExpect("hello %v", "username").Return("hello username")
	formatter.FormatMock.SeriesExpect("goodbye %v", "username").Return("goodbye username")

	assert.Equal(t, "hello username", formatter.Format("hello %v", "username"))
	assert.Equal(t, "goodbye username", formatter.Format("goodbye %v", "username"))
}

func TestMFormatterMockFormat_SeriesExpect_NoReturn(t *testing.T) {
	tester := NewTesterMock(t)

	formatter := NewFormatterMock(tester)
	defer formatter.Finish()

	formatter.FormatMock.SeriesExpect("hello %v", "username").Return("hello username")
	formatter.FormatMock.SeriesExpect("goodbye %v", "username") // no return here will produce error

	// return is set for this invocation
	formatter.Format("hello %v", "username")

	// return is not set for this invocation
	tester.FatalMock.SeriesExpect("No results are set for the FormatterMock.Format")
	formatter.Format("goodbye %v", "username")
}

func TestMFormatterMockFormat_SeriesExpect_NotEnoughCalls(t *testing.T) {
	tester := NewTesterMock(t)

	formatter := NewFormatterMock(tester)

	formatter.FormatMock.SeriesExpect("hello %v", "username").Return("hello username")
	formatter.FormatMock.SeriesExpect("goodbye %v", "username")

	formatter.Format("hello %v", "username")

	// expected two invocations of Format, but did only one
	tester.FatalMock.SeriesExpect("Expected call to FormatterMock.Format")
	formatter.Finish()
}

func TestMFormatterMockFormat_SeriesExpect_TooMuchCalls(t *testing.T) {
	tester := NewTesterMock(t)

	formatter := NewFormatterMock(tester)

	formatter.FormatMock.SeriesExpect("hello %v", "username").Return("hello username")
	formatter.FormatMock.SeriesExpect("goodbye %v", "username").Return("goodbye username")

	formatter.Format("hello %v", "username")
	formatter.Format("goodbye %v", "username")

	// expected two invocations to Format, but do three
	tester.FatalfMock.SeriesExpect("Unexpected call to FormatterMock.Format. %v %v", "hello again", []interface{}{"username"})
	formatter.Format("hello again", "username")
}

type dummyFormatter struct {
	Formatter
}

type dummyMockController struct {
	minimock.MockController
	registerCounter int
}

func (dmc *dummyMockController) RegisterMocker(m minimock.Mocker) {
	dmc.registerCounter++
}

func TestFormatterMock_RegistersMocker(t *testing.T) {
	mockController := &dummyMockController{}

	NewFormatterMock(mockController)
	assert.Equal(t, 1, mockController.registerCounter)
}
