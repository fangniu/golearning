package main

import (
	"fmt"
	"time"
	"sync"
)


var wg sync.WaitGroup

var ch = make(chan int, 11)

func f1(i int)  {
	defer wg.Done()
	time.Sleep(time.Second * 1)
	ch <- i
}

func f2(i int)  {
	defer wg.Done()
	time.Sleep(time.Second * 1)
	ch <- i
}

func main() {
	wg.Add(10)
	go f1(1)
	go f2(2)
	go f1(3)
	go f2(4)
	go f1(5)
	go f2(6)
	go f1(7)
	go f2(8)
	go f1(9)
	go f2(10)
	wg.Wait()
	close(ch)
	for i := range ch {
		fmt.Println(i)
	}
}