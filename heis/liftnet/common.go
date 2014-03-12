package liftnet

import (
	"net"
	"os"
	"log"
	"strconv"
	"strings"
	"errors"
)

//Returns IPv4 address for lift
func FindIP() (string, error) {
	name, err := os.Hostname()
	if err != nil {
		return "", err
	}

	addrs, err := net.LookupHost(name)
	if err != nil {
		return "", err
	}
	for _, a := range addrs {
		if strings.Contains(a, ".") {
			return a, nil
		}
	}
	return "", errors.New("Unable to find IPv4 address")
}

//converts IPv4 address to ID string.
// 3 last digits from IPv4 address
func FindID(a string) int {
	id, err := strconv.Atoi(strings.Split(a, ".")[3])
	if err != nil {
		log.Fatal("Error converting IP to ID", err)
	}
	return id
}
