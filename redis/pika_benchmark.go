package main

import (
	"math/rand"
	"github.com/gomodule/redigo/redis"
	"time"
	"log"
	"flag"
	"runtime"
	"sync"
)

const randBytes = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

var wg sync.WaitGroup

type Client struct {
	redis redis.Conn
}

func NewClient(addr string) Client {
	conn, err := redis.Dial("tcp", addr)
	if err != nil {
		log.Fatalln(err)
	}
	return Client{conn}
}

func (c *Client) run() {
	defer wg.Done()
	for {
		c.set()
		c.hset()
	}
}

func (c *Client) hset() {
	c.redis.Do("HSET", genRandString(16), "hello", "world")
}

func (c *Client) set() {
	_, err := c.redis.Do("SET", genRandString(16), genRandString(16))
	if err != nil {
		log.Fatalln("SET ERROR", err)
		return
	}
}


func genRandString(n int) string {
	rand.Seed(time.Now().UnixNano())
	b := make([]byte, n)
	for i := range b {
		b[i] = randBytes[rand.Intn(len(randBytes))]
	}
	return string(b)
}

func main() {
	n := flag.Int("cpuNum", 2, "runtime.GOMAXPROCS")
	addr := flag.String("addr", "127.0.0.1:6379", "redis address")
	flag.Parse()
	runtime.GOMAXPROCS(*n)
	for i:=0; i < *n; i ++ {
		wg.Add(1)
		c := NewClient(*addr)
		go c.run()
	}
	wg.Wait()
}