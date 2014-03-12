package elevatorControl

import (
	"io"
	"localQueue"
)


func elevatorControl(){
	// localQueue indexing: [floor][0]=Up, [floor][1]=Down, [floor][2]=Command
	// stupid way to init
	// should it be a type instead?
	var localQueue := [4][3]bool {{false, false, false},{false, false, false},{false, false, false},{false, false, false}}
	
	//Check file if there is an existing queue
	
	
	floorOrder := make(chan uint)
	buttonPress := make(chan Button)
	floor := make(chan FloorStatus)

	init(floorOrder, floor)

	go ReadButtons(buttonPress) // puts pressed buttons in the buttonPress chanel
	go GetOrder(floorOrder)

	// func that gets something from manager and runs the different localQueue funcs

	
	
	for{
		select{
		case buttonPressed := <- buttonPress:
			if buttonPressed.button == Up || buttonPressed.button == Down {
				// to network or take self
			} else if buttonPressed.button == Command {
				addLockalCommand(buttonPressed, localQueue)
			}
			
		}
	}
	
}

