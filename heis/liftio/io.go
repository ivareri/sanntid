package liftio

import (
	"log"
	"time"
)

const MAXFLOOR = 4
const DEFAULTSPEED = 200

type ButtonType int

var (
	floorSeen chan uint
	motor chan int
	// TODO: fix this
	dir chan bool
	lastPress [12]bool
)
const (
	Up ButtonType = iota
	Down
	Command
	Stop
	Obstruction
)

// Floor should be ignored when Button is  Stop or Obstruction
// Used for passing around keypresses
type Button struct {
	Floor  uint
	Button ButtonType
}

type LightType int

// Used for setting command and order lights
// TODO: Fix type name
type Light struct {
	Floor uint
	Light ButtonType
	On    bool
}

// Running is false when lift is stationary
type FloorStatus struct {
	Running   bool
	Floor     uint
	Direction bool
}

// Initilazes elevator, starts runElevator routine, wich in turn starts readFloorSensor routine.
// Do not write or read from floorOrder and floor untill this function returns true
func Init(floorOrder *chan uint, floor *chan FloorStatus, button *chan Button) bool {
	// Init hardware
	if !io_init() {
		log.Fatal("Error during HW init")
	}
//	motor<-0
//	dir<-false
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
	go RunIO(button)
	go runElevator(floorOrder, floor)
	// wait for eleavtor to arrive at a floor
	<-*floor
	return true
}

func RunIO(button *chan Button) {
	i := 0
	for {
		if i == 20  {
			log.Println("Looped 20 times")
			i = 0
		}
		readButtons(button)
		readFloorSensor()
//		runMotor()
		i++
	}
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
//	runMotor(DEFAULTSPEED, lastFloor.Direction)
//	for {
//		currentFloor = <-floorSeen
//		if currentFloor != 0 {
//			log.Println("Found floor")
//			break
//		}
//	}
//	runMotor(0, lastFloor.Direction)
//	lastFloor.Floor = currentFloor
//	lastFloor.Running = false
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
				motor<-0
				dir<-lastFloor.Direction
				lastFloor.Floor = currentFloor
				lastFloor.Running = false
				*floor <- lastFloor
			}
			if currentFloor == floorStop {
				motor<-0
				dir<-lastFloor.Direction
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
			if !io_read_bit(DOOR_OPEN) && !lastFloor.Running {
				motor<-DEFAULTSPEED
				dir<-lastFloor.Direction
				lastFloor.Running = true
				*floor <- lastFloor
			}
		}
	}
}

// Called from runElevator
func runMotor() {
	// Invert direction in order to break elevator before stopping
	select {
	case speed :=<-motor:
				direction := <-dir
				if speed == 0 {
					direction = !direction
				}
			if direction {
				io_set_bit(MOTORDIR)
			} else {
				io_clear_bit(MOTORDIR)
			}
			time.Sleep(5*time.Millisecond)
			io_write_analog(MOTOR, 2048+4*int(speed))
	default:
			return
		}
}

// Sets order\command lights
// TODO: ugly beast. Should be a cleaner way of doing this
func SetLight(light Light) {
	if light.On {
		switch light.Floor {
		case 1:
			switch light.Light {
			case Command:
				io_set_bit(LIGHT_COMMAND1)
			case Up:
				io_set_bit(LIGHT_UP1)
			}
		case 2:
			switch light.Light {
			case Command:
				io_set_bit(LIGHT_COMMAND2)
			case Up:
				io_set_bit(LIGHT_UP2)
			case Down:
				io_set_bit(LIGHT_DOWN2)
			}
		case 3:
			switch light.Light {
			case Command:
				io_set_bit(LIGHT_COMMAND3)
			case Up:
				io_set_bit(LIGHT_UP3)
			case Down:
				io_set_bit(LIGHT_DOWN3)
			}
		case 4:
			switch light.Light {
			case Command:
				io_set_bit(LIGHT_COMMAND4)
			case Down:
				io_set_bit(LIGHT_DOWN4)
			}
		}
	} else {
		switch light.Floor {
		case 1:
			switch light.Light {
			case Command:
				io_clear_bit(LIGHT_COMMAND1)
			case Up:
				io_clear_bit(LIGHT_UP1)
			}
		case 2:
			switch light.Light {
			case Command:
				io_clear_bit(LIGHT_COMMAND2)
			case Up:
				io_clear_bit(LIGHT_UP2)
			case Down:
				io_clear_bit(LIGHT_DOWN2)
			}
		case 3:
			switch light.Light {
			case Command:
				io_clear_bit(LIGHT_COMMAND3)
			case Up:
				io_clear_bit(LIGHT_UP3)
			case Down:
				io_clear_bit(LIGHT_DOWN3)
			}
		case 4:
			switch light.Light {
			case Command:
				io_clear_bit(LIGHT_COMMAND4)
			case Down:
				io_clear_bit(LIGHT_DOWN4)
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
func doorOpen(open bool) {
	if open {
		io_set_bit(DOOR_OPEN)
	} else {
		io_clear_bit(DOOR_OPEN)
	}
}

// Called from RunIO.
func readButtons(keypress *chan Button) {
	buttons := []int{
		FLOOR_COMMAND1,
		FLOOR_COMMAND2,
		FLOOR_COMMAND3,
		FLOOR_COMMAND4,
		FLOOR_UP1,
		FLOOR_UP2,
		FLOOR_UP3,
		FLOOR_DOWN2,
		FLOOR_DOWN3,
		FLOOR_DOWN4,
		STOP,
		OBSTRUCTION}

	keyType := []ButtonType{
		Command,
		Command,
		Command,
		Command,
		Up,
		Up,
		Up,
		Down,
		Down,
		Down,
		Stop,
		Obstruction}

//	for {
		for index, key := range buttons{
			if readbutton(key, &lastPress[index]) {
				log.Println("Keypress", key)
				*keypress <-Button{uint(index + 1), keyType[index]}
			}
		}
//	}
}

// Called from ReadButtons
func readbutton(key int, lastPress *bool) bool {
	if io_read_bit(key) {
		if !*lastPress {
			*lastPress = true
			return true
		}
	} else if *lastPress {
		*lastPress = false
	}
	return false
}

// Called from RunIO
func readFloorSensor() {
	currentFloor := -1
	sensormap := []int{
		SENSOR1,
		SENSOR2,
		SENSOR3,
		SENSOR4}

//	for {
		atfloor := false
		for i := 0; i < 4; i++ {
			if io_read_bit(sensormap[i]) {
				floorsensor(i+1, &currentFloor)
				atfloor = true
				break
			}
			// No floor sensors active
		}
		if !atfloor {
			floorsensor(0, &currentFloor)
		}
//	}
}
// Called from readFloorSensor
func floorsensor(sensor int, currentFloor *int) {
	if *currentFloor != sensor {
		if sensor != 0 {
			setFloorLight(sensor)
		}
		*currentFloor = sensor
		log.Println("Floor ", sensor, ":", currentFloor)
		floorSeen <- uint(sensor)
	}
}
