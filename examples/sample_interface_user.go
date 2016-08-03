package examples

import "github.com/gojuno/minimock/examples/sample"

func SampleInterfaceUser(si sample.Interface) []*string {
	s := si.GetString()
	sum := si.CalculateSum(1, 2, 3)

	return si.GetArrayOfStrings(s, sum)
}
