package liftnet

import (
	"encoding/json"
	"log"
	"net"
	"time"
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
	Floor     uint
	Direction bool
	Status    Orderstatus
	TimeSent  time.Time
	TimeRecv  time.Time
}

const multicastaddr = "239.0.0.148:49153"

// Sets up network sender and reciver
func MulticastInit(send chan Message, recieved chan Message, iface *net.Interface) {
	group, err := net.ResolveUDPAddr("udp", multicastaddr)
	if err != nil {
		log.Println("error from ResolveUDPAddr:", err)
		return
	}

	conn, err := net.ListenMulticastUDP("udp", iface, group)
	if err != nil {
		log.Println("error from ListenMulticastUDP:", err)
		return
	}
	defer conn.Close()
	log.Println("STarting reader")
	go multicastRead(recieved, conn)
	log.Println("Starting sender")
	go multicastSend(send, conn, group)
	<-quit
}
func multicastSend(send chan Message, conn *net.UDPConn, addr *net.UDPAddr) {
	for {
		m := <-send
		log.Println("Sending msg")
		buf, err := json.Marshal(m)
		if err != nil {
			log.Println("Error encoding message: ", err)
		} else {
			n, err := conn.WriteToUDP(buf, addr)
			if err != nil {
				log.Println("Error sending message", err)
			} else {
				log.Println("Sendt bytes", n)
			}
		}
	}
}
func multicastRead(recieved chan Message, conn *net.UDPConn) {
	for {
		buf := make([]byte, 512)
		l, _, err := conn.ReadFrom(buf)
		if err != nil {
			log.Println("error from ReadFrom:", err)
		}
		log.Println("Found message, lenght: ", l)
		var m Message
		er := json.Unmarshal(buf[:l], &m)
		if er != nil {
			log.Println("Error unpacking", er)
		} else {
			m.TimeRecv = time.Now()
			recieved <- m
		}
	}
}
