package main

import (
	"flag"
			"github.com/golearning/thrift/protoservice"
	"fmt"
	"log"
	"git.apache.org/thrift.git/lib/go/thrift"
			"errors"
)

const (
	ResponseType  int32 = 199999
	RequestType int32 = 199998
)

var (
	host *string
	port *int
	shardingNum *int64
	retry *int
	client *protoservice.ProtoRpcServiceClient
)

func check(id int64) error {
	loRequest := protoservice.NewProtoRequest()
	loRequest.Type = RequestType
	loRequest.ShardingID = id
	fmt.Println(loRequest.String())

	r, err := client.DealTwowayMessage(loRequest)
	if err != nil {
		return err
	}
	fmt.Println(r.String())
	if r.Type != ResponseType {
		return errors.New(fmt.Sprintf("response type error: %d", r.Type))
	}
	return nil
}

func main()  {
	host = flag.String("h", "127.0.0.1", "host ip")
	port = flag.Int("p", 123, "host port")
	shardingNum = flag.Int64("s", 4, "sharding Number")
	retry = flag.Int("r", 3, "check retry")

	flag.Parse()
	addr := fmt.Sprintf("%s:%d", *host, *port)
	trans, err := thrift.NewTSocket(addr)
	if err != nil {
		log.Fatalln("[ERROR] tSocket:", err)
	}

	protocolFactory := thrift.NewTBinaryProtocolFactoryDefault()
	client = protoservice.NewProtoRpcServiceClientFactory(trans, protocolFactory)

	if err := trans.Open(); err != nil {
		log.Fatalln("[ERROR] opening:", addr)
	}
	defer trans.Close()

	//for i := int64(1); i <= *shardingNum; i ++ {
	//	var err error
	//	var count int
	//	for j := 0; j < *retry; j ++ {
	//		if e := check(i); e == nil {
	//			break
	//		} else {
	//			count ++
	//			err = e
	//		}
	//	}
	//	if count == *retry {
	//		log.Fatalln(fmt.Sprintf("[ERROR] sharding[%d]:", i), err)
	//	} else {
	//		log.Println(fmt.Sprintf("[INFO] sharding[%d]: OK", i))
	//	}
	//}
	check(*shardingNum)
}
