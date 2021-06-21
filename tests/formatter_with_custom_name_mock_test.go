package tests

import (
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCustomFormatterNameMock_ImplementsStringer(t *testing.T) {
	v := NewCustomFormatterNameMock(NewTesterMock(t))
	assert.True(t, reflect.TypeOf(v).Implements(reflect.TypeOf((*Formatter)(nil)).Elem()))
}

func TestCustomFormatterNameMock_UnmockedCallFailsTest(t *testing.T) {
	var mockCalled bool
	tester := NewTesterMock(t)
	tester.FatalfMock.Set(func(s string, args ...interface{}) {
		assert.Equal(t, "Unexpected call to CustomFormatterNameMock.Format. %v %v", s)
		assert.Equal(t, "this call fails because Format method isn't mocked", args[0])

		mockCalled = true
	})

	defer tester.MinimockFinish()

	formatterMock := NewCustomFormatterNameMock(tester)
	dummyFormatter{formatterMock}.Format("this call fails because Format method isn't mocked")
	assert.True(t, mockCalled)
}

func TestCustomFormatterNameMock_MockedCallSucceeds(t *testing.T) {
	tester := NewTesterMock(t)

	formatterMock := NewCustomFormatterNameMock(tester)
	formatterMock.FormatMock.Set(func(format string, args ...interface{}) string {
		return "mock is successfully called"
	})
	defer tester.MinimockFinish()

	df := dummyFormatter{formatterMock}
	assert.Equal(t, "mock is successfully called", df.Format(""))
}

func TestCustomFormatterNameMock_Wait(t *testing.T) {
	tester := NewTesterMock(t)

	formatterMock := NewCustomFormatterNameMock(tester)
	formatterMock.FormatMock.Set(func(format string, args ...interface{}) string {
		return "mock is successfully called from the goroutine"
	})

	go func() {
		df := dummyFormatter{formatterMock}
		assert.Equal(t, "mock is successfully called from the goroutine", df.Format(""))
	}()

	formatterMock.MinimockWait(time.Second)
}

func TestCustomFormatterNameMock_Expect(t *testing.T) {
	tester := NewTesterMock(t)

	formatterMock := NewCustomFormatterNameMock(tester).FormatMock.Expect("Hello", "world", "!").Return("")

	df := dummyFormatter{formatterMock}
	df.Format("Hello", "world", "!")

	assert.EqualValues(t, 1, formatterMock.FormatBeforeCounter())
	assert.EqualValues(t, 1, formatterMock.FormatAfterCounter())
}

func TestCustomFormatterNameMock_ExpectDifferentArguments(t *testing.T) {
	assert.Panics(t, func() {
		tester := NewTesterMock(t)
		defer tester.MinimockFinish()

		tester.ErrorfMock.Set(func(s string, args ...interface{}) {
			assert.Equal(t, "CustomFormatterNameMock.Format got unexpected parameters, want: %#v, got: %#v%s\n", s)
			require.Len(t, args, 3)
			assert.Equal(t, CustomFormatterNameMockFormatParams{s1: "expected"}, args[0])
			assert.Equal(t, CustomFormatterNameMockFormatParams{s1: "actual"}, args[1])
		})

		tester.FatalMock.Expect("No results are set for the CustomFormatterNameMock.Format").Return()

		formatterMock := NewCustomFormatterNameMock(tester)
		formatterMock.FormatMock.Expect("expected")
		formatterMock.Format("actual")
	})
}

func TestCustomFormatterNameMock_ExpectAfterSet(t *testing.T) {
	tester := NewTesterMock(t)
	defer tester.MinimockFinish()

	tester.FatalfMock.Expect("CustomFormatterNameMock.Format mock is already set by Set").Return()

	formatterMock := NewCustomFormatterNameMock(tester)
	formatterMock.FormatMock.Set(func(string, ...interface{}) string { return "" })

	formatterMock.FormatMock.Expect("Should not work")
}

func TestCustomFormatterNameMock_ExpectAfterWhen(t *testing.T) {
	tester := NewTesterMock(t)
	defer tester.MinimockFinish()

	tester.FatalfMock.Expect("Expectation set by When has same params: %#v", CustomFormatterNameMockFormatParams{s1: "Should not work", p1: nil}).Return()

	formatterMock := NewCustomFormatterNameMock(tester)
	formatterMock.FormatMock.When("Should not work").Then("")

	formatterMock.Format("Should not work")

	formatterMock.FormatMock.Expect("Should not work")
}

func TestCustomFormatterNameMock_Return(t *testing.T) {
	tester := NewTesterMock(t)

	formatterMock := NewCustomFormatterNameMock(tester).FormatMock.Return("Hello world!")
	df := dummyFormatter{formatterMock}
	assert.Equal(t, "Hello world!", df.Format(""))
}

func TestCustomFormatterNameMock_ReturnAfterSet(t *testing.T) {
	tester := NewTesterMock(t)
	defer tester.MinimockFinish()

	tester.FatalfMock.Expect("CustomFormatterNameMock.Format mock is already set by Set").Return()

	formatterMock := NewCustomFormatterNameMock(tester)
	formatterMock.FormatMock.Set(func(string, ...interface{}) string { return "" })

	formatterMock.FormatMock.Return("Should not work")
}

func TestCustomFormatterNameMock_ReturnWithoutExpectForFixedArgsMethod(t *testing.T) {
	// Test for issue https://github.com/gojuno/minimock/issues/31

	tester := NewTesterMock(t)
	defer tester.MinimockFinish()

	tester.ErrorMock.Expect("Expected call to CustomFormatterNameMock.Format")
	tester.FailNowMock.Expect()

	formatterMock := NewCustomFormatterNameMock(tester)
	formatterMock.FormatMock.Return("")
	formatterMock.MinimockFinish()
}

func TestCustomFormatterNameMock_Set(t *testing.T) {
	tester := NewTesterMock(t)

	formatterMock := NewCustomFormatterNameMock(tester).FormatMock.Set(func(string, ...interface{}) string {
		return "set"
	})

	df := dummyFormatter{formatterMock}
	assert.Equal(t, "set", df.Format(""))
}

func TestCustomFormatterNameMock_SetAfterExpect(t *testing.T) {
	tester := NewTesterMock(t)
	defer tester.MinimockFinish()

	tester.FatalfMock.Expect("Default expectation is already set for the Formatter.Format method").Return()

	formatterMock := NewCustomFormatterNameMock(tester).FormatMock.Expect("").Return("")

	//second attempt should fail
	formatterMock.FormatMock.Set(func(string, ...interface{}) string { return "" })
}

func TestCustomFormatterNameMock_SetAfterWhen(t *testing.T) {
	tester := NewTesterMock(t)
	defer tester.MinimockFinish()

	tester.FatalfMock.Expect("Some expectations are already set for the Formatter.Format method").Return()

	formatterMock := NewCustomFormatterNameMock(tester).FormatMock.When("").Then("")

	//second attempt should fail
	formatterMock.FormatMock.Set(func(string, ...interface{}) string { return "" })
}

func TestCustomFormatterNameMockFormat_WhenThen(t *testing.T) {
	formatter := NewCustomFormatterNameMock(t)
	defer formatter.MinimockFinish()

	formatter.FormatMock.When("hello %v", "username").Then("hello username")
	formatter.FormatMock.When("goodbye %v", "username").Then("goodbye username")

	assert.Equal(t, "hello username", formatter.Format("hello %v", "username"))
	assert.Equal(t, "goodbye username", formatter.Format("goodbye %v", "username"))
}

func TestCustomFormatterNameMockFormat_WhenAfterSet(t *testing.T) {
	tester := NewTesterMock(t)
	defer tester.MinimockFinish()

	tester.FatalfMock.Expect("CustomFormatterNameMock.Format mock is already set by Set").Return()

	formatterMock := NewCustomFormatterNameMock(tester)
	formatterMock.FormatMock.Set(func(string, ...interface{}) string { return "" })

	formatterMock.FormatMock.When("Should not work")
}

func TestCustomFormatterNameMock_MinimockFormatDone(t *testing.T) {
	formatterMock := NewCustomFormatterNameMock(t)

	formatterMock.FormatMock.expectations = []*CustomFormatterNameMockFormatExpectation{{}}
	assert.False(t, formatterMock.MinimockFormatDone())

	formatterMock = NewCustomFormatterNameMock(t)
	formatterMock.FormatMock.defaultExpectation = &CustomFormatterNameMockFormatExpectation{}
	assert.False(t, formatterMock.MinimockFormatDone())
}

func TestCustomFormatterNameMock_MinimockFinish(t *testing.T) {
	tester := NewTesterMock(t)
	defer tester.MinimockFinish()

	tester.ErrorMock.Expect("Expected call to CustomFormatterNameMock.Format").Return()
	tester.FailNowMock.Expect().Return()

	formatterMock := NewCustomFormatterNameMock(tester)
	formatterMock.FormatMock.Set(func(string, ...interface{}) string { return "" })

	formatterMock.MinimockFinish()
}

func TestCustomFormatterNameMock_MinimockFinish_WithNoMetExpectations(t *testing.T) {
	tester := NewTesterMock(t)
	defer tester.MinimockFinish()

	tester.ErrorfMock.Set(func(m string, args ...interface{}) {
		assert.Equal(t, m, "Expected call to CustomFormatterNameMock.Format with params: %#v")
	})
	tester.FailNowMock.Expect().Return()

	formatterMock := NewCustomFormatterNameMock(tester)
	formatterMock.FormatMock.Expect("a").Return("a")
	formatterMock.FormatMock.When("b").Then("b")

	formatterMock.MinimockFinish()
}

func TestCustomFormatterNameMock_MinimockWait(t *testing.T) {
	tester := NewTesterMock(t)
	defer tester.MinimockFinish()

	tester.ErrorMock.Expect("Expected call to CustomFormatterNameMock.Format").Return()
	tester.FailNowMock.Expect().Return()

	formatterMock := NewCustomFormatterNameMock(tester)
	formatterMock.FormatMock.Set(func(string, ...interface{}) string { return "" })

	formatterMock.MinimockWait(time.Millisecond)
}

// Verifies that Calls() doesn't return nil if no calls were made
func TestCustomFormatterNameMock_CallsNotNil(t *testing.T) {
	tester := NewTesterMock(t)
	defer tester.MinimockFinish()

	formatterMock := NewCustomFormatterNameMock(tester)
	calls := formatterMock.FormatMock.Calls()

	assert.NotNil(t, calls)
	assert.Empty(t, calls)
}

// Verifies that Calls() returns the correct call args in the expected order
func TestCustomFormatterNameMock_Calls(t *testing.T) {
	tester := NewTesterMock(t)
	defer tester.MinimockFinish()

	// Arguments used for each mock call
	expected := []*CustomFormatterNameMockFormatParams{
		{"a1", []interface{}{}},
		{"b1", []interface{}{"b2"}},
		{"c1", []interface{}{"c2", "c3"}},
		{"d1", []interface{}{"d2", "d3", "d4"}},
	}

	formatterMock := NewCustomFormatterNameMock(tester)

	for _, p := range expected {
		formatterMock.FormatMock.Expect(p.s1, p.p1...).Return("")
		formatterMock.Format(p.s1, p.p1...)
	}

	assert.Equal(t, expected, formatterMock.FormatMock.Calls())
}

// Verifies that Calls() returns a new shallow copy of the params list each time
func TestCustomFormatterNameMock_CallsReturnsCopy(t *testing.T) {
	tester := NewTesterMock(t)
	defer tester.MinimockFinish()

	expected := []*CustomFormatterNameMockFormatParams{
		{"a1", []interface{}{"a1"}},
		{"b1", []interface{}{"b2"}},
	}

	formatterMock := NewCustomFormatterNameMock(tester)
	callHistory := [][]*CustomFormatterNameMockFormatParams{}

	for _, p := range expected {
		formatterMock.FormatMock.Expect(p.s1, p.p1...).Return("")
		formatterMock.Format(p.s1, p.p1...)
		callHistory = append(callHistory, formatterMock.FormatMock.Calls())
	}

	assert.Equal(t, len(expected), len(callHistory))

	for i, c := range callHistory {
		assert.Equal(t, i+1, len(c))
	}
}
