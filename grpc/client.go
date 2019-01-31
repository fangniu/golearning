package main


import (
	"log"
	"os"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"github.com/golearning/grpc/hello"
	"fmt"
	"time"
	"encoding/base64"
)

const (
	address     = "localhost:50051"
	defaultName = "World"
)

func main() {
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatal("did not connect:", err)
	}
	defer conn.Close()
	c := hello.NewHelloServiceClient(conn)

	name := defaultName
	if len(os.Args) >1 {
		name = os.Args[1]
	}
	fmt.Println(time.Now())
	h := hello.HelloRequest{Greeting: name}
	data := base64.StdEncoding.EncodeToString([]byte(h.String()))
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	fmt.Println("aaa", data)
	r, err := c.SayHello(context.Background(), &h)
	fmt.Println(time.Now())
	if err != nil {
		log.Fatal("could not greet:", err)
	}
	log.Printf("Greeting: %s, %d", r.Reply, r.Number[2])
}
