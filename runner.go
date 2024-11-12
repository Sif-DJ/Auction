package main

import (
	client "Auction/Client"
	"log"
)

func main() {

	var clients []string
	clients = append(clients, "5050")
	clients = append(clients, "5051")
	clients = append(clients, "5052")
	clients = append(clients, "5053")
	clients = append(clients, "5054")
	for i := 0; i < len(clients); i++ {
		go client.StartNode(clients[i], clients)
	}

	log.Println("Now in forever loop")
	for {

	}

}
