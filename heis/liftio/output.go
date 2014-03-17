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
		if !bla.direction {
			io_set_bit(MOTORDIR)
		} else {
			io_clear_bit(MOTORDIR)
		}
		time.Sleep(8 * time.Millisecond)
		io_write_analog(MOTOR, 2048+4*int(bla.speed))
	default:
		return
	}
}

// Sets order\command lights
// TODO: ugly beast. Should be a cleaner way of doing this
func setLight(lightch chan Light) {
	lightmap := []int {
		LIGHT_COMMAND1,
		LIGHT_COMMAND2,
		LIGHT_COMMAND3,
		LIGHT_COMMAND4,
		LIGHT_UP1,
		LIGHT_UP2,
		LIGHT_UP3,
		LIGHT_UP4,
		LIGHT_DOWN1,
		LIGHT_DOWN2,
		LIGHT_DOWN3,
		LIGHT_DOWN4,
		LIGHT_STOP,
		DOOR_OPEN}
	keyType := []int {
		Command: -1,
		Up: 3,
		Down: 7,
		Stop: 12,
		Door: 13}
	select {
	default:
		return
	case light := <-lightch:
		if light.On {
			io_set_bit(lightmap[keyType[int(light.Button)] + int(light.Floor)])
		} else {
			io_clear_bit(lightmap[keyType[int(light.Button)]+int(light.Floor)])
		}
		if light.Button == Door {
			doorOpen = light.On
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
