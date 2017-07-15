//Package tests contains tests for minimock tool and demonstrates minimock features
package tests

import "fmt"

type (
	//Stringer embeds fmt.Stringer interface to check that minimock processes embedded interfaces correctly
	Stringer interface {
		fmt.Stringer
	}

	//EmptyStringer implements Stringer
	EmptyStringer struct {
		Stringer Stringer
	}
)

//String method returns "empty string" string in case when es.Stringer returns an empty string
func (es EmptyStringer) String() string {
	if result := es.Stringer.String(); result != "" {
		return result
	}
	return "empty string"
}
