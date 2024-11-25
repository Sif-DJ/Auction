package main

import (
	proto "Auction/grpc"
	"context"
	"log"
	"strconv"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	address := "5050"
	client := StartNode(address)
	_, err := client.SendBid(context.Background(), &proto.Bid{Node: address, Amount: 500})
	if err != nil {
		log.Println("SendBid error: ", err)
	}
	outcome, _ := client.Result(context.Background(), &proto.Empty{})
	log.Println("is auction done? [" + strconv.FormatBool(outcome.AuctionFinished) + "]")

}

func StartNode(address string) proto.AuctionServiceClient {
	conn, err := grpc.NewClient("localhost:"+address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Exception Error")
	}
	return proto.NewAuctionServiceClient(conn)
}
