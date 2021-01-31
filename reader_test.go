package main

import (
	"testing"
)

func TestCapabilityToParseFiles(t *testing.T) {

	//Testing
	file := openFile("assets/products.xlsx")

	for i := 0; i < 100000000; i++ {
		parseExcelFile(file)
	}

	//It takes 1.460s to parse 1 000 000 files
	// 		   12.134s to parse  10 000 000 files
	//		   123.241s to parse 100 000 000 files
}
