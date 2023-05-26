package main

import (
	"chan_net/tcpc"
	"fmt"
	"log"
)

func main() {
	channel, err := tcpc.New[string](":4000", ":3000")

	if err != nil {
		log.Fatal(err)
	}

	msg := <- channel.Recvch

	fmt.Println("received msg from channel over 3000 :", msg)

	channel.Sendch <- "hi"
}