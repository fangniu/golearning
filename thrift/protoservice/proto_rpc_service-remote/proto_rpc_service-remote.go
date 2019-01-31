// Autogenerated by Thrift Compiler (1.0.0-dev)
// DO NOT EDIT UNLESS YOU ARE SURE THAT YOU KNOW WHAT YOU ARE DOING

package main

import (
	"flag"
	"fmt"
	"git.apache.org/thrift.git/lib/go/thrift"
	"math"
	"net"
	"net/url"
	"os"
	"github.com/golearning/thrift/protoservice"
	"strconv"
	"strings"
)

func Usage() {
	fmt.Fprintln(os.Stderr, "Usage of ", os.Args[0], " [-h host:port] [-u url] [-f[ramed]] function [arg1 [arg2...]]:")
	flag.PrintDefaults()
	fmt.Fprintln(os.Stderr, "\nFunctions:")
	fmt.Fprintln(os.Stderr, "  ProtoReply dealTwowayMessage(ProtoRequest msg)")
	fmt.Fprintln(os.Stderr, "  void dealOnewayMessage(ProtoRequest msg)")
	fmt.Fprintln(os.Stderr)
	os.Exit(0)
}

func main() {
	flag.Usage = Usage
	var host string
	var port int
	var protocol string
	var urlString string
	var framed bool
	var useHttp bool
	var parsedUrl url.URL
	var trans thrift.TTransport
	_ = strconv.Atoi
	_ = math.Abs
	flag.Usage = Usage
	flag.StringVar(&host, "h", "localhost", "Specify host and port")
	flag.IntVar(&port, "p", 9090, "Specify port")
	flag.StringVar(&protocol, "P", "binary", "Specify the protocol (binary, compact, simplejson, json)")
	flag.StringVar(&urlString, "u", "", "Specify the url")
	flag.BoolVar(&framed, "framed", false, "Use framed transport")
	flag.BoolVar(&useHttp, "http", false, "Use http")
	flag.Parse()

	if len(urlString) > 0 {
		parsedUrl, err := url.Parse(urlString)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error parsing URL: ", err)
			flag.Usage()
		}
		host = parsedUrl.Host
		useHttp = len(parsedUrl.Scheme) <= 0 || parsedUrl.Scheme == "http"
	} else if useHttp {
		_, err := url.Parse(fmt.Sprint("http://", host, ":", port))
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error parsing URL: ", err)
			flag.Usage()
		}
	}

	cmd := flag.Arg(0)
	var err error
	if useHttp {
		trans, err = thrift.NewTHttpClient(parsedUrl.String())
	} else {
		portStr := fmt.Sprint(port)
		if strings.Contains(host, ":") {
			host, portStr, err = net.SplitHostPort(host)
			if err != nil {
				fmt.Fprintln(os.Stderr, "error with host:", err)
				os.Exit(1)
			}
		}
		trans, err = thrift.NewTSocket(net.JoinHostPort(host, portStr))
		if err != nil {
			fmt.Fprintln(os.Stderr, "error resolving address:", err)
			os.Exit(1)
		}
		if framed {
			trans = thrift.NewTFramedTransport(trans)
		}
	}
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error creating transport", err)
		os.Exit(1)
	}
	defer trans.Close()
	var protocolFactory thrift.TProtocolFactory
	switch protocol {
	case "compact":
		protocolFactory = thrift.NewTCompactProtocolFactory()
		break
	case "simplejson":
		protocolFactory = thrift.NewTSimpleJSONProtocolFactory()
		break
	case "json":
		protocolFactory = thrift.NewTJSONProtocolFactory()
		break
	case "binary", "":
		protocolFactory = thrift.NewTBinaryProtocolFactoryDefault()
		break
	default:
		fmt.Fprintln(os.Stderr, "Invalid protocol specified: ", protocol)
		Usage()
		os.Exit(1)
	}
	client := protoservice.NewProtoRpcServiceClientFactory(trans, protocolFactory)
	if err := trans.Open(); err != nil {
		fmt.Fprintln(os.Stderr, "Error opening socket to ", host, ":", port, " ", err)
		os.Exit(1)
	}

	switch cmd {
	case "dealTwowayMessage":
		if flag.NArg()-1 != 1 {
			fmt.Fprintln(os.Stderr, "DealTwowayMessage requires 1 args")
			flag.Usage()
		}
		arg4 := flag.Arg(1)
		mbTrans5 := thrift.NewTMemoryBufferLen(len(arg4))
		defer mbTrans5.Close()
		_, err6 := mbTrans5.WriteString(arg4)
		if err6 != nil {
			fmt.Println(err6)
			Usage()
			return
		}
		factory7 := thrift.NewTSimpleJSONProtocolFactory()
		jsProt8 := factory7.GetProtocol(mbTrans5)
		argvalue0 := protoservice.NewProtoRequest()
		err9 := argvalue0.Read(jsProt8)
		if err9 != nil {
			fmt.Println(err9)
			Usage()
			return
		}
		value0 := argvalue0
		fmt.Print(client.DealTwowayMessage(value0))
		fmt.Print("\n")
		break
	case "dealOnewayMessage":
		if flag.NArg()-1 != 1 {
			fmt.Fprintln(os.Stderr, "DealOnewayMessage requires 1 args")
			flag.Usage()
		}
		arg10 := flag.Arg(1)
		mbTrans11 := thrift.NewTMemoryBufferLen(len(arg10))
		defer mbTrans11.Close()
		_, err12 := mbTrans11.WriteString(arg10)
		if err12 != nil {
			Usage()
			return
		}
		factory13 := thrift.NewTSimpleJSONProtocolFactory()
		jsProt14 := factory13.GetProtocol(mbTrans11)
		argvalue0 := protoservice.NewProtoRequest()
		err15 := argvalue0.Read(jsProt14)
		if err15 != nil {
			fmt.Println(err15)
			Usage()
			return
		}
		value0 := argvalue0
		fmt.Print(client.DealOnewayMessage(value0))
		fmt.Print("\n")
		break
	case "":
		Usage()
		break
	default:
		fmt.Fprintln(os.Stderr, "Invalid function ", cmd)
	}
}