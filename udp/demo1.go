package main

import (
	//"fmt"
	//"net"
	//"bufio"

	"fmt"
	"github.com/golang/protobuf/proto"
	"net"
	"bufio"
)


var (
	requestUdp = []byte{0, 4, 1, 0, 0, 0, 0}
)


func main() {

	conn, err := net.Dial("udp", "10.31.63.11:3658")
	defer conn.Close()
	reader := bufio.NewReaderSize(conn, 8192)
	_, err = conn.Write(requestUdp)
	if err != nil {
		fmt.Printf("Some error %v", err)
		return
	}
	fmt.Println(conn.LocalAddr())
	resp :=  make([]byte, 7)
	n, err := reader.Read(resp)
	if err == nil && n == 7 && resp[1] == 5 && resp[2] == 1 {
		fmt.Println(resp[3:])
		var length int
		for _, v := range resp[3:] {
			length = length*256 + int(v)
			fmt.Println(length)
		}
		fmt.Println("length", length)

		var a AgentServerStatusResponse
		buf := make([]byte, length)
		fmt.Println("ready to read")
		n, err := reader.Read(buf)
		fmt.Println("read", n)
		if err != nil {
			fmt.Println("read error", err)
			return
		}
		if n != length {
			fmt.Println("length error", length, n)
		}
		err = proto.Unmarshal(buf, &a)
		if err != nil {
			fmt.Println(string(buf))
			fmt.Println("unmarshaling error: ", err)
		} else {
			fmt.Println(a.GetZkHosts(), *a.ProcessCount, *a.QueueEleCount)
		}
	} else {
		fmt.Println("response error", err, string(resp))
	}

}
