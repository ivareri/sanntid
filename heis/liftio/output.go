package liftio

import (
	"log"
	"time"
)

// Called from init and on <-*quit
func ioShutDown() {
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
}

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
func setLight(lightch chan Light) {
	lightmap := []int{
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
	keyType := []int{
		Command: -1,
		Up:      3,
		Down:    7,
		Stop:    12,
		door:    13}
	select {
	default:
		return
	case light := <-lightch:
		if light.On {
			io_set_bit(lightmap[keyType[int(light.Button)]+int(light.Floor)])
		} else {
			io_clear_bit(lightmap[keyType[int(light.Button)]+int(light.Floor)])
		}
	}
}

// Called from readFloorSensor
func setFloorLight(floor int) {
	if (floor < 1) || (floor > 4) {
		log.Fatal("Floororder out of range: ", floor)
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
