package main

import (
	"fmt"
	"sync"
	"time"
)

var (
	mutex sync.RWMutex
	wait sync.WaitGroup
	number int
)


func f1() {
	go func() {
		mutex.RLock()
		defer wait.Done()
		defer mutex.RUnlock()
		fmt.Println(number)
		time.Sleep(time.Second * 3)
		fmt.Println(number)
	}()
}

func f2()  {
	go func() {
		time.Sleep(time.Second * 1)
		mutex.Lock()
		defer wait.Done()
		defer mutex.Unlock()
		number = 100
	}()
}

func main() {
	wait.Add(2)
	f1()
	f2()
	wait.Wait()
	fmt.Println(number)
}