package main

import (
	"math/bits"
	"fmt"
	"github.com/gomodule/redigo/redis"
)

type AAA struct {
	redis.ConnWithTimeout
}

func main() {
	var aaa AAA
	aaa.Send("PING")

	var a uint = 31
	fmt.Printf("bits.OnesCount(%d) = %d\n", a, bits.OnesCount(a))

	a++
	a++
	fmt.Printf("bits.OnesCount(%d) = %d\n", a, bits.OnesCount(a))
}