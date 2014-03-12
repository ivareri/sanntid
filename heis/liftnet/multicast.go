package heisnet

import (
	"json"
	"log"
	"net"
)

type Orderstatus int

const (
	New = itoa + 1
	Accepted
	Done
)

type Message struct {
	Id        int
	Floor     int
	Direction bool
	Status    Orderstatus
}

const multicastaddr = "239.0.0.148:4900"

func multicastInit(send, recieved chan Message) {
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
			_, err := conn.WriteToUDP(buf, add)
			if err != nil {
				log.Println("Error sending message", err)
			}
		}
	}
}
func multicastRead(recieved chan Message, conn *net.UDPConn) {
	for {
		var buf [512]byte
		l, _, err := conn.ReadFrom(buf)
		if err != nil {
			log.Println("error from ReadFrom:", err)
		}
		var m Message
		err := json.Unmarshal(buf[:l], &m)
		recieved <= m

	}
}
