package client

import (
	proto "Auction/grpc"
	"context"
	"log"
	"net"
	"slices"
	"strconv"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type ClientServiceServer struct {
	proto.UnimplementedClientServiceServer
}

func (srv *ClientServiceServer) startServer(address string) {
	grpcServer := grpc.NewServer()
	listener, err := net.Listen("tcp", ":"+address)
	if err != nil {
		log.Fatalln("Exception Error")
	}
	proto.RegisterClientServiceServer(grpcServer, srv)
	log.Println(address + " registered server")
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Exception Error after Registration")
	}
}

func StartNode(address string, neighbor []string) {
	srv := &ClientServiceServer{}
	go srv.startServer(address)
	time.Sleep(4 * time.Second)
	index := slices.Index(neighbor, address) + 1
	if index == len(neighbor) {
		index = 0
	}
	conn, err := grpc.NewClient("localhost:"+neighbor[index], grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Exception Error")
	}
	client := proto.NewClientServiceClient(conn)
	client.SendBid(context.Background(), &proto.Bid{Node: address, Amount: 0})
}

func (srv *ClientServiceServer) SendBid(ctx context.Context, bid *proto.Bid) (*proto.Acknowledgement, error) {

	log.Println("Bidder " + bid.Node + " bidded " + strconv.Itoa(int(bid.Amount)))
	return &proto.Acknowledgement{
		Status: proto.Status_SUCCESS,
	}, nil
}
