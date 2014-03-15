package localQueue

import (
	"encoding/json"
	"log"
	"os"
)


// Writes localQueue.command to file for backup
func writeQueueToFile(localQueue Queue) {
	cmdQueue, err := json.Marshal(localQueue.Command)
	if err != nil{
		log.Println(err)
	}
	file, err := os.Create("localQueue.txt")  
	if err != nil {
		log.Println("Error in opening file ", err)
	}
	n, err := file.Write(cmdQueue) // overwrites existing file 
	if err != nil {
		log.Println("Error in writing to file ", err)
	}
	log.Println("wrote %d bytes\n to localQueue.txt", n)
	defer file.Close()
}

// Reads localQueue.Command from backup file
func ReadQueueFromFile(localQueue Queue){
	input, err := os.Open("queue.txt")
	if err != nil {
		log.Println("Error in opening file: ", err)
	}
	byt := make([]byte, 23)
	dat, err := input.Read(byt)
	if err != nil {
		log.Println("Error in reading file: ", err)
	}
	defer input.Close()
	log.Println("Read %d bytes: %s from file\n", dat, string(byt))
	if err :=json.Unmarshal(d, &localQueue.commandQueue); err != nil {
		log.Println(err)
	}
}

// Adds floor to local Queue and writes to file
func AddLocalCommand(buttonPressed button, localQueue Queue) {
	SetLight(Light{buttonPressed.Floor, Command, true})
	localQueue.Command[buttonPressed.Floor-1] = true
	writeQueueToFile(localQueue, "localQueue")
}

// Deletes floor from local Queue and writes to file
func DeleteLocalCommand(floor uint, localQueue Queue) {
	SetLight(Light{floor, Command, false})
	localQueue.Command[floor-1] = false
	writeQueueToFile(localQueue, "localQueue")
}

// Adds request to localQueue and writes to file
func AddLocalRequest(manager chan button, localQueue Queue) {
	buttonPressed := <-manager
	SetLight(Light{buttonPressed.Floor, buttonPressed.Button, true})
	if buttonPressed.button == Up {
		localQueue.Up[buttonPressed.Floor] = true
	} else {
		localQueue.Down[buttonPressed.Floor] = true
	}
	
}

// Deletes requests from localQueue and writes to file
// Is direction ok to use to decide on which queue or does the braking fuck it up? 
func DeleteLocalRequest(Direction bool, floor uint, localQueue Queue) {
	if Direction {
		SetLights(Light{floor, Up, false})
		localQueue.Up[floor-1] = false
		writeQueueToFile(localQueue, "localQueue")
	} else {
		SetLights(Light{floor, Down, false})
		localQueue.Down[floor-1] = false
		writeQueueToFile(localQueue, "localQueue")
	}
}

// Returns next floor ordered from the local queue or 0 if empty
func GetOrder(floorOrder chan uint, status chan FloorStatus, localQueue Queue) {
	status := <-status
	currentFloor := status.floor
	currentIndex = int(currentFloor - 1)

	if status.direction {
		if next := checkUp(currentIndex, 3, localQueue); next && currentIndex != 3 {
			floorOrder <- next
		} else if next := checkDown(3, 0, localQueue); next {
			floorOrder <- next
		} else {
			floorOrder <- checkUp(0, currentIndex, localQueue)
		}
	} else {
		if next := checkDown(currentIndex, 0, localQueue); next && currentIndex != 3 {
			floorOrder <- next
		} else if next := checkUp(0, 3, localQueue); next {
			floorOrder <- next
		} else {
			floorOrder <- checkDown(3, currentIndex, localQueue)
		}
	}
}

// Returns next floor ordered above current in UP queue or 0 if empty
func checkUp(start int, stop int, lockalQueue Queue) int {
	for i := start; i <= stop; i++ {
		if localQueue.Up[i] || localQueue.Command[i] {
			return i + 1
		}
	}
	return 0
}

// Returns next floor ordered below current in DOWN queue or 0 if empty
func checkDown(start int, stop int, lockalQueue Queue) int {
	for i := floor; i >= stop; i-- {
		if localQueue.Down[i] || localQueue.Command[i] {
			return i + 1
		}
	}
	return 0
}





