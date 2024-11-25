package main

import (
	proto "Auction/grpc"
	"bufio"
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
	"time"

	"google.golang.org/grpc"
)

type AuctionServiceServer struct {
	proto.UnimplementedAuctionServiceServer
	winingBid         proto.Bid
	auctionFinished   bool
	isLeader          bool
	currentLeader     proto.AuctionServiceClient
	servernodes       []proto.AuctionServiceClient
	servernodesString []string
}

func main() {
	srv := &AuctionServiceServer{
		winingBid:         proto.Bid{Amount: 0, Node: "Server"},
		auctionFinished:   false,
		isLeader:          false,
		servernodesString: []string{"5050", "5051", "5052", "5053"},
	}
	var input string
	reader := bufio.NewReader(os.Stdin)
	read, _ := reader.ReadString('\n')
	input, _, _ = strings.Cut(read, "\r\n")
	go srv.startServer(input)
	go srv.healthcheck()

	for {

	}
}

func (srv *AuctionServiceServer) healthcheck() {
	for {
		time.Sleep(1 * time.Second)
		if !srv.isLeader {
			_, err := srv.currentLeader.SendBid(context.Background(), &proto.Bid{Node: "Server 5050", Amount: 0})
			if err != nil {
				srv.selectleader()
			}
		}
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
		log.Println("Bidder bidded more than the previous winning bid " + fmt.Sprint(srv.winingBid.Amount))
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

func (srv *AuctionServiceServer) selectleader() {
	for _, elem := range srv.servernodes {
		_, err := elem.SendBid(context.Background(), &proto.Bid{Node: "Server", Amount: 0})
	}
}
