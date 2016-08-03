//go:generate minimock -f ./sample/interface.go -i Interface -o ./sample_interface_mock.go -p examples
package examples

import "testing"

func TestSampleInterfaceUser(t *testing.T) {
	mock := NewInterfaceMock(t)

	mock.GetStringFunc = func() string {
		return "Hello"
	}
	mock.CalculateSumFunc = func(ints ...int) int {
		return 1
	}
	mock.GetArrayOfStringsFunc = func(s string, i int) []*string {
		return []*string{&s}
	}

	res := SampleInterfaceUser(mock)
	if len(res) != 1 {
		t.Fatalf("expected: %d, got: %d", 1, len(res))
	}

	if res[0] == nil {
		t.Fatalf("expected non-nil value")
	}

	if *res[0] != "Hello" {
		t.Errorf("expected: Hello, got: %s", *res[0])
	}

	//checks that all mocked functions have been called,
	mock.ValidateCallCounters()

	//you can perform more precise check of the calls counter
	if mock.GetStringCounter != 1 {
		t.Errorf("expected one call to GetString, got %d", mock.GetStringCounter)
	}
}
