package main

import (
	"chan_net/tcpc"
	"fmt"
	"log"
	//"time"
)

func main() {
	channel, err := tcpc.New[string](":3000", ":4000")

	if err != nil {
		log.Fatal(err)
	}
	
	channel.Sendch <- "Hello"

	msg := <- channel.Recvch

	fmt.Println("recived msg from channel over 4000:", msg)

}