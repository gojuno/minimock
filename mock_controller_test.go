package minimock

import (
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewController(t *testing.T) {
	c := NewController(t)
	assert.Equal(t, t, c.Tester)
}

func TestController_RegisterMocker(t *testing.T) {
	c := &Controller{}
	c.RegisterMocker(nil)
	assert.Len(t, c.mockers, 1)
}

type dummyMocker struct {
	Mocker
	finishCounter int32
	waitCounter   int32
}

func (dm *dummyMocker) MinimockFinish() {
	atomic.AddInt32(&dm.finishCounter, 1)
}

func (dm *dummyMocker) MinimockWait(time.Duration) {
	atomic.AddInt32(&dm.waitCounter, 1)
}

func TestController_Finish(t *testing.T) {
	dm := &dummyMocker{}
	c := &Controller{
		mockers: []Mocker{dm, dm},
	}

	c.Finish()
	assert.Equal(t, int32(2), atomic.LoadInt32(&dm.finishCounter))
}

func TestController_Wait(t *testing.T) {
	dm := &dummyMocker{}
	c := &Controller{
		mockers: []Mocker{dm, dm},
	}

	c.Wait(0)
	assert.Equal(t, int32(2), atomic.LoadInt32(&dm.waitCounter))
}
