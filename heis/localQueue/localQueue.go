package localQueue

import (
	"../liftio"
	"encoding/json"
	"log"
	"os"
)

type Queue struct {
	Up      [4]bool
	Down    [4]bool
	Command [4]bool
}

const backupFile = "backupQueue.q"

var localQueue = Queue{}

// Called by elevatorControl
// Write localQueue.command to backup file
func writeQueueToFile() {
	commandQueue, err := json.Marshal(localQueue.Command)
	if err != nil {
		log.Println(err)
	}
	file, err := os.Create(backupFile)
	if err != nil {
		log.Println("Error in opening file ", err)
	}
	n, err := file.Write(commandQueue) // overwrites existing file
	if err != nil {
		log.Println("Error in writing to file ", err)
	}
	log.Println("wrote %d bytes\n to %q", n, backupFile)
	defer file.Close()
}

// Called by elevatorControl
// Read Command queue from backup file
func ReadQueueFromFile() {
	input, err := os.Open(backupFile)
	if err != nil {
		log.Println("Error in opening file: ", err)
		return
	}
	byt := make([]byte, 23)
	dat, err := input.Read(byt)
	if err != nil {
		log.Println("Error in reading file: ", err)
		return
	}
	defer input.Close()
	log.Println("Read %d bytes: %s from file\n", dat, string(byt))
	if err := json.Unmarshal(byt, &localQueue.Command); err != nil {
		log.Println(err)
	}
}

// Called by elevatorControl
// Add command to local Queue and writes to backup file
func AddLocalCommand(buttonPressed liftio.Button) {
	localQueue.Command[buttonPressed.Floor-1] = true
	writeQueueToFile()
}

// Called by elevatorControl
// Add request to localQueue
func AddLocalRequest(floor uint, direction bool) {
	if direction {
		localQueue.Up[floor] = true
	} else {
		localQueue.Down[floor] = true
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
func GetOrder(currentFloor uint, direction bool) uint {
	currentIndex := int(currentFloor - 1)
	if direction {
		if nextStop := checkUp(currentIndex, 3, localQueue); nextStop > 0 && currentIndex != 3 {
			return nextStop
		} else if nextStop := checkDown(3, 0, localQueue); nextStop > 0 {
			return nextStop
		} else {
			return checkUp(0, currentIndex, localQueue)
		}
	} else {
		if nextStop := checkDown(currentIndex, 0, localQueue); nextStop > 0 && currentIndex != 3 {
			return nextStop
		} else if nextStop := checkUp(0, 3, localQueue); nextStop > 0 {
			return nextStop
		} else {
			return checkDown(3, currentIndex, localQueue)
		}
	}
}

// Called by GetOrder()
// Returns next floor ordered above current in Up queue or 0 if empty
func checkUp(start int, stop int, lockalQueue Queue) uint {
	for i := start; i <= stop; i++ {
		if localQueue.Up[i] || localQueue.Command[i] {
			return  uint(i + 1)
		}
	}
	return 0
}

// Called by GetOrder()
// Returns next floor ordered below current in Down queue or 0 if empty
func checkDown(start int, stop int, lockalQueue Queue) uint {
	for i := start; i >= stop; i-- {
		if localQueue.Down[i] || localQueue.Command[i] {
			return uint(i + 1)
		}
	}
	return 0
}
