package main

import (
	"fmt"
	"time"
	"sync"
)

var (
	ch = make(chan int, 1000)
	wg sync.WaitGroup
)

func f1()  {
	defer wg.Done()
	for {
		select {
		default:
			fmt.Println("aaa")
			value, ok := <-ch
			fmt.Println("bbb")

			fmt.Println(ok, value)
		}
	}
}

func f2()  {
	defer wg.Done()
	for {
		ch <- 1
		time.Sleep(time.Second * 2)
	}

}

func main() {
	wg.Add(2)
	go f1()
	go f2()
	wg.Wait()
}