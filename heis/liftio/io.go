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

// Floor isLiftStatus ignored when Button is  Stop or Obstruction
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
)

// Initilazes hardware and starts IO routines
// Do not write or read to any channels untill this function returns true
func Init(floorOrder *chan uint, light *chan Light, floor *chan LiftStatus, button *chan Button) bool {
	// Init hardware
	if !io_init() {
		log.Fatal("Error during HW init")
	}
	// Stop motor and turn off all lights
	io_write_analog(MOTOR, 0)
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
	log.Println("Cleared lights. Starting go routines")
	lightch = light
	floorch = floor
	go runIO(button)
	go runElevator(floorOrder)
	return true
}

// Threads writing\reading from IO might cause bugs
func runIO(button *chan Button) {
	for {
		setLight(*lightch)
		readFloorSensor()
		runMotor()
		readButtons(*button)
		time.Sleep(5 * time.Millisecond)
	}
}

// Listens on floorOrder, and runs lift to given floor
// Returns status as it arrives at any floor
// Called from Init
func runElevator(floorOrder *chan uint) {
	var currentFloor, stopFloor uint
	var status LiftStatus
	status.Direction = false

	// Go to closest floor downwards.
	// Do this to get a known state
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

	// Elevator should be in known state. Starting loop
	for {
		select {
		case newStopFloor := <-*floorOrder:
			if checkFloorRange(newStopFloor) {
				stopFloor = newStopFloor
			}
		case currentFloor = <-floorSeen:
			newFloor(currentFloor, &status, &stopFloor)
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
}

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

func goToFloor(currentFloor uint, status *LiftStatus, stopFloor *uint) {
	if !status.Door && !status.Running {
		status.Direction = status.Floor < *stopFloor
		motor <- motorType{DEFAULTSPEED, status.Direction}
		status.Floor = currentFloor
		status.Running = true
		*floorch <- *status
	}
}

func checkFloorRange(floor uint) bool {
	if floor < 1 || floor > MAXFLOOR {
		log.Println("FloorOrder out of range:", floor)
		return false
	} else {
		return true
	}
}

func newFloor(currentFloor uint, status *LiftStatus, stopFloor *uint) {
	switch currentFloor {
	case 0:
		if status.Door {
			log.Fatal("FATAL ERROR: Elevator should not be moving. Door is open")
		}
		if !status.Running {
			log.Fatal("FATAL ERROR: Elevator should not be moving. Motor is off")
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
