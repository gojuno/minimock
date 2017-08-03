package minimock

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewController(t *testing.T) {
	c := NewController(t)
	assert.Equal(t, c, &Controller{T: t})
}

func TestController_RegisterMocker(t *testing.T) {
	c := &Controller{}
	c.RegisterMocker(nil)
	assert.Len(t, c.mockers, 1)
}

type dummyMocker struct {
	Mocker
	finishCounter int
	waitCounter   int
}

func (dm *dummyMocker) MinimockFinish() {
	dm.finishCounter++
}

func (dm *dummyMocker) MinimockWait(time.Duration) {
	dm.waitCounter++
}

func TestController_Finish(t *testing.T) {
	dm := &dummyMocker{}
	c := &Controller{
		mockers: []Mocker{dm, dm},
	}

	c.Finish()
	assert.Equal(t, 2, dm.finishCounter)
}

func TestController_Wait(t *testing.T) {
	dm := &dummyMocker{}
	c := &Controller{
		mockers: []Mocker{dm, dm},
	}

	c.Wait(0)
	assert.Equal(t, 2, dm.waitCounter)
}
