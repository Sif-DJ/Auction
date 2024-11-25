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
	"google.golang.org/grpc/credentials/insecure"
)

type AuctionServiceServer struct {
	proto.UnimplementedAuctionServiceServer
	winingBid         proto.Bid
	auctionFinished   bool
	isLeader          bool
	currentLeader     proto.AuctionServiceClient
	servernodes       []proto.AuctionServiceClient
	servernodesString []string
	serverIndex       int
	serverTime        int
}

func main() {
	srv := &AuctionServiceServer{
		winingBid:         proto.Bid{Amount: 0, Node: "Server"},
		auctionFinished:   false,
		isLeader:          false,
		servernodesString: []string{"5050", "5051", "5052", "5053"},
		servernodes:       make([]proto.AuctionServiceClient, 4),
		serverTime:        0,
	}
	var input string
	reader := bufio.NewReader(os.Stdin)
	read, _ := reader.ReadString('\n')
	input, _, _ = strings.Cut(read, "\r\n")
	for i, elem := range srv.servernodesString {
		if elem == input {
			srv.serverIndex = i
		}
	}
	go srv.startServer(input)
	srv.ConnectToNodes()
	go srv.healthcheck()

	for srv.serverTime < 100 { // 1 minute and 40 seconds to bid before closing of auction
		srv.serverTime++
		time.Sleep(1 * time.Second)
	}

	srv.auctionFinished = true
	for input != "EXIT" {
		read, _ = reader.ReadString('\n')
		input, _, _ = strings.Cut(read, "\r\n")
	}
}

func (srv *AuctionServiceServer) ConnectToNodes() {
	for i, elem := range srv.servernodesString {
		if srv.serverIndex == i {
			srv.servernodes[i] = nil
			continue
		}
		conn, err := grpc.NewClient("localhost:"+elem, grpc.WithTransportCredentials(insecure.NewCredentials()))
		for err != nil {
			conn, err = grpc.NewClient("localhost:"+elem, grpc.WithTransportCredentials(insecure.NewCredentials()))
		}
		log.Println("Connected to " + elem)
		srv.servernodes[i] = (proto.NewAuctionServiceClient(conn))
	}

	srv.selectleader()
}

// Local health check of other servers
func (srv *AuctionServiceServer) healthcheck() {
	for {
		time.Sleep(1 * time.Second)
		if !srv.isLeader {
			outcome, err := srv.currentLeader.HealthCheck(context.Background(), &proto.Empty{})
			if err != nil {
				srv.selectleader()
			} else {
				srv.winingBid = *outcome.GetWiningbid()
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
	if srv.isLeader {
		log.Println("Bidder " + bid.Node + " bidded " + strconv.Itoa(int(bid.Amount)))
		if bid.Amount > srv.winingBid.Amount && !srv.auctionFinished {
			log.Println("Bidder bidded more than the previous winning bid " + fmt.Sprint(srv.winingBid.Amount))
			srv.winingBid = *bid
			return &proto.Acknowledgement{
				Status: proto.Status_SUCCESS,
			}, nil
		}
		return &proto.Acknowledgement{
			Status: proto.Status_FAIL,
		}, nil
	} else {
		ackno, err := srv.currentLeader.SendBid(ctx, bid)
		if err != nil {
			srv.selectleader()
			ackno, _ = srv.currentLeader.SendBid(ctx, bid)
		}
		return ackno, nil
	}
}

func (srv *AuctionServiceServer) Result(ctx context.Context, _ *proto.Empty) (*proto.Outcome, error) {
	if srv.isLeader {
		return &proto.Outcome{
			Winingbid:       &srv.winingBid,
			AuctionFinished: srv.auctionFinished,
		}, nil
	} else {
		outcome, err := srv.currentLeader.Result(ctx, &proto.Empty{})
		if err != nil {
			srv.selectleader()
			outcome, _ = srv.currentLeader.Result(ctx, &proto.Empty{})
		}
		srv.winingBid = *outcome.GetWiningbid()
		return outcome, nil
	}

}

func (srv *AuctionServiceServer) selectleader() {
	for i, elem := range srv.servernodes {
		if i == srv.serverIndex {
			srv.isLeader = true
			log.Println("I AM THE LEADER NOW")
			return
		}
		_, err := elem.HealthCheck(context.Background(), &proto.Empty{})
		if err == nil {
			srv.currentLeader = elem
			return
		}
	}
}

// GRPC for either getting the current winning bid from Leader or to register a crash
func (srv *AuctionServiceServer) HealthCheck(ctx context.Context, _ *proto.Empty) (*proto.Outcome, error) {
	//log.Println("HealthCheck called")
	return &proto.Outcome{
		Winingbid:       &srv.winingBid,
		AuctionFinished: srv.auctionFinished,
	}, nil
}
