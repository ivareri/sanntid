package liftio

import (
	"log"
	"time"
)

const MAXFLOOR = 4
const DEFAULTSPEED = 200

type ButtonType int

const (
	Up ButtonType = iota
	Down
	Command
	Stop
	Obstruction
	Door
)

// Floor should be ignored when Button is  Stop or Obstruction
// Used for passing around keypresses
type Button struct {
	Floor  uint
	Button ButtonType
}

// Used for setting command and order lights
type Light struct {
	Floor  uint
	Button ButtonType
	On     bool
}

// Running is false when lift is stationary
type FloorStatus struct {
	Running   bool
	Floor     uint
	Direction bool
}

type motorType struct {
	speed     int
	direction bool
}

var (
	floorSeen = make(chan uint, 5)
	motor = make(chan motorType, 5)
)

// Initilazes elevator, starts runElevator routine, wich in turn starts readFloorSensor routine.
// Do not write or read from floorOrder and floor untill this function returns true
func Init(floorOrder *chan uint, light *chan Light, floor *chan FloorStatus, button *chan Button) bool {
	// Init hardware
	if !io_init() {
		log.Fatal("Error during HW init")
	}
	// turn off all lights
	io_clear_bit(LIGHT_STOP)
	io_clear_bit(DOOR_OPEN)
	io_clear_bit(LIGHT_COMMAND1)
	io_clear_bit(LIGHT_COMMAND2)
	io_clear_bit(LIGHT_COMMAND3)
	io_clear_bit(LIGHT_COMMAND4)
	io_clear_bit(LIGHT_UP1)
	io_clear_bit(LIGHT_UP2)
	io_clear_bit(LIGHT_UP3)
	io_clear_bit(LIGHT_DOWN2)
	io_clear_bit(LIGHT_DOWN3)
	io_clear_bit(LIGHT_DOWN4)
	go runIO(button, light)
	go runElevator(floorOrder, floor)
	return true
}

func runIO(button *chan Button, light *chan Light) {
	log.Println("entering for loop")
	for {
		readFloorSensor()
		runMotor()
		readButtons(*button)
		setLight(*light)
		time.Sleep(5*time.Millisecond)
	}
	log.Println("Exited for loop")
}

// Listens on floorOrder, and runs lift to given floor
// Returns status as it arrives at any floor
// Called from Init
func runElevator(floorOrder *chan uint, floor *chan FloorStatus) {
	var currentFloor, floorStop uint
	var lastFloor FloorStatus
	lastFloor.Direction = false

	// Go to closest floor downwards.
	// Do this to get a known state
		motor<-motorType{DEFAULTSPEED, lastFloor.Direction}
	  for {
			currentFloor = <-floorSeen
			if currentFloor != 0 {
				log.Println("Found floor")
				break
			}
		}
		motor<-motorType{0, lastFloor.Direction}
		lastFloor.Floor = currentFloor
		lastFloor.Running = false
	*floor <- lastFloor

	// Elevator should be in known state. Starting loop
	for {
		select {
		case newFloorStop := <-*floorOrder:
			if newFloorStop < 1 || newFloorStop > MAXFLOOR {
				log.Println("FloorOrder out of range:", newFloorStop)
			} else {
				floorStop = newFloorStop
			}
		case currentFloor = <-floorSeen:
			if currentFloor == 0 {
				break
			}
			if currentFloor == 1 || currentFloor == MAXFLOOR {
				motor <- motorType{DEFAULTSPEED, lastFloor.Direction}
				lastFloor.Floor = currentFloor
				lastFloor.Running = false
				*floor <- lastFloor
			}
			if currentFloor == floorStop {
				motor <- motorType{0, lastFloor.Direction}
				lastFloor.Floor = currentFloor
				floorStop = 0
				lastFloor.Floor = currentFloor
				lastFloor.Running = false
				*floor <- lastFloor
			}
			if currentFloor > 1 && currentFloor <= MAXFLOOR {
				lastFloor.Floor = currentFloor
				*floor <- lastFloor
			}
		default:
			if floorStop == 0 {
				break
			}
			if floorStop < lastFloor.Floor {
				lastFloor.Direction = true
			} else {
				lastFloor.Direction = false
			}
			if doorOpen && !lastFloor.Running {
				motor <- motorType{DEFAULTSPEED, lastFloor.Direction}
				lastFloor.Floor = currentFloor
				lastFloor.Running = true
				*floor <- lastFloor
			}
		}
	}
}
