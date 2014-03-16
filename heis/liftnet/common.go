package liftnet

import (
	"errors"
	"log"
	"net"
	"strconv"
	"strings"
)

var quit = make(chan bool)
//Returns IPv4 address for lift
func FindIP() (string, *net.Interface, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return "", nil, err
	}
	for _, iface := range ifaces {
		addrs, _ := iface.Addrs()
		for _, a := range addrs {
			if strings.Contains(a.String(), "129.") {
				return a.String(), &iface, nil
			}
		}
	}
	return "", nil, errors.New("Unable to find IPv4 address")
}

//converts IPv4 address to ID string.
// 3 last digits from IPv4 address
func FindID(a string) int {
	log.Println(a)
	id, err := strconv.Atoi(strings.Split(a, ".")[3][:3])
	if err != nil {
		log.Fatal("Error converting IP to ID", err)
	}
	return id
}
