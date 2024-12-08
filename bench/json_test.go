package main

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/hxkhan/jsonx"
)

const fileN = "./input.json"

var file, err = os.ReadFile(fileN)

func BenchmarkCustom(b *testing.B) {
	if err != nil {
		panic(err)
	}
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		CustomDecode()
	}
}

func BenchmarkStandard(b *testing.B) {
	if err != nil {
		panic(err)
	}
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		StandardDecode()
	}
}

func CustomDecode() interface{} {
	structure, err := jsonx.Decode(file)
	if err != nil {
		panic(err)
	}
	return structure
}

func StandardDecode() interface{} {
	var structure interface{}
	err = json.Unmarshal(file, &structure)
	if err != nil {
		panic(err)
	}
	return structure
}
