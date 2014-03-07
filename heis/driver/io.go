package io

import (
	"log"
)

type buttonType int
const (
	up Button = itoa+1
	down
	command
	stop
	obstruction
)

type floorButton struct {
	floor int
	button buttonType
}

func init(void) {
	// Init hardware
	if (!io_init()) {
		log.Fatal("Error during HW init")
	}

	// Zero all floor button lamps
  // Clear stop lamp, door open lamp, and set floor indicator to ground floor.
  // Return success.
    return true;
}

func runMotor(speed uint, direction bool) {

	// Invert direction in order to break elevator before stopping
	if speed == 0 {
		direction = !direction
	}
	if direction {
		io_set_bit(MOTORDIR)
	} else {
		io_clear_bit(MOTORDIR)
	}
	io_write_analog(MOTOR, 2048 + 4 * speed)
}

func readFloorsensor(floor chan uint) {
	currenFloor := -1;
	if io_read_bit(SENSOR1) && ( currentFloor != 1) {
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
		io_set_bit(DOOR_OPEN);
	} else {
		io_clear_bit(DOOR_OPEN);
	}
}
