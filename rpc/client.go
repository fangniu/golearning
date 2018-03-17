package main

import (
	"os"
	"fmt"
	"log"
	"net/rpc"
	"time"
)

type Args struct {
	A, B int
}

type Quotient struct {
	Quo, Rem int
}


func main() {

	if len(os.Args) != 2 {
		fmt.Println("Usage: ", os.Args[0], "server:port")
		os.Exit(1)
	}
	service := os.Args[1]
	client, err := rpc.Dial("tcp", service)
	now := time.Now()
	if err != nil {
		log.Fatalln("dialing error", err)
	}
	args := &Args{7, 8}
	var replay int
	err = client.Call("Arith.Multiply", args, &replay)
	if err != nil {
		log.Fatalln("arith error:", err)
	}
	fmt.Printf("Arith: %d*%d=%d\n", args.A, args.B, replay)

	var quot Quotient
	err = client.Call("Arith.Divide", args, &quot)
	if err != nil {
		log.Fatalln("arith error:", err)
	}
	fmt.Printf("Arith: %d/%d=%d remainder %d\n", args.A, args.B, quot.Quo, quot.Rem)
	fmt.Println(time.Now().Sub(now))
}