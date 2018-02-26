package main

import "fmt"

func main() {
	origin, wait := make(chan int, 0), make(chan struct{}, 0)
	Processor(origin, wait)
	for num := 2; num < 1000; num++ {
		origin <- num
	}
	close(origin)
	<-wait
	fmt.Println("end.")
}

func Processor(seq chan int, wait chan struct{}) {
	go func() {
		prime, ok := <-seq
		if !ok {
			close(wait)
			return
		}
		/*prime := <-seq
		defer close(wait)*/
		fmt.Println(prime)
		out := make(chan int)
		Processor(out, wait) // 此处为什么要在递归调用一次呢
		for num := range seq {
			if num%prime != 0 {
				out <- num
			}
		}
		close(out)
	}()
}