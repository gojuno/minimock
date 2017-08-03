package tests

import (
	"reflect"
	"testing"
	"time"

	"github.com/gojuno/minimock"
	"github.com/stretchr/testify/assert"
)

func TestStringerMock_ImplementsStringer(t *testing.T) {
	v := NewStringerMock(NewTesterMock(t))
	assert.True(t, reflect.TypeOf(v).Implements(reflect.TypeOf((*Stringer)(nil)).Elem()))
}

func TestStringerMock_UnmockedCallFailsTest(t *testing.T) {
	var mockCalled bool
	tester := NewTesterMock(t)
	tester.FatalFunc = func(args ...interface{}) {
		assert.Len(t, args, 1)
		assert.Equal(t, "Unexpected call to StringerMock.String", args[0])

		mockCalled = true
	}

	defer tester.Finish()

	m := NewStringerMock(tester)

	es := EmptyStringer{Stringer: m}
	assert.Equal(t, "empty string", es.String())
	assert.True(t, mockCalled)
}

func TestStringerMock_MockedCallSucceeds(t *testing.T) {
	tester := NewTesterMock(t)

	m := NewStringerMock(tester)
	m.StringFunc = func() string {
		return ""
	}
	defer tester.Finish()

	es := EmptyStringer{Stringer: m}
	assert.Equal(t, "empty string", es.String())
}

func TestStringerMock_Wait(t *testing.T) {
	tester := NewTesterMock(t)

	m := NewStringerMock(tester)
	m.StringFunc = func() string {
		return ""
	}

	go func() {
		es := EmptyStringer{Stringer: m}
		assert.Equal(t, "empty string", es.String())
	}()

	m.Wait(time.Second)
}

func TestStringerMock_MockReturn(t *testing.T) {
	tester := NewTesterMock(t)

	m := NewStringerMock(tester).StringMock.Return("Hello world!")
	defer m.Finish()

	es := EmptyStringer{Stringer: m}
	assert.Equal(t, "Hello world!", es.String())
}

func TestStringerMock_MockSet(t *testing.T) {
	tester := NewTesterMock(t)

	m := NewStringerMock(tester).StringMock.Set(func() string {
		return "set"
	})
	defer m.Finish()

	es := EmptyStringer{Stringer: m}
	assert.Equal(t, "set", es.String())
}

func TestStringerMock_AllMocksCalled(t *testing.T) {
	tester := NewTesterMock(t)

	m := NewStringerMock(tester).StringMock.Return("")
	assert.False(t, m.AllMocksCalled())

	assert.Equal(t, "", m.String())
	assert.True(t, m.AllMocksCalled())
}

func TestStringerMock_Finish(t *testing.T) {
	var mockCalled bool

	tester := NewTesterMock(t)
	tester.FatalMock.Set(func(args ...interface{}) {
		assert.Len(t, args, 1)
		assert.Equal(t, "Expected call to StringerMock.String", args[0])
		mockCalled = true
	})

	m := NewStringerMock(tester).StringMock.Return("")
	m.Finish()
	assert.True(t, mockCalled)
}

type dummyMockController struct {
	minimock.MockController
	registerCounter int
}

func (dmc *dummyMockController) RegisterMocker(m minimock.Mocker) {
	dmc.registerCounter++
}

func TestStringerMock_RegistersMocker(t *testing.T) {
	mockController := &dummyMockController{}

	NewStringerMock(mockController)
	assert.Equal(t, 1, mockController.registerCounter)
}
