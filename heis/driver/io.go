package io

import (
	"log"
)

const MAXFLOOR = 4
const DEFAULTSPEED = 200

type buttonType int

const (
	up buttonType = itoa + 1
	down
	command
	stop
	obstruction
)

type Button struct {
	floor  int
	button buttonType
}

type lightType int
const (
	up lightTYpe = itoa +1
	down
	command
	stop
	door
)
type Light struct {
	floor int
	light lightType
	on bool
}

func init(floorOrder chan uint, floor chan uint) {
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

	go runElevator(floorOrder, floor)
	// wait for eleavtor to arrive at a floor
	<-floor
	return true
}

func runElevator(floorOrder chan uint, floor chan uint) {
	floorSeen := make(chan uint)
	var currentFloor, lastFloor, floorStop uint = 0
	direction := false
	go readFloorSensor(floorSeen)

	// Go to closest floor downwards.
	// Do this to get a known state
	runMotor(DEFAULTSPEED, direction)
	for {
		currentFloor <- floorSeen
		if currentFloor != 0 {
			break
		}
	}
	floor <- currentFloor
	runMotor(0, direction)

	for {
		select {
		case newFloorStop := <-floorOrder:
			if newFloorStop < 1 || newFloorStop > 4 {
				log.Println("FloorOrder out of range:", newFloorStop)
			} else {
				floorStop = newFloorStop
			}
		case currentFloor <- floorSeen:
			if currentFloor == floorStop {
				runMotor(0, direction)
				floorStop = 0
			}
			if currentFloor > 1 && currentFloor < MAXFLOOR {
				lastFloor = currentFloor
				floor <- lastFloor // TODO: finn på noka lurt her, som å sende en struct med heisstatus
			}
		default:
			if floorStop == 0 {
				break
			}
			if floorStop < lastFloor {
				direction = false
			} else {
				direction = true
			}
			runMotor(DEFAULTSPEED, direction)
		}
	}
}

func runMotor(speed uint, direction bool) {
	//TODO: Should save last direction for breaking)

	// Invert direction in order to break elevator before stopping
	if speed == 0 {
		direction = !direction
	}
	if direction {
		io_set_bit(MOTORDIR)
	} else {
		io_clear_bit(MOTORDIR)
	}
	io_write_analog(MOTOR, 2048+4*speed)
}

func readFloorsensor(floor chan uint) {
	currenFloor := -1
	if io_read_bit(SENSOR1) && (currentFloor != 1) {
		setFloorLight(1)
		floor <- 1
	} else if io_read_bit(SENSOR2) && (currentFloor != 2) {
		setFloorLight(2)
		floor <- 2
	} else if io_read_bit(SENSOR3) && (currentFloor != 3) {
		setFloorLight(3)
		floor <- 3
	} else if io_read_bit(SENSOR4) && (currentFloor != 4) {
		setFloorLight(4)
		floor <- 4
	} else if currentFloor != 0 {
		floor <- 0
	}
}

func setFloorLight(floor int) {
	if (floor < 1) || (floor > 4) {
		log.Fatal("Floor out of range: ", floor)
	}
	switch floor {
	case 1:
		io_clear_bit(FLOOR_IND1)
		io_clear_bit(FLOOR_IND2)
	case 2:
		io_set_bit(FLOOR_IND1)
		io_clear_bit(FLOOR_IND2)
	case 3:
		io_clear_bit(FLOOR_IND1)
		io_set_bit(FLOOR_IND2)
	case 4:
		io_set_bit(FLOOR_IND1)
		io_set_bit(FLOOR_IND2)
	}
}

func setLight(light Light) {
// not done yet
		io_set_bit(LIGHT_COMMAND1)
		io_set_bit(LIGHT_COMMAND2)
		io_set_bit(LIGHT_COMMAND3)
		io_set_bit(LIGHT_COMMAND4)
		io_set_bit(LIGHT_UP1)
		io_set_bit(LIGHT_UP2)
		io_set_bit(LIGHT_DOWN2)
		io_set_bit(LIGHT_UP3)
		io_set_bit(LIGHT_DOWN3)
		io_set_bit(LIGHT_DOWN4)
		io_set_bit(DOOR_OPEN)

		io_clear_bit(LIGHT_COMMAND1)
		io_clear_bit(LIGHT_COMMAND2)
		io_clear_bit(LIGHT_COMMAND3)
		io_clear_bit(LIGHT_COMMAND4)
		io_clear_bit(LIGHT_UP1)
		io_clear_bit(LIGHT_UP2)
		io_clear_bit(LIGHT_DOWN2)
		io_clear_bit(LIGHT_UP3)
		io_clear_bit(LIGHT_DOWN3)
		io_clear_bit(LIGHT_DOWN4)
		io_clear_bit(DOOR_OPEN)

}

func emergencyStop(bool stop) {
	if stop {
		io_set_bit(LIGHT_STOP)
	} else {
		io_clear_bit(LIGHT_STOP)
	}
}

func doorOpen(bool open) {
	if open {
		io_set_bit(DOOR_OPEN)
	} else {
		io_clear_bit(DOOR_OPEN)
	}
}
