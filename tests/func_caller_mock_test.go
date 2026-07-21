package tests

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFuncCallerMock_AcceptsFunc(t *testing.T) {
	tester := NewTesterMock(t)
	tester.CleanupMock.Return().HelperMock.Return()

	validFunc := func() { _ = 1 }
	mock := NewFuncCallerMock(tester)
	mock.CallFuncMock.Expect(validFunc).Return(1)

	assert.Equal(t, 1, mock.CallFunc(validFunc))
}

func TestFuncCallerMock_ChecksNilFunc(t *testing.T) {
	tester := NewTesterMock(t)
	tester.CleanupMock.Return().HelperMock.Return()

	var f func()
	mock := NewFuncCallerMock(tester)
	mock.CallFuncMock.Expect(f).Return(1)

	assert.Equal(t, 1, mock.CallFunc(nil))
}

func TestFuncCallerMock_ChecksNil(t *testing.T) {
	tester := NewTesterMock(t)
	tester.CleanupMock.Return().HelperMock.Return()

	var f func()
	mock := NewFuncCallerMock(tester)
	mock.CallFuncMock.Expect(nil).Return(1)

	assert.Equal(t, 1, mock.CallFunc(f))
}

func TestFuncCallerMock_WhenThen(t *testing.T) {
	tester := NewTesterMock(t)
	tester.CleanupMock.Return().HelperMock.Return()

	var (
		f1 func()
		f2 = func() { _ = 2 }
	)
	mock := NewFuncCallerMock(tester)
	mock.
		CallFuncMock.When(f1).Then(1).
		CallFuncMock.When(f2).Then(2)

	assert.Equal(t, 1, mock.CallFunc(f1))
	assert.Equal(t, 2, mock.CallFunc(f2))
}

func TestFuncCallerMock_Method(t *testing.T) {
	tester := NewTesterMock(t)
	tester.CleanupMock.Return().HelperMock.Return()

	mock := NewFuncCallerMock(tester)

	var obj object
	mock.
		CallFuncMock.When(obj.Method1).Then(1).
		CallFuncMock.When(obj.Method2).Then(2)

	assert.Equal(t, 1, mock.CallFunc(obj.Method1))
	assert.Equal(t, 2, mock.CallFunc(obj.Method2))
}

func TestFuncCallerMock_SameVTable(t *testing.T) {
	tester := NewTesterMock(t)
	tester.CleanupMock.Return().HelperMock.Return()

	mock := NewFuncCallerMock(tester)

	var (
		obj1 object
		obj2 object
	)
	mock.CallFuncMock.Expect(obj1.Method1).Times(2).Return(1)

	assert.Equal(t, 1, mock.CallFunc(obj1.Method1))
	assert.Equal(t, 1, mock.CallFunc(obj2.Method1))
}

func TestFuncCallerMock_Closures(t *testing.T) {
	tester := NewTesterMock(t)
	tester.CleanupMock.Return().HelperMock.Return().ErrorfMock.Return()

	mock := NewFuncCallerMock(tester)

	mock.CallFuncMock.Expect(func() {}).Return(1)

	_ = mock.CallFunc(func() {})
}

func TestFuncCallerMock_ClosuresPartial(t *testing.T) {
	tester := NewTesterMock(t)
	tester.CleanupMock.Return().HelperMock.Return().ErrorfMock.Return()

	mock := NewFuncCallerMock(tester)

	mock.CallFuncMock.ExpectFParam1(func() {}).Return(1)

	_ = mock.CallFunc(func() {})
}

func TestFuncCallerMock_Partial(t *testing.T) {
	tester := NewTesterMock(t)
	tester.CleanupMock.Return().HelperMock.Return()

	mock := NewFuncCallerMock(tester)

	f := func() {}
	mock.CallFuncMock.ExpectFParam1(f).Return(1)

	assert.Equal(t, 1, mock.CallFunc(f))
}

func TestFuncCallerMock_Functions(t *testing.T) {
	tester := NewTesterMock(t)
	tester.CleanupMock.Return().HelperMock.Return().ErrorfMock.Return()

	mock := NewFuncCallerMock(tester)

	mock.CallFuncMock.Expect(function1).Return(1)

	assert.Equal(t, 1, mock.CallFunc(function1))
	mock.CallFunc(function2)
}

func TestFuncCallerMock_ClosuresFromFunc(t *testing.T) {
	tester := NewTesterMock(t)
	tester.CleanupMock.Return().HelperMock.Return().ErrorfMock.Return()

	mock := NewFuncCallerMock(tester)

	mock.CallFuncMock.Expect(newClosure()).Return(1)

	mock.CallFunc(newClosure())
}

type object struct{}

func (o object) Method1() {
}

func (o object) Method2() {
}

func function1() {}

func function2() {}

func newClosure() func() {
	return func() {}
}
