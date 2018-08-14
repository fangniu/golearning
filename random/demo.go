package main

import (
	"math/rand"
	"fmt"
	"time"
	"strconv"
)

func main() {
	//rand.Seed(time.Now().UnixNano())
	//fmt.Println(rand.Perm(6))
	timestamp := time.Now().Unix()
	rand.Seed(timestamp)
	fmt.Println(strconv.Itoa(rand.Int()))
}
