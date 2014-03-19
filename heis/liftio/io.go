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
	door // Not an actual button. Used for door light only, hence not exported
)

// Used for passing around keypresses
// Floor is ignored when Button is Stop or Obstruction
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

type LiftStatus struct {
	Running   bool
	Floor     uint
	Direction bool
	Door      bool
}

type motorType struct {
	speed     int
	direction bool
}

var (
	floorSeen = make(chan uint, 5)
	motor     = make(chan motorType, 5)
	lightch   *chan Light
	floorch   *chan LiftStatus
	doorTimer = make(chan bool, 2)
	quit *chan bool
)

// Called by Runlift
// Initilazes hardware and starts IO routines
// Do not write or read to any channels untill this function returns true
func IOInit(floorOrder *chan uint, light *chan Light, floor *chan LiftStatus, button *chan Button, quit *chan bool) bool {
	// Init hardware
	if !io_init() {
		log.Fatal("Error during HW init")
	}
	// Stops motor and turns off all lights
	ioShutDown()
	log.Println("Cleared lights. Starting go routines")
	lightch = light
	floorch = floor
	go runIO(button)
	go executeOrder(floorOrder)
	return true
}

// Called by IOInit as go routine
// Threads writing\reading from IO might cause bugs
func runIO(button *chan Button) {
	for {
		select {
		case <-*quit:
			ioShutDown()
			return
		default:
			setLight(*lightch)
			readFloorSensor()
			runMotor()
			readButtons(*button)
			time.Sleep(5 * time.Millisecond)
		}
	}
}

// Called by IOInit as go routine
// Listens on floorOrder, and runs lift to given floor
// Returns status as it arrives at any floor
func executeOrder(floorOrder *chan uint) {
	var currentFloor, stopFloor uint
	var status LiftStatus
	status.Direction = false

	// Go to closest floor downwards to get a known state.
	motor <- motorType{DEFAULTSPEED, status.Direction}
	for {
		currentFloor = <-floorSeen
		if currentFloor != 0 {
			log.Println("Found floor")
			break
		}
	}
	motor <- motorType{0, status.Direction}
	status.Floor = currentFloor
	status.Running = false
	*floorch <- status

	// Lift in known state. Starting loop
	for {
		select {
		case <-*quit:
			return
		case newStopFloor := <-*floorOrder:
			if checkFloorRange(newStopFloor) {
				stopFloor = newStopFloor
			}
		case currentFloor = <-floorSeen:
			updateStatus(currentFloor, &status)
		case <-doorTimer:
			*lightch <- Light{0, door, false}
			status.Door = false
			*floorch <- status
		default:
			time.Sleep(5 * time.Millisecond)
			if stopFloor != 0 {
				stopAtFloor(currentFloor, &status, &stopFloor)
				goToFloor(currentFloor, &status, &stopFloor)
			}
		}
	}
}

// Called by executeOrder
func stopAtFloor(currentFloor uint, status *LiftStatus, stopFloor *uint) {
	if status.Floor == *stopFloor {
		motor <- motorType{0, status.Direction}
		status.Floor = currentFloor
		status.Running = false
		status.Door = true
		*lightch <- Light{0, door, true}
		go func() {
			time.Sleep(3 * time.Second)
			doorTimer <- true
		}()
		*stopFloor = 0
		*floorch <- *status
	}
}

// Called by executeOrder
func goToFloor(currentFloor uint, status *LiftStatus, stopFloor *uint) {
	if !status.Door && !status.Running {
		status.Direction = status.Floor < *stopFloor
		motor <- motorType{DEFAULTSPEED, status.Direction}
		status.Floor = currentFloor
		status.Running = true
		*floorch <- *status
	}
}

// Called by executeOrder
func checkFloorRange(floor uint) bool {
	if floor < 1 || floor > MAXFLOOR {
		log.Println("FloorOrder out of range:", floor)
		return false
	} else {
		return true
	}
}

// Called by executeOrder
func updateStatus(currentFloor uint, status *LiftStatus) {
	switch currentFloor {
	case 0:
		if status.Door {
			log.Fatal("FATAL ERROR: Lift should not be moving. Door is open")
		}
		if !status.Running {
			log.Fatal("FATAL ERROR: Lift should not be moving. Motor is off")
		}
		return
	case 1, MAXFLOOR:
		log.Println("Maxfloor")
		motor <- motorType{0, status.Direction}
		status.Floor = currentFloor
		status.Running = false
		*floorch <- *status
		log.Println("Floor range")
	case 2, 3:
		log.Println("Floor 2,3")
		if currentFloor != status.Floor {
			status.Floor = currentFloor
			*floorch <- *status
		}
	default:
		log.Println("Detected floor out of range, ignoring :", currentFloor)
	}
}
