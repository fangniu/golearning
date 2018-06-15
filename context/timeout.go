package main

import (
	"context"
	"log"
	"os"
	"time"
	"runtime"
	"fmt"
)

var logg *log.Logger

func someHandler() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*1)
	go doStuff(ctx)
	fmt.Println(runtime.NumGoroutine())
	//10秒后取消doStuff
	time.Sleep(5 * time.Second)
	cancel()

}

//每1秒work一下，同时会判断ctx是否被取消了，如果是就退出
func doStuff(ctx context.Context) {
	select {
	case <-ctx.Done():
		logg.Printf("timeout")
		return
	case <-getData():
		logg.Printf("work")
	}
}

func getData() <-chan struct{} {
	ch := make(chan struct{}, 1)
	ch <- struct{}{}
	go func() {
		time.Sleep(time.Second*1)
	}()
	return ch
}

func main() {
	logg = log.New(os.Stdout, "", log.Ltime)
	someHandler()
	logg.Printf("down")
}