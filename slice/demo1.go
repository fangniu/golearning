package main

import "fmt"

type myType struct {
	s []string
}

func main() {
	t := myType{}
	t.s = []string{"aa"}
	t.s = nil
	var s1 *string
	s1 = nil
	s2 := *s1
	t.s = append(t.s, s2)
	fmt.Println(t.s)
}
