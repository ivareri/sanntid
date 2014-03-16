package main

import (
	"../liftnet/"
	"fmt"
	"log"
	"time"
)

func main() {
	a, iface, err := liftnet.FindIP()
	if err != nil {
		fmt.Println("Error", err)
	}
	fmt.Println(liftnet.FindID(a))
	send := make(chan liftnet.Message)
	recieved := make(chan liftnet.Message)
	go liftnet.MulticastInit(send, recieved, iface)
	time.Sleep(20 * time.Millisecond)
	log.Println(iface.MulticastAddrs())

	var bla liftnet.Message
	bla.Id = 1
	bla.Floor = 2
	bla.Direction = false
	bla.Status = liftnet.New
	bla.TimeSent = time.Now()
	bla.TimeRecv = time.Now()
	send <- bla
	for {
		as := <-recieved
		fmt.Println("Id: " + string(as.Id) + " Floor: " + string(as.Floor) + "Status: " + string(as.Status))

	}
}
