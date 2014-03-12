package main

import (
	"fmt"
	"../liftnet/"
)
func main () {
	a, err := liftnet.FindIP()
	if err != nil {
		fmt.Println("Error", err)
	}
	fmt.Println(liftnet.FindID(a))
	send := make(chan liftnet.Message)
	recieved := make(chan liftnet.Message)
	liftnet.MulticastInit(send, recieved)
	var bla liftnet.Message
	bla.Id = 1
	bla.Floor = 2
	bla.Direction = false
	bla.Status = liftnet.New
	send <- bla
	<-recieved
}
