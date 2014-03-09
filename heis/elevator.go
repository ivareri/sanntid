package elevatorControl

import (
	"io"
	"localQueue"
)

func elevatorControl(){
	floorOrder := make(chan uint)
	buttonPress := make(chan Button)
	floor := make(chan FloorStatus)
	
	// [floor][0]=up, [floor][1]=down, [floor][2]= command 
	// localqueue[4][3] bool
	
	init(floorOrder, floor)
	go ReadButtons(buttonPress) // puts pressed buttons in the buttonPress chanel
	go GetOrder(floorOrder)
	
	// func that gets something from manager and runs the different localQueue funcs
	
}