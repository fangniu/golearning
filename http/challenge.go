package main

import (
	"net"
	"fmt"
	"encoding/binary"
	"golang.org/x/text/encoding/simplifiedchinese"
)

func f1()  {
	conn, err := net.Dial("tcp", "challenge.yuansuan.cn:7042")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer conn.Close()
	buff := make([]byte, 1024)
	l, err := conn.Read(buff)
	if err != nil {
		fmt.Println(err)
		return
	}
	if l > 32 {
		fmt.Println("length error", l)
		return
	}
	id := string(buff[18:(l-1)])
	content := fmt.Sprintf("IAM:%v:%v\n", id, "838307912@qq.com")
	fmt.Println(content)
	l, err = conn.Write([]byte(content))
	if err != nil {
		fmt.Println(err)
		return
	}
	l, err = conn.Read(buff)
	if err != nil {
		fmt.Println(err)
		return
	}
	if string(buff[:l]) == "SUCCESS!\n" {
		fmt.Println("connect success!")
	}
	l, err = conn.Read(buff)
	fmt.Println("length", l)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(binary.LittleEndian.Uint32(buff[:4]))
	//fmt.Println(string(buff[4:8]))
	fmt.Println(binary.LittleEndian.Uint32(buff[8:12]))
	decodeBytes, _ := simplifiedchinese.GB18030.NewDecoder().Bytes(buff[12:26])
	fmt.Println(string(decodeBytes))
	//fmt.Println("length", l)
	//fmt.Println(string(buff))

}




func main() {
	f1()
}