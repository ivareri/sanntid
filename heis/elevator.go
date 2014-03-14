package elevatorControl

import (
	"io"
	"localQueue"
	"time"
	"io/ioutil"
)

type Queue struct {
	upQueue [4]bool
	downQueue [4]bool
	commandQueue [4]bool
}

func elevatorControl(){
	localQueue := Queue{}

	//Check if there is an existing queue on file
	
	
	floorOrder := make(chan uint) 			// channeling floor orders to io 
	buttonPress := make(chan Button)		// channeling button presses from io
	status := make(chan FloorStatus)		// channeling the lifts status
	sendToNet := make(chan Message)			// channeling messages to the network
	recieveFromNet := make(chan Message)	// channeling mesages to the network

	init(floorOrder, status)			

	go ReadButtons(buttonPress) 		
	go GetOrder(floorOrder)
	
	for{
		select{
		case buttonPressed := <- buttonPress:
			if buttonPressed.button == Up || buttonPressed.button == Down {
				log.Println("Request button %v pressed.", buttonPressed.button)
				// to network or take self
			} else if buttonPressed.button == Command {
				log.Println("Command button %v pressed.", buttonPressed.Floor)
				addLockalCommand(buttonPressed, localQueue)
			} else if buttonPressed.button == Stop {
				log.Println("Stop button pressed")
				EmergencyStop(true)
				// do or check something
			} else if buttonPressed.button == Obstuction {
				//does this even belong here ?
			}
		default
			log.Println("No buttons pressed.")	
		}
	}
	
}

// timer thingy
// now := time.Now()
// diff := now.sub(then)
// sum := then.add(diff) 
// diff.Hours() osv -> Nanoseconds


