package elevatorControl

import (
	"liftio"
	"localQueue"
	"time"
	"log"
	"os"
)

// change names to avoid confusion with buttons?
type Queue struct {
	Up [4]bool
	Down [4]bool
	Command [4]bool
}

// Initialize lift
// Send orders to liftio
// Asign lifts to requests
// Add/delete orders/requests to/from localQueue
func elevatorControl(){
	localQueue := Queue{}
	ReadQueueFromFile(localQueue) 			// If no previous queue:
											// logs error: "queue.txt doesn't exitst" 
		
	floorOrder := make(chan uint) 			// channeling floor orders to io 
	buttonPress := make(chan Button)		// channeling button presses from io
	status := make(chan FloorStatus)		// channeling the lifts status
	// Rename to LiftStatus?
	toNetwork := make(chan Message)				// channeling messages to the network
	fromNetwork := make(chan Message)	  		// channeling messages to the network

	init(&floorOrder, &status)
	MultiCastInit(&toNetwork, &fromNetwork)			

	go ReadButtons(&buttonPress) 		
	go GetOrder(&floorOrder)
	
	for{
		select{
		case buttonPressed := <-buttonPress:
			if buttonPressed.button == Up || buttonPressed.button == Down {
				log.Println("Request button %v pressed.", buttonPressed.button)
				// to network or take self
				// tell net anyhow
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
		case request := <-fromNet:
			FS := figureOfSuitability(request,)
			
		default:
			log.Println("No buttons pressed.")
			// teit å skrive  ut hele tiden?	
		}
	}
	
}

func assignLift(toNetwork ){
	toNetwork <- jhk
	
}

// Returns int fs
// Lift with largest fs should accept the request 
func figureOfSuitability(request Message, status FloorStatus) int {
	reqDir := request.Direction
	reqFlr := request.Floor
	statDir := status.Direction
	statFlr := status.Floor
	if reqDir == statDir && reqFlr > statFlr { // if lift moving towards req flr and req in same dir: N+1-d
		fs := MAXFLOOR + 1 - diff(reqFlr,statFlr)
	} else if {
		
	} 
	
	
	else if  !reqDir && statDir && requFlr < statFlr {
		fs := MAXFLOOR + 1 - diff(reqFlr,statFlr)
	} else {
		fs := 1
	}
	return fs
}

// timer thingy
// now := time.Now()
// diff := now.sub(then)
// sum := then.add(diff) 
// diff.Hours() osv -> Nanoseconds


