package liftio

var lastPress [14]bool
var currentFloor = -1

// Called from RunIO.
func readButtons(keypress chan Button) {
	buttons := []int{
		FLOOR_COMMAND1,
		FLOOR_COMMAND2,
		FLOOR_COMMAND3,
		FLOOR_COMMAND4,
		FLOOR_UP1,
		FLOOR_UP2,
		FLOOR_UP3,
		FLOOR_UP4,
		FLOOR_DOWN1,
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
		Up,
		Down,
		Down,
		Down,
		Down,
		Stop,
		Obstruction}

	for index, key := range buttons {
		if readbutton(key, index) {
			keypress <- Button{uint(index%4 + 1), keyType[index]}
		}
	}
}

// Called from ReadButtons
func readbutton(key int, index int) bool {
	if io_read_bit(key) {
		if !lastPress[index] {
			lastPress[index] = true
			return true
		}
	} else if lastPress[index] {
		lastPress[index] = false
	}
	return false
}

// Called from RunIO
func readFloorSensor() {
	sensormap := []int{
		SENSOR1,
		SENSOR2,
		SENSOR3,
		SENSOR4}

	//  for {
	atfloor := false
	for i := 0; i < 4; i++ {
		if io_read_bit(sensormap[i]) {
			floorsensor(i + 1)
			atfloor = true
			return
		}
		// No floor sensors active
	}
	if !atfloor {
		floorsensor(0)
	}
	//  }
}

// Called from readFloorSensor
func floorsensor(sensor int) {
	if currentFloor != sensor {
		if sensor != 0 {
			setFloorLight(sensor)
		}
		currentFloor = sensor
		floorSeen <- uint(sensor)
	}
}
