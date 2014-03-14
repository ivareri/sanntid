package main

import (
	"../liftio"
	"log"
)
func main() {
	order := make (chan uint)
	status := make(chan liftio.FloorStatus)
	keypress := make(chan liftio.Button)
	if !liftio.Init(&order, &status, &keypress) {
		log.Fatal("Error starting lift")
	}
	log.Println("Lift started")
	for {
		select {
			case bla:=<-keypress:
			log.Println("Keypress: ", bla)
			if bla.Button == liftio.Command {
				order <-bla.Floor
			}
			case bla:=<-status:
			log.Println("Status:", bla)
		}
	}
}
