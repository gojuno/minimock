package tests

//Tester interface contains methods of testing.T that are used by code generated with minimock
type Tester interface {
	Fatal(args ...interface{})
	Fatalf(format string, args ...interface{})
	Error(...interface{})
}
