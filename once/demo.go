package main

import (
	"sync"
	"fmt"
)

type Client struct {
	addr string
	sync.Once
}

func f()  {
	fmt.Println("aaa")
}

func newClient(addr string) *Client{
	c := Client{}
	c.addr = addr
	return &c
}

func main() {
	//c1 := newClient("127.0.0.1:6379")
	//c1.Do(f)
	//c1.Do(f)
	var o sync.Once
	o.Do(f)
	o.Do(f)
}
