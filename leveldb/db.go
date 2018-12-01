package leveldb

import "fmt"

var A int


func init() {
	A := 1
	fmt.Printf("Hello world! %d \n", A)
}

func Test() {
	fmt.Printf("Hello world! %d \n", A)
}