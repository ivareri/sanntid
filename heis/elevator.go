package elevatorControl

import (
	"io"
	"localQueue"
)

func elevatorControl(){
	
	floorOrder := make(chan uint)
	buttonPress := make(chan Button)
	floor := make(chan FloorStatus)
	
	init(floorOrder, floor)
	go readButtons(buttonPress) // puts pressed buttons in the buttonPress chanel
	
	
}