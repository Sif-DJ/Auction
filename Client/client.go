package main

import (
	proto "Auction/grpc"
	"bufio"
	"context"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type clientstr struct {
	servernodesString []string
	connection        proto.AuctionServiceClient
}

func main() {
	clnt := &clientstr{
		servernodesString: []string{"5053", "5052", "5051", "5050"},
	}
	clientName := "Client" + fmt.Sprint(rand.Intn(256))
	log.Println("You are", clientName)
	clnt.StartNode()
	outcome := &proto.Outcome{
		AuctionFinished: false,
	}
	reader := bufio.NewReader(os.Stdin)
	for !outcome.AuctionFinished {
		read, _ := reader.ReadString('\n')       // Get Read
		input, _, _ := strings.Cut(read, "\r\n") // Cut off weird extras
		if input == "outcome" {
			outcome, err := clnt.connection.Result(context.Background(), &proto.Empty{}) //tries to get the current outcome
			if err != nil {                                                              //if it fails to get the outcome because server is down
				clnt.StartNode()
				outcome, _ = clnt.connection.Result(context.Background(), &proto.Empty{}) //tries to get the current outcome after finding a new server
			}
			if outcome.GetAuctionFinished() {
				log.Println("Auction finished and the winning bid is " + outcome.GetWiningbid().GetNode() + " with a bid of " + strconv.Itoa(int(outcome.GetWiningbid().GetAmount())))
			} else {
				log.Println("Highest bid is " + outcome.GetWiningbid().GetNode() + " with a bid of " + strconv.Itoa(int(outcome.GetWiningbid().GetAmount())))
			}
		} else {
			value, err := strconv.Atoi(input) // Convert to Int
			if err != nil {
				log.Println("Bad input, try again")
				continue
			}
			endValue := int32(value) // Convert to Int32
			bid := &proto.Bid{
				Node:   clientName,
				Amount: endValue,
			}
			ack, err := clnt.connection.SendBid( // Send current bid and receive an acknowledgement
				context.Background(),
				bid,
			)
			if err != nil { // Maybe there is a wrong client, so it is time to connect to another
				clnt.StartNode()
				ack, _ = clnt.connection.SendBid( // Send current bid to new client and receive an acknowledgement
					context.Background(),
					bid,
				)
			}
			if ack.Status == proto.Status_SUCCESS {
				log.Println("Succesfully bidded", endValue)
			}
			if ack.Status == proto.Status_FAIL {
				log.Print("Bidded lower than or equal to the current highest or the auction is closed")
			}
		}
	}

}

func (clnt *clientstr) StartNode() {
	for _, elem := range clnt.servernodesString {
		conn, err := grpc.NewClient("localhost:"+elem, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			log.Fatalf("Exception Error")
		}
		_connection := proto.NewAuctionServiceClient(conn)
		_, err = _connection.HealthCheck(context.Background(), &proto.Empty{})
		if err == nil {
			clnt.connection = _connection
			return
		}
	}
	log.Fatalln("Could not connect to any server")
}
