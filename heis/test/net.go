package main

import (
	"../liftnet/"
	"fmt"
	"log"
	"strconv"
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
	go liftnet.MulticastInit(&send, &recieved, iface)
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
		fmt.Println("Id: " + strconv.Itoa(int(as.Id)) + " Floor: " + strconv.Itoa(int(as.Floor)) + "Status: " + strconv.Itoa(int(as.Status)))

	}
}
