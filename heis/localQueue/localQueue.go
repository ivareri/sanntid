package localQueue

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

type OrderQueue struct {
	Up      [4]bool // Requests
	Down    [4]bool // Requests
	Command [4]bool // Commands
}

const backupFile = "backupQueue.q"

var localQueue = OrderQueue{}

// Called by liftControl
// Writes localQueue.command to backup file
func writeQueueToFile() {
	commandQueue, err := json.Marshal(localQueue.Command)
	if err != nil {
		log.Println(err)
	}
	err = ioutil.WriteFile(backupFile, commandQueue, 0600)
	if err != nil {
		log.Println("Error wirting to file", err)
	}
}

// Called by restoreBackup
// Returns bool struct with commands saved before shutdown
func ReadQueueFromFile() []bool {
	byt, err := ioutil.ReadFile(backupFile)
	if err != nil {
		log.Println("Error reading from backupfile", err)
	}
	var cmd []bool
	if err := json.Unmarshal(byt, &cmd); err != nil {
		log.Println("Error during unmarshal: ", err)
		log.Println("Got: ", cmd)
	}
	return cmd
}

// Called by liftControl
// Adds command to local Queue and writes to backup file
func AddLocalCommand(floor uint) {
	localQueue.Command[floor-1] = true
	writeQueueToFile()
}

// Called by liftControl
// Adds request to localQueue
func AddLocalRequest(floor uint, direction bool) {
	if direction {
		localQueue.Up[floor-1] = true
	} else {
		localQueue.Down[floor-1] = true
	}
}

// TODO: Called by ...
// Deletes requests reassigned to other lifts from localQueue
func DeleteLocalRequest(floor uint, Direction bool){
	if Direction{
		localQueue.Up[floor-1] = false
	} else {
		localQueue.Down[floor-1] = false
	}
}

// Called by liftControl
// Deletes orders from localQueue and writes to backup file
func DeleteLocalOrder(floor uint, Direction bool) {
	localQueue.Command[floor-1] = false
	writeQueueToFile()
	if Direction {
		localQueue.Up[floor-1] = false
	} else {
		localQueue.Down[floor-1] = false
	}
}

// Called by liftControl
// Returns next floor ordered from localQueue, 0 if empty
// and bool indicating that order is above/below currentFloor
func GetOrder(currentFloor uint, direction bool) (uint, bool) {
	if direction {
		if nextStop := checkUp(currentFloor, 4); nextStop > 0 {
			return nextStop, true
		} else if nextStop := checkDown(4, 1); nextStop > 0 {
			return nextStop, false
		} else {
			return checkUp(1, 4), true
		}
	} else {
		if nextStop := checkDown(currentFloor, 1); nextStop > 0 {
			return nextStop, false
		} else if nextStop := checkUp(1, 4); nextStop > 0 {
			return nextStop, true
		} else {
			return checkDown(4, 1), false
		}
	}
}

// Called by GetOrder()
// Returns floor for next order above current in Up queue or 0 if empty
func checkUp(start uint, stop uint) uint {
	for i := int(start) - 1; i <= int(stop)-1; i++ {
		if i > 3 || i < 0 {
			log.Println("Out of bounds UP. Stop: ", stop, " start: ", start, " i: ", i)
			return 0
		} else if localQueue.Up[i] || localQueue.Command[i] {
			return uint(i + 1)
		}
	}
	return 0
}

// Called by GetOrder()
// Returns floor for next floor order below current in Down queue or 0 if empty
func checkDown(start uint, stop uint) uint {
	for i := int(start) - 1; i >= int(stop)-1; i-- {
		if i > 3 || i < 0 {
			log.Println("Out of bounds Down. Stop: ", stop, " start: ", start, " i: ", i)
			return 0
		} else if localQueue.Down[i] || localQueue.Command[i] {
			return uint(i + 1)
		}
	}
	return 0
}
