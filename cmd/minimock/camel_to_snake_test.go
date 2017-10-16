package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_camelToSnake(t *testing.T) {
	assert.Equal(t, "i_love_golang_and_json_so_much", camelToSnake("ILoveGolangAndJSONSoMuch"))
	assert.Equal(t, "i_love_json", camelToSnake("ILoveJSON"))
	assert.Equal(t, "json", camelToSnake("json"))
	assert.Equal(t, "json", camelToSnake("JSON"))
	assert.Equal(t, "привет_мир", camelToSnake("ПриветМир"))
}

func Benchmark_CamelToSnake(b *testing.B) {
	for n := 0; n < b.N; n++ {
		camelToSnake("TestTable")
	}
}
