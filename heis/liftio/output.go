package liftio

import (
	"log"
	"time"
)

var doorOpen bool

// Called from runElevator
func runMotor() {
	// Invert direction in order to break elevator before stopping
	select {
	case bla := <-motor:
		if bla.speed == 0 {
			bla.direction = !bla.direction
		}
		if bla.direction {
			io_set_bit(MOTORDIR)
		} else {
			io_clear_bit(MOTORDIR)
		}
		time.Sleep(1 * time.Millisecond)
		io_write_analog(MOTOR, 2048+4*int(bla.speed))
	default:
		return
	}
}

// Sets order\command lights
// TODO: ugly beast. Should be a cleaner way of doing this
func setLight(lightch chan Light) {
	select {
	default:
		return
	case light := <-lightch:
		if light.On {
			switch light.Floor {
			case 1:
				switch light.Button {
				case Command:
					io_set_bit(LIGHT_COMMAND1)
				case Up:
					io_set_bit(LIGHT_UP1)
				}
			case 2:
				switch light.Button {
				case Command:
					io_set_bit(LIGHT_COMMAND2)
				case Up:
					io_set_bit(LIGHT_UP2)
				case Down:
					io_set_bit(LIGHT_DOWN2)
				}
			case 3:
				switch light.Button {
				case Command:
					io_set_bit(LIGHT_COMMAND3)
				case Up:
					io_set_bit(LIGHT_UP3)
				case Down:
					io_set_bit(LIGHT_DOWN3)
				}
			case 4:
				switch light.Button {
				case Command:
					io_set_bit(LIGHT_COMMAND4)
				case Down:
					io_set_bit(LIGHT_DOWN4)
				}
			}
		} else {
			switch light.Floor {
			case 1:
				switch light.Button {
				case Command:
					io_clear_bit(LIGHT_COMMAND1)
				case Up:
					io_clear_bit(LIGHT_UP1)
				}
			case 2:
				switch light.Button {
				case Command:
					io_clear_bit(LIGHT_COMMAND2)
				case Up:
					io_clear_bit(LIGHT_UP2)
				case Down:
					io_clear_bit(LIGHT_DOWN2)
				}
			case 3:
				switch light.Button {
				case Command:
					io_clear_bit(LIGHT_COMMAND3)
				case Up:
					io_clear_bit(LIGHT_UP3)
				case Down:
					io_clear_bit(LIGHT_DOWN3)
				}
			case 4:
				switch light.Button {
				case Command:
					io_clear_bit(LIGHT_COMMAND4)
				case Down:
					io_clear_bit(LIGHT_DOWN4)
				}
			}
		}
	}
}

// Called from readFloorSensor
func setFloorLight(floor int) {
	if (floor < 1) || (floor > 4) {
		log.Fatal("Floor out of range: ", floor)
	}
	switch floor {
	case 1:
		io_clear_bit(FLOOR_IND1)
		io_clear_bit(FLOOR_IND2)
	case 2:
		io_set_bit(FLOOR_IND2)
		io_clear_bit(FLOOR_IND1)
	case 3:
		io_clear_bit(FLOOR_IND2)
		io_set_bit(FLOOR_IND1)
	case 4:
		io_set_bit(FLOOR_IND1)
		io_set_bit(FLOOR_IND2)
	}
}

// Currently only sets the stop light.
// Future revisions should implement actual emergency stop procedures
func EmergencyStop(stop bool) {
	if stop {
		io_set_bit(LIGHT_STOP)
	} else {
		io_clear_bit(LIGHT_STOP)
	}
}

// Open\close door. Does not automaticly close door
// Lift will not run while door open
// TODO: Check obstruction before closing door
func openDoor(open bool) {
	if open {
		io_set_bit(DOOR_OPEN)
		doorOpen = true
	} else {
		io_clear_bit(DOOR_OPEN)
		doorOpen = false
	}
}
