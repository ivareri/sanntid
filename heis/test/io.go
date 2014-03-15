package main

import (
	"../liftio"
	"log"
)

var (
	order    = make(chan uint, 5)
	status   = make(chan liftio.FloorStatus, 10)
	keypress = make(chan liftio.Button, 10)
	light    = make(chan liftio.Light, 10)
)

func main() {

	if !liftio.Init(&order, &light, &status, &keypress) {
		log.Fatal("Error starting lift")
	}
	log.Println("Lift started")
	for {
		select {
		case bla := <-keypress:
			testkeypres(bla)
		case bla := <-status:
			log.Println("Status:", bla)
		}
	}
}

func testkeypres(bla liftio.Button) {
	log.Println("Keypress: ", bla)
	light <- liftio.Light{bla.Floor, bla.Button, true}
	if bla.Button == liftio.Command {
		order <- bla.Floor
	}

}
