package liftnet

import (
	"encoding/json"
	"log"
	"net"
)

type Orderstatus int

const (
	New Orderstatus = iota
	Accepted
	Done
)

// Used for communication between lifts
type Message struct {
	Id        int
	Floor     int
	Direction bool
	Status    Orderstatus
}

const multicastaddr = "239.0.0.148:49153"

// Sets up network sender and reciver
func MulticastInit(send chan Message, recieved chan Message) {
	group, err := net.ResolveUDPAddr("udp", multicastaddr)
	if err != nil {
		log.Println("error from ResolveUDPAddr:", err)
		return
	}

	conn, err := net.ListenMulticastUDP("udp", nil, group)
	if err != nil {
		log.Println("error from ListenMulticastUDP:", err)
		return
	}
	go multicastRead(recieved, conn)
	go multicastSend(send, conn, group)
}
func multicastSend(send chan Message, conn *net.UDPConn, addr *net.UDPAddr) {
	for {
		m := <-send
		buf, err := json.Marshal(m)
		if err != nil {
			log.Println("Error encoding message: ", err)
		} else {
			_, err := conn.WriteToUDP(buf, addr)
			if err != nil {
				log.Println("Error sending message", err)
			}
		}
	}
}
func multicastRead(recieved chan Message, conn *net.UDPConn) {
	for {
		buf :=make([]byte, 512)
		l, _, err := conn.ReadFrom(buf)
		if err != nil {
			log.Println("error from ReadFrom:", err)
		}
		var m Message
		er := json.Unmarshal(buf[:l], &m)
		if er != nil {
			log.Println("Error unpacking", er)
		} else {
		//	recieved <= m
		}
	}
}
