package localQueue

import (
	"encoding/json"
	"log"
	"os"
)

type Queue struct {
	Up      [4]bool
	Down    [4]bool
	Command [4]bool
}

const BackupFile = "backupQueue.q"

var localQueue = Queue{}

// Called by elevatorControl
// Write localQueue.command to backup file
func writeQueueToFile() {
	commandQueue, err := json.Marshal(localQueue.Command)
	if err != nil {
		log.Println(err)
	}
	file, err := os.Create(BackupFile)
	if err != nil {
		log.Println("Error in opening file ", err)
	}
	n, err := file.Write(commandQueue) // overwrites existing file
	if err != nil {
		log.Println("Error in writing to file ", err)
	}
	log.Println("wrote, ", n, " bytes to ", BackupFile)
}

// Called by elevatorControl
// Add command to local Queue and writes to backup file
func AddLocalCommand(floor uint) {
	localQueue.Command[floor-1] = true
	writeQueueToFile()
}

// Called by elevatorControl
// Add request to localQueue
func AddLocalRequest(floor uint, direction bool) {
	if direction {
		localQueue.Up[floor-1] = true
	} else {
		localQueue.Down[floor-1] = true
	}
}

// Called by elevatorControl
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

// Called by elevatorControl
// Returns next floor ordered from the local queue or 0 if empty
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
		if nextStop := checkDown(currentFloor, 1); nextStop > 0  {
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
	for i := int(start) - 1; i <= int(stop) - 1; i++ {
		if i > 3 || i < 0 {
			log.Println("Out of bounds UP. Stop: ", stop, " start: ", start, " i: ", i)
			return 0
		} else if localQueue.Up[i] || localQueue.Command[i] {
			return uint(i+1)
		}
	}
	return 0
}

// Called by GetOrder()
// Returns floor for next floor order below current in Down queue or 0 if empty
func checkDown(start uint, stop uint) uint {
	for i := int(start)-1; i >= int(stop)-1; i-- {
		if i > 3 || i < 0 {
			log.Println("Out of bounds Down. Stop: ", stop, " start: ", start, " i: ", i)
			return 0
		} else if localQueue.Down[i] || localQueue.Command[i] {
			return uint(i+1)
		}
	}
	return 0
}
