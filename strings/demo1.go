package main

import (
	"strings"
	"fmt"
	"reflect"
)

func f1()  {
	s := "aaa.bbb"
	s1 := strings.Replace(s,".", `"."`, 1)
	fmt.Println(s1)
}

type TypeA struct {
	i int
}

func newTypeA() *TypeA {
	return &TypeA{}
}

func main() {
	//var t *TypeA
	t := newTypeA()

	fmt.Println(reflect.TypeOf(t))
	fmt.Println(t.i)


}