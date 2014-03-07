package io

import (
	"log"
)

type buttonType int

const (
	up Button = itoa + 1
	down
	command
	stop
	obstruction
)

type floorButton struct {
	floor  int
	button buttonType
}

func init(floor chan uint) {
	// Init hardware
	if !io_init() {
		log.Fatal("Error during HW init")
	}

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

	go readFloorSensor(floor)
	runMotor(200, false)
	<-floor
	runMotor(0, false)
	// Return success.
	return true
}
func runElevator(floorOrder chan uint, lastFloor chan uint) {
	floor := make(chan uint)
	var currentFloor uint = 0
	var floorStop [4]bool = false
	go readFloorSensor(floor)
	runMotor(200, false)
	for {
		currentFloor<-floor
		if currentFloor != 0 {
			break
		}
	}
	runMotor(0, false)
	lastFloor <- currentFloor
	for {
		select {
		case i := <-floorOrder:
			floorStop[i] = True
		case i := <-lastFloor:
			if floorStop[i] {
				runMotor(0, direction)
			}
		default:
			// noka lurt for å kjøre\stoppe heis
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

func setLight() {
	// write some fancy code
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
