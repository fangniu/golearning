package main

import (
	"math/rand"
	"fmt"
	"time"
)

func main() {
	rand.Seed(time.Now().UnixNano())
	fmt.Println(rand.Perm(6))
}
