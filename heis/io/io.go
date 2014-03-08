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
	floor  uint
	button buttonType
}

type lightType int

const (
	up lightTYpe = itoa + 1
	down
	command
)

type Light struct {
	floor uint
	light lightType
	on    bool
}

type Status struct {
	running   bool
	floor     uint
	direction bool
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

func runElevator(floorOrder chan uint, floor chan Status) {
	floorSeen := make(chan uint)
	var currentFloor, floorStop uint = 0
	var lastFloor Status
	lastFloor.direction = false

	go readFloorSensor(floorSeen)
	// Go to closest floor downwards.
	// Do this to get a known state
	runMotor(DEFAULTSPEED, lastFloor.direction)
	for {
		currentFloor <- floorSeen
		if currentFloor != 0 {
			break
		}
	}
	lastFloor.floor = currentFloor
	lastFloor.running = false
	floor <- lastFloor
	runMotor(0, lastFloor.direction)

	// Elevator should be in known state. Starting loop
	for {
		select {
		case newFloorStop := <-floorOrder:
			if newFloorStop < 1 || newFloorStop > MAXFLOOR {
				log.Println("FloorOrder out of range:", newFloorStop)
			} else {
				floorStop = newFloorStop
			}
		case currentFloor <- floorSeen:
			if currentFloor == floorStop {
				runMotor(0, lastFloor.direction)
				floorStop = 0
			}
			if currentFloor > 1 && currentFloor < MAXFLOOR {
				lastFloor.floor = currentFloor
				floor <- lastFloor
			}
		default:
			if floorStop == 0 {
				break
			}
			if floorStop < lastFloor.floor {
				lastFloor.direction = false
			} else {
				lastFloor.direction = true
			}
			if !io_read_bit(DOOR_OPEN) {
				runMotor(DEFAULTSPEED, lastFloor.direction)
			}
		}
	}
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
func readButtons(keypress chan Button) {

	floor_command := [4]int{FLOOR_COMMAND1,
		FLOOR_COMMAND2,
		FLOOR_COMMAND3,
		FLOOR_COMMAND4}

	floorUp := [3]int{
		FLOOR_UP1,
		FLOOR_UP2,
		FLOOR_UP3}

	floorDown := [3]int{
		FLOOR_DOWN2,
		FLOOR_DOWN3,
		FLOOR_DOWN4}

	lastPress := create(map[int]bool)

	for {
		for i := uint; i <= MAXFLOOR; i++ {
			if io_read_bit(floorCommand[i]) && !lastPress[floorCommand[i]] {
				lastPress[floorCommand[i]] = true
				keypress <- Button{i + 1, command}
			} else if !io_read_bit(floorCommand[i]) && lastPress[floorCommand[i]] {
				lastPress[floorCommand[i]] = false
			}
		}

		for i := uint; i < MAXFLOOR; i++ {
			if io_read_bit(floorUP[i]) && !lastPress[floorUP[i]] {
				lastPress[floorUP[i]] = true
				keypress <- Button{i, up}
			} else if !io_read_bit(floorUP[i]) && lastPress[floorUP[i]] {
				lastPress[floorUP[i]] = false
			}
		}

		for i := uint; i < MAXFLOOR; i++ {
			if io_read_bit(floorDown[i]) && !lastPress[floorDown[i]] {
				lastPress[floorDown[i]] = true
				keypress <- Button{i + 2, down}
			} else if !io_read_bit(floorDown[i]) && lastPress[floorDown[i]] {
				lastPress[floorDown[i]] = false
			}
		}

		if io_read_bit(STOP) && !lastPress[STOP] {
			lastPress[STOP] = true
			keypress <- Button{0, stop}
		} else if (!io_read_bit(STOP)) && (lastPress[STOP]) {
			lastPress[STOP] = false
		}

		if io_read_bit(OBSTRUCTION) && !lastPress[OBSTRUCTION] {
			lastPress[OBSTRUCTION] = true
			keypress <- Button{0, obstruction}
		} else if !io_read_bit(OBSTRUCTION) && lastPress[OBSTRUCTION] {
			lastPress[OBSTRUCTION] = false
		}
	}
}

func setLight(light Light) {
	if light.on {
		switch light.floor {
		case 1:
			switch light.light {
			case command:
				io_set_bit(LIGHT_COMMAND1)
			case up:
				io_set_bit(LIGHT_UP1)
			}
		case 2:
			switch light.light {
			case command:
				io_set_bit(LIGHT_COMMAND2)
			case up:
				io_set_bit(LIGHT_UP2)
			case down:
				io_set_bit(LIGHT_DOWN2)
			}
		case 3:
			switch light.light {
			case command:
				io_set_bit(LIGHT_COMMAND3)
			case up:
				io_set_bit(LIGHT_UP3)
			case down:
				io_set_bit(LIGHT_DOWN3)
			}
		case 4:
			switch light.light {
			case command:
				io_set_bit(LIGHT_COMMAND4)
			case down:
				io_set_bit(LIGHT_DOWN4)
			}
		}
	} else {
		switch light.floor {
		case 1:
			switch light.light {
			case command:
				io_clear_bit(LIGHT_COMMAND1)
			case up:
				io_clear_bit(LIGHT_UP1)
			}
		case 2:
			switch light.light {
			case command:
				io_clear_bit(LIGHT_COMMAND2)
			case up:
				io_clear_bit(LIGHT_UP2)
			case down:
				io_clear_bit(LIGHT_DOWN2)
			}
		case 3:
			switch light.light {
			case command:
				io_clear_bit(LIGHT_COMMAND3)
			case up:
				io_clear_bit(LIGHT_UP3)
			case down:
				io_clear_bit(LIGHT_DOWN3)
			}
		case 4:
			switch light.light {
			case command:
				io_clear_bit(LIGHT_COMMAND4)
			case down:
				io_clear_bit(LIGHT_DOWN4)
			}
		}
	}
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
