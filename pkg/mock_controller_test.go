package minimock

import (
	"sync/atomic"
	"testing"
	"time"
)

func TestNewController(t *testing.T) {
	c := NewController(t)
	if !Equal(&safeTester{Tester: t}, c.Tester) {
		t.Error()
	}
}

func TestController_RegisterMocker(t *testing.T) {
	c := &Controller{}
	c.RegisterMocker(nil)
	if len(c.mockers) != 1 {
		t.Error()
	}
}

type dummyMocker struct {
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
	if !Equal(int32(2), atomic.LoadInt32(&dm.finishCounter)) {
		t.Error()
	}
}

func TestController_Wait(t *testing.T) {
	dm := &dummyMocker{}
	c := &Controller{
		mockers: []Mocker{dm, dm},
	}

	c.Wait(0)
	if !Equal(int32(2), atomic.LoadInt32(&dm.waitCounter)) {
		t.Error()
	}
}

func TestController_WaitConcurrent(t *testing.T) {
	um1 := &unsafeMocker{}
	um2 := &unsafeMocker{}

	c := &Controller{
		Tester:  newSafeTester(&unsafeTester{}),
		mockers: []Mocker{um1, um2},
	}

	um1.tester = c
	um2.tester = c

	c.Wait(0) //shouln't produce data races
}

type unsafeMocker struct {
	Mocker
	tester Tester
}

func (um *unsafeMocker) MinimockWait(time.Duration) {
	um.tester.FailNow()
}

type unsafeTester struct {
	Tester

	finished bool
}

func (u *unsafeTester) FailNow() {
	u.finished = true
}
