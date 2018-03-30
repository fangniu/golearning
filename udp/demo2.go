package main

import (
	//"fmt"
	//"net"
	//"bufio"

	"fmt"
	//"net"
	"github.com/golang/protobuf/proto"
)

func send()  {
	bs := make([]byte, 7)
	bs[1] = 4
	bs[2] = 1
	fmt.Println(bs)
}

var (
	RequestUdp = []byte{0, 4, 1, 0, 0, 0, 0}
)


func main() {



	//p :=  make([]byte, 2048)
	//conn, err := net.Dial("udp", "10.31.63.11:3658")
	//fmt.Println(conn.LocalAddr())
	//if err != nil {
	//	fmt.Printf("Some error %v", err)
	//	return
	//}
	//_, err = conn.Write(RequestUdp)
	//n, err := conn.Read(p)
	if true {
		sli := "192.168.101.58"
		port := uint32(3658)
		ci := uint32(30)
		zkHosts := "zkhosts1,2,3"
		zkAgentPath := "ZkAgentPath"
		zkProxyPath := "ZkProxyPath"
		zkCheckInterval := uint32(30)
		queueName := "QueueName"
		queueEleSize := uint32(30)
		queueEleCount := uint32(30)
		processCount := uint64(30)
		processPerSec := uint64(30)
		proxySize := uint32(30)
		attrSize := uint32(30)

		a := AgentServerStatusResponse{
			StrLocalIp: &sli,
			LocalPort:        &port,
			CollectInterval:  &ci,
			ZkHosts:          &zkHosts,
			ZkAgentPath:      &zkAgentPath,
			ZkProxyPath:     &zkProxyPath,
			ZkCheckInterval:  &zkCheckInterval,
			QueueName:        &queueName,
			QueueEleSize:     &queueEleSize,
			QueueEleCount:    &queueEleCount,
			ProcessCount:     &processCount,
			ProcessPerSec:    &processPerSec,
			ProxySize:        &proxySize,
			AttrSize:         &attrSize,
		}
		data, err := proto.Marshal(&a)
		if err != nil {
			fmt.Println("marshaling error: ", err)
			return
		}

		newData := &AgentServerStatusResponse{}
		//err = proto.Unmarshal(p[:n], &a)
		err = proto.Unmarshal(data, newData)
		if err != nil {
			fmt.Println("unmarshaling error: ", err)
		} else {
			fmt.Println(a.AttrSize, *a.ProcessCount, a.QueueEleCount)
		}

	} else {
		fmt.Printf("Some error %v\n", "aa")
	}
	//conn.Close()
}
