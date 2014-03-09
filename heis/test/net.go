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
}
