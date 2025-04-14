package minimock

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCheckFunctions(t *testing.T) {
	t.Parallel()
	var (
		validFunc        = func() {}
		nilFunc   func() = nil
	)

	tests := []struct {
		name string
		a    interface{}
		b    interface{}
		want bool
	}{{
		name: "one is nil, another is nil func",
		a:    nil,
		b:    nilFunc,
		want: false,
	}, {
		name: "arguments are nil",
		a:    nil,
		b:    nil,
		want: false,
	}, {
		name: "nil func and nil",
		a:    nilFunc,
		b:    nil,
		want: false,
	}, {
		name: "arguments are not functions",
		a:    1,
		b:    "string",
		want: false,
	}, {
		name: "both functions are nil",
		a:    nilFunc,
		b:    nilFunc,
		want: true,
	}, {
		name: "a is nil",
		a:    nilFunc,
		b:    validFunc,
		want: false,
	}, {
		name: "b is nil",
		a:    validFunc,
		b:    nilFunc,
		want: false,
	}, {
		name: "functions are anonymous",
		a:    func() {},
		b:    func() {},
		want: false,
	}, {
		name: "functions are equal",
		a:    validFunc,
		b:    validFunc,
		want: true,
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.want, checkFunctions(tt.a, tt.b))
		})
	}
}
