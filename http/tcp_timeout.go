package main

import (
	"net"
	"bufio"
	"time"
	"fmt"
)

func tcp1() (err error) {
	var conn net.Conn
	conn, err = net.DialTimeout("tcp", "127.0.0.1:6379", time.Millisecond * 100)
	if err != nil {
		return
	}
	conn.Close()
	defer conn.Close()
	reader := bufio.NewReaderSize(conn, 8192)
	//_, err = conn.Write([]byte("sdfsdfsdfsdflkjdflkgjldskfjg"))
	if err != nil {
		return
	}
	resp :=  make([]byte, 7)
	var n int
	fmt.Println("aaa")
	n, err = reader.Read(resp)
	fmt.Println("bbb")
	fmt.Println(n)
	return
}

func main() {
	err := tcp1()
	fmt.Println(err)
}