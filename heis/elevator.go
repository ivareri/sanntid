package elevatorControl

import (
	"./liftio"
	"./localQueue"
	"log"
	"os"
	"time"
)

// Initialize lift
// Send orders to liftio
// Asign lifts to requests
// Add/delete orders/requests to/from locallocalQueue.Queue
func elevatorControl() {
	locallocalQueue.Queue := localQueue.Queue{}
	ReadlocalQueue.QueueFromFile(locallocalQueue.Queue) // If no previous queue:
	// logs error: "queue.txt doesn't exitst"

	floorOrder := make(chan uint)    	// channeling floor orders to io
	buttonPress := make(chan Button) 	// channeling button presses from io
	status := make(chan FloorStatus) 	// channeling the lifts status
	// Rename to LiftStatus?
	toNetwork := make(chan Message)   	// channeling messages to the network
	fromNetwork := make(chan Message) 	// channeling messages to the network
	light := make(chan Light)			

	init(&floorOrder, &status)
	MultiCastInit(&toNetwork, &fromNetwork)

	go ReadButtons(&buttonPress)
	//go GetOrder(floorOrder)?

	for {
		select {
		case buttonPressed := <-buttonPress:
			if buttonPressed.button == Up || buttonPressed.button == Down {
				log.Println("Request button %v pressed.", buttonPressed.button)
				// to network or take self
				// tell net anyhow
				// if taken self: call GetOrder
			} else if buttonPressed.button == Command {
				log.Println("Command button %v pressed.", buttonPressed.Floor)
				addLockalCommand(buttonPressed, locallocalQueue.Queue)
				// call GetOrder
			} else if buttonPressed.button == Stop {
				log.Println("Stop button pressed")
				// optional
			} else if buttonPressed.button == Obstuction {
				log.Println("Obstruction")
				// optional
			}
		case message := <-fromNet: 
			if message.Status == liftnet.Done {
				if message.Direction{
					light <- Light{message.Floor, Up , false}
				}else{
					light <- Light{message.Floor, Down , false}
				}
				
			}
			
			
		default:
			log.Println("No buttons pressed.")
			// teit å skrive  ut hele tiden?
		}
	}

}

func assignLift(toNetwork) {
	toNetwork <- jhk

}

// Nearest Car algorithm, returns Figure of Suitability
// Lift with largest FS should accept the request
func figureOfSuitability(request Message, status FloorStatus) int {
	reqDir := request.Direction
	reqFlr := request.Floor
	statDir := status.Direction
	statFlr := status.Floor
	if reqDir == statDir {
		// lift moving towards the requested floor and the request is in the same direction
		if (statDir && reqFlr > statFlr) || (!statDir && reqFlr < statFlr) {
			fs := MAXFLOOR + 1 - diff(reqFlr, statFlr)
		}
	} else {
		// lift moving towards the requested floor, but the request is in oposite direction
		if (statDir && reqFlr > statFlr) || (!statDir && reqFlr < statFlr) {
			fs := MAXFLOOR - diff(reqFlr, statFlr)
		} else {
			fs := 1
		}
	}
	return fs
}

// timer thingy
// now := time.Now()
// diff := now.sub(then)
// sum := then.add(diff)
// diff.Hours() osv -> Nanoseconds
