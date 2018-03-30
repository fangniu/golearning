package main

import (
	"time"
	"fmt"
	"sync"
)

var (
	ch2 = make(chan int, 1000)
	wg2 sync.WaitGroup
)

func push() {
	defer wg2.Done()
	time.Sleep(time.Second * 1)
	count := 1000
	for count > 0 {
		ch2 <- 1
		count --
	}
}

func pop()  {
	defer wg2.Done()
	for {
		_, ok := <- ch2
		if !ok {
			return
		}
		fmt.Println(len(ch2))
		time.Sleep(time.Second*1)
	}
}

func main() {
	wg2.Add(2)
	go push()
	go pop()
	wg2.Wait()
}