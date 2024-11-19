package main

import (
	proto "Auction/grpc"
	"bufio"
	"context"
	"log"
	"net"
	"os"
	"strconv"
	"strings"

	"google.golang.org/grpc"
)

type AuctionServiceServer struct {
	proto.UnimplementedAuctionServiceServer
	winingBid       proto.Bid
	auctionFinished bool
}

func main() {
	srv := &AuctionServiceServer{
		winingBid:       proto.Bid{Amount: 0, Node: "Server"},
		auctionFinished: false,
	}
	var input string
	reader := bufio.NewReader(os.Stdin)
	read, _ := reader.ReadString('\n')
	input, _, _ = strings.Cut(read, "\r\n")
	go srv.startServer(input)

	for {

	}
}

func (srv *AuctionServiceServer) startServer(address string) {
	grpcServer := grpc.NewServer()
	listener, err := net.Listen("tcp", ":"+address)
	if err != nil {
		log.Fatalln("Exception Error")
	}
	proto.RegisterAuctionServiceServer(grpcServer, srv)
	log.Println(address + " registered server")
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Exception Error after Registration")
	}
}

func (srv *AuctionServiceServer) SendBid(ctx context.Context, bid *proto.Bid) (*proto.Acknowledgement, error) {

	log.Println("Bidder " + bid.Node + " bidded " + strconv.Itoa(int(bid.Amount)))
	if bid.Amount > srv.winingBid.Amount {
		srv.winingBid = *bid
		return &proto.Acknowledgement{
			Status: proto.Status_SUCCESS,
		}, nil
	}
	return &proto.Acknowledgement{
		Status: proto.Status_FAIL,
	}, nil
}

func (srv *AuctionServiceServer) Result(ctx context.Context, _ *proto.Empty) (*proto.Outcome, error) {

	return &proto.Outcome{
		Winingbid:       &srv.winingBid,
		AuctionFinished: srv.auctionFinished,
	}, nil
}
