package tests

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestActorMock_TestPassedWithBothExpectedParams(t *testing.T) {
	tester := NewTesterMock(t)
	tester.CleanupMock.Return().HelperMock.Return()

	mock := NewActorMock(tester).
		ActionMock.ExpectFirstParamParam1("abc").
		ExpectSecondParamParam2(24).Return(1, nil)

	a, err := mock.Action("abc", 24)
	assert.NoError(t, err)
	assert.Equal(t, 1, a)
}

func TestActorMock_TestPassedWithOneExpectedParams(t *testing.T) {
	tester := NewTesterMock(t)
	tester.CleanupMock.Return().HelperMock.Return()

	mock := NewActorMock(tester).
		ActionMock.ExpectFirstParamParam1("abc").Return(1, nil)

	a, err := mock.Action("abc", 24)
	assert.NoError(t, err)
	assert.Equal(t, 1, a)
}

func TestActorMock_TestFailedWithExpectedParams(t *testing.T) {
	tester := NewTesterMock(t)
	tester.CleanupMock.Return().HelperMock.Return()
	tester.ErrorfMock.Set(func(format string, args ...interface{}) {
		assert.Equal(t, "ActorMock.Action got unexpected parameter secondParam, expected at\n%s:\nwant: %#v\n got: %#v%s\n", format)

		assert.Equal(t, 24, args[1])
		assert.Equal(t, 25, args[2])
		assert.Equal(t, "", args[3])
	})
	mock := NewActorMock(tester).
		ActionMock.ExpectFirstParamParam1("abc").
		ExpectSecondParamParam2(24).Return(1, nil)

	a, err := mock.Action("abc", 25)
	assert.NoError(t, err)
	assert.Equal(t, 1, a)
}

func TestActorMock_FailedToUseExpectAfterExpectParams(t *testing.T) {
	tester := NewTesterMock(t)
	tester.CleanupMock.Return()
	tester.FatalfMock.
		Expect("ActorMock.Action mock is already set by ExpectParams functions").
		Return()

	_ = NewActorMock(tester).
		ActionMock.ExpectFirstParamParam1("abc").
		Expect("aaa", 123).Return(1, nil)
}

func TestActorMock_FailedToUseExpectParamsAfterExpect(t *testing.T) {
	tester := NewTesterMock(t)
	tester.CleanupMock.Return()
	tester.FatalfMock.
		Expect("ActorMock.Action mock is already set by Expect").
		Return()

	_ = NewActorMock(tester).
		ActionMock.Expect("aaa", 123).
		ExpectFirstParamParam1("abc").Return(1, nil)
}

func TestActorMock_FailedToUseExpectParamsAfterSet(t *testing.T) {
	tester := NewTesterMock(t)
	tester.CleanupMock.Return()
	tester.FatalfMock.
		Expect("ActorMock.Action mock is already set by Set").
		Return()

	_ = NewActorMock(tester).
		ActionMock.Set(func(firstParam string, secondParam int) (i1 int, err error) {
		return
	}).ActionMock.ExpectFirstParamParam1("abc").Return(1, nil)
}

func TestActorMock_Optional(t *testing.T) {
	tester := NewTesterMock(t)
	tester.CleanupMock.Return()

	mock := NewActorMock(tester).ActionMock.Optional().Return(1, nil)

	mock.MinimockFinish()
}
