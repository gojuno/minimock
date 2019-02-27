package minimock

import (
	"reflect"
)

// Equal returns true if a equals b
func Equal(a, b interface{}) bool {
	if a == nil && b == nil {
		return a == b
	}

	return reflect.DeepEqual(a, b)
}
