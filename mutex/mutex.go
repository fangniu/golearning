package main

import (
	"sync"
	"fmt"
	"time"
)

var m sync.Mutex
var wg sync.WaitGroup

func f3()  {
	time.Sleep(time.Second*1)
	m.Lock()
	defer wg.Done()
	defer m.Unlock()
	fmt.Println("aaa")
}

func f4()  {
	time.Sleep(time.Second*2)
	m.Lock()
	defer wg.Done()
	defer m.Unlock()
	fmt.Println("bbb")
}

func f5()  {
	time.Sleep(time.Second*2)
	m.Lock()
	defer wg.Done()
	defer m.Unlock()
	fmt.Println("ccc")
}

func main() {
	wg.Add(3)
	go f3()
	go f5()
	go f4()
	wg.Wait()
}