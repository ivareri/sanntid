package liftnet

import (
	"errors"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

const bcastPort = 63000

func sendStatus(statusch chan string, stopUDPBroadcast chan bool) {
	broadcast, err := net.ResolceUDPAddr("udp", IPv4bcast+":"+bcastPort)
	if err != nil {
		log.Fatal("Fatal error:", err)
	}
	con, err := net.DialUDP("udp", nil, broadcast)
	if err != nil {
		log.Fatal("Fatal error:", err)
	}
	defer con.Close()
	for {
		select {
		case <-stopUDPBroadcast:
			return
		case status <- statusch:
			enc, err := json.Marshal(status)
			if err != nil {
				log.Println("Error: ", err)
			} else {
				encoded = en
			}
		default:
			_, err := con.Write(encoded)
			if err != nil {
				log.Println("Error sending status:", err)
			}
			time.Sleep(500 * time.Millisecond)
		}
	}
}

func listenStatus(closeUDPListen chan bool, statuschan chan string) {
	udpAddr, err := net.ResolveUDPAddr("udp", IPv4addr+":"+bcastPort)
	if err != nil {
		log.Fatal("Error: ", err)
	}
	netlisten, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		log.Fatal("Error: ", err)
	}
	defer netlisten.Close()
	buffer, oob := make([]byte, 1024)
	for {
		netlisten.SetDeadline(time.Now().Add(1e9))
		n, _, addr, _, err := netlisten.ReadMsgUDP(buffer, oob)
		if err != nil {
			if opErr, ok := err.(*net.OpError); ok && opErr.Timeout() {
				continue
			} else {
				log.Fatal("Error:", err)
			}
		} else {
			status, _ := unpackUDP(addr, buffer)
			statuschan <- status
		}
		switch {
		case <-closeUDPListen:
			return
		}
	}
}

func unpackUDP(addr net.IPAddr, buffer []byte) (string, error) {

	status, err := json.Unmarshal(buffer[:n])
	if err != nil {
		log.Println("Error unmarshaling: ", err)
	}

}
