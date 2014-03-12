package main

import (
	"../liftio"
	"log"
)
func main() {
//	order := make (chan uint)
//	status := make(chan liftio.FloorStatus)
	keypress := make(chan liftio.Button)
//	if !liftio.Init(&order, &status) {
//		log.Fatal("Error starting lift")
//	}
//	log.Println("Lift started")
go	liftio.ReadButtons(keypress)
	for {
		select {
			case bla:=<-keypress:
			log.Println("Keypress")
			log.Println(bla)
//			case bla:=<-status:
//			log.Println(bla)
		}
	}
}
