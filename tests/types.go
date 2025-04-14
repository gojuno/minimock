// Package tests contains tests for minimock tool and demonstrates minimock features
package tests

import (
	"context"
	"io"

	"github.com/gofrs/uuid/v5"
	"google.golang.org/protobuf/proto"
)

type (
	// Alias for package with version
	gen = uuid.Gen

	//Formatter interface is used to test code generated by minimock
	Formatter interface {
		formatter //to check if variadic functions are supported
	}

	formatter interface {
		Format(string, ...interface{}) string //to check if variadic functions are supported
	}

	// these generic types provide all possible cases of type params declarations
	// This produces invalid Go code when used with minimock v3.0.10 and below
	genericInout[T any] interface {
		Name(T) T
	}
	genericOut[T any] interface {
		Name() T
	}

	genericIn[T any] interface {
		Name(T)
	}

	// following types reference some specific constraints on generic types
	// These are not generated properly with minimock v3.1.1 and below

	// Reference a specific type
	genericSpecific[T proto.Message] interface {
		Name(T)
	}

	// Reference a single type as a simple union
	simpleUnion interface {
		int
	}

	// Reference a composite union of multiple types
	complexUnion interface {
		int | float64
	}

	genericSimpleUnion[T simpleUnion] interface {
		Name(T)
	}

	genericComplexUnion[T complexUnion] interface {
		Name(T)
	}

	genericInlineUnion[T int | float64] interface {
		Name(T)
	}

	genericInlineUnionWithManyTypes[T int | float64 | string] interface {
		Name(T)
	}

	genericMultipleTypes[T proto.Message, K any] interface {
		Name(T, K)
	}

	formatterAlias = Formatter

	formatterType Formatter

	reader = io.Reader

	structArg struct {
		a int
		b string
	}

	contextAccepter interface {
		AcceptContext(context.Context)
		AcceptContextWithOtherArgs(context.Context, int) (int, error)
		AcceptContextWithStructArgs(context.Context, structArg) (int, error)
	}

	actor interface {
		Action(firstParam string, secondParam int) (int, error)
	}

	funcCaller interface {
		CallFunc(f func()) int
	}
)
