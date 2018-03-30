package main

import (
	//"fmt"
	//"net"
	//"bufio"

	"fmt"
	//"net"
	"github.com/golang/protobuf/proto"
	"net"
)

func send()  {
	bs := make([]byte, 7)
	bs[1] = 4
	bs[2] = 1
	fmt.Println(bs)
}

var (
	RequestUdp = []byte{0, 4, 1, 0, 0, 0, 0}
	i1 = uint64(1)
	i2 = int(i1)
)


func main() {



	p :=  make([]byte, 2048)
	conn, err := net.Dial("udp", "10.31.63.11:3658")
	fmt.Println(conn.LocalAddr())
	if err != nil {
		fmt.Printf("Some error %v", err)
		return
	}
	_, err = conn.Write(RequestUdp)
	n, err := conn.Read(p)
	if err == nil {
		var a AgentServerStatusResponse
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
