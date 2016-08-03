package sample

type Interface interface {
	GetString() string
	CalculateSum(ints ...int) int
	GetArrayOfStrings(s string, i int) []*string
}
