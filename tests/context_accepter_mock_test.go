package tests

import (
	"context"
	"testing"

	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/assert"
)

func TestContextAccepterMock_AnyContext(t *testing.T) {
	tester := NewTesterMock(t)

	var mockCalled bool
	tester.ErrorfMock.Set(func(s string, args ...interface{}) {
		assert.Equal(t, "ContextAccepterMock.AcceptContext got unexpected parameters, want: %#v, got: %#v%s\n", s)

		mockCalled = true
	}).CleanupMock.Return()

	mock := NewContextAccepterMock(tester).
		AcceptContextMock.Expect(context.Background()).Return()

	mock.AcceptContext(context.TODO())

	assert.True(t, mockCalled)
}

func TestContextAccepterMock_TodoContextMatchesAnycontext(t *testing.T) {
	tester := NewTesterMock(t)
	tester.CleanupMock.Return()

	mock := NewContextAccepterMock(tester).
		AcceptContextMock.Expect(minimock.AnyContext).Return()

	mock.AcceptContext(context.TODO())
}

func TestContextAccepterMock_WhenThenMatchAnycontext(t *testing.T) {
	tester := NewTesterMock(t)
	tester.CleanupMock.Return()

	mock := NewContextAccepterMock(tester).
		AcceptContextWithOtherArgsMock.When(minimock.AnyContext, 1).Then(42, nil)

	result, err := mock.AcceptContextWithOtherArgs(context.TODO(), 1)
	assert.NoError(t, err)
	assert.Equal(t, 42, result)
}

//func TestContextAccepterMock_Anycontext(t *testing.T) {
//	tester := NewTesterMock(t)
//	tester.CleanupMock.Return()
//
//	mock := NewContextAccepterMock(tester).
//		AcceptContextWithOtherArgsMock.Expect(minimock.AnyContext, 2422).Return(1, nil)
//
//	mock.AcceptContextWithOtherArgs(context.WithValue(context.TODO(), "a", 123), 123)
//}
