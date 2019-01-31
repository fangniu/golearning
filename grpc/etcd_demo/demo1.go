package main

import (
		"time"
		"context"
		"log"
		"go.etcd.io/etcd/clientv3"
		"fmt"
)

func main() {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"192.168.8.65:22379"},
		DialTimeout: 5 * time.Second,
		Username: "cs1",
		Password: "123456",
	})
	if err != nil {
		log.Fatalln(err)
	}
	defer cli.Close()
	leaseApi := clientv3.NewLease(cli)
	resp, err := leaseApi.Grant(context.TODO(), 10)
	if err != nil {
		log.Fatalln(err)
	}
	leaseID := resp.ID
	go func(leaseID clientv3.LeaseID) {
		ch, err := leaseApi.KeepAlive(context.TODO(), leaseID)
		if err != nil {
			log.Fatalln(err)
		}
		for lease := range ch {
			log.Println(lease.ID, lease.TTL)
		}
	}(leaseID)

	//ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	rch := cli.Watch(context.Background(), "e3w_test/services/", clientv3.WithPrefix())
	log.Println("开始监控...")
	for wresp := range rch {
		for _, ev := range wresp.Events {
			log.Println(fmt.Sprintf("%s %q : %q", ev.Type, ev.Kv.Key, ev.Kv.Value))
		}
	}
	//cancel()
	//if err != nil {
	//	// handle error!
	//	switch err {
	//	case context.Canceled:
	//		log.Fatalf("ctx is canceled by another routine: %v", err)
	//	case context.DeadlineExceeded:
	//		log.Fatalf("ctx is attached with a deadline is exceeded: %v", err)
	//	case rpctypes.ErrEmptyKey:
	//		log.Fatalf("client-side error: %v", err)
	//	default:
	//		log.Fatalf("bad cluster endpoints, which are not etcd servers: %v", err)
	//	}
	//}
	// use the response
	//for _, ev := range resp.Kvs {
	//	fmt.Printf("%s : %s\n", ev.Key, ev.Value)
	//}
}