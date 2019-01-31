package main

import (
	"time"
	"log"
	"go.etcd.io/etcd/clientv3"
	"context"
	"fmt"
	)

func main() {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"192.168.8.64:2379", "192.168.8.65:22379", "192.168.100.7:32379"},
		DialTimeout: 5 * time.Second,
		Username: "root",
		Password: "123456",
	})
	if err != nil {
		log.Fatalln(err)
	}
	defer cli.Close()
	resp, err := cli.Put(context.Background(), "/foo", "111")
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(resp.Header.String())
}