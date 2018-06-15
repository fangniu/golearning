package main

import (
	"fmt"
	"time"
)

func main() {
	ch := make(chan int, 10000)

	go func() {
		for i:=0;i<10000;i++ {
			ch <- 1
		}
	}()
	<- ch
	time.Sleep(time.Microsecond*10)
	fmt.Println(len(ch))
}