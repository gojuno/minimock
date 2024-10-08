package minimock

import "sync"

type safeTester struct {
	Tester
	m sync.Mutex
}

func newSafeTester(t Tester) *safeTester {
	return &safeTester{Tester: t}
}

// Error implements Tester
func (st *safeTester) Error(args ...interface{}) {
	st.m.Lock()
	defer st.m.Unlock()
	st.Tester.Helper()

	st.Tester.Error(args...)
}

// Errorf implements Tester
func (st *safeTester) Errorf(format string, args ...interface{}) {
	st.m.Lock()
	defer st.m.Unlock()
	st.Tester.Helper()

	st.Tester.Errorf(format, args...)
}

// Fatal implements Tester
func (st *safeTester) Fatal(args ...interface{}) {
	st.m.Lock()
	defer st.m.Unlock()
	st.Tester.Helper()

	st.Tester.Fatal(args...)
}

// Fatalf implements Tester
func (st *safeTester) Fatalf(format string, args ...interface{}) {
	st.m.Lock()
	defer st.m.Unlock()
	st.Tester.Helper()

	st.Tester.Fatalf(format, args...)
}
