package exporter

import (
	pb "proto"
	"fmt"
	"google.golang.org/grpc"
	"golang.org/x/net/context"
	"log"
)

var address string
var ResourceInfo ServerStatItem

func Report (warning bool, serverState, address string) (bool) {
	pioResourceInfo := &ResourceInfo
	pioResourceInfo = new(ServerStatItem)
	if warning == true {
	    pioResourceInfo.Tag = false
	    name := serverState
	    conn, err := grpc.Dial(address, grpc.WithInsecure())
	    if err != nil {
	        fmt.Println(err)
	    }
	    defer conn.Close()
	    client := pb.NewHelloClient(conn)
	    request, err := client.SayHello(context.Background(), &pb.HelloRequest{Name: name})
	    if err != nil {
	        log.Fatal(err)
	    }
	    fmt.Println(request.Message)
	} else {
        pioResourceInfo.Tag = true
        name := serverState
        conn, err := grpc.Dial(address, grpc.WithInsecure())
        if err != nil {
            log.Fatal(err)
        }
        defer conn.Close()
        client := pb.NewHelloClient(conn)
        request, err := client.SayHello(context.Background(), &pb.HelloRequest{Name: name})
        if err != nil {
            log.Fatal(err)
        }
        fmt.Println(request.Message)
	}
	return pioResourceInfo.Tag
}
