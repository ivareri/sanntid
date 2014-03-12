package localQueue

import (
	"log"
	"encoding/json"
	"os"
)

// Adds floor to local Queue and tells network
func addLocalCommand(buttonPressed button, localQueue [][]bool){
	SetLight(Light{buttonPressed.Floor, Command, true})
	localQueue[floor-1][Command] = true 
	//TODO: Skriv til fil
}

// Deletes floor from local Queue and tells network
func deleteLocalCommand(floor uint ,localQueue [][]bool){
	SetLight(Light{floor, Command, false})
	localQueue[floor-1][Command] = false
	//TODO: skriv til fil
}

// Adds request to localQueue
func addLocalRequest(manager chan button, localQueue [][]bool){
	buttonPressed := <- manager
	SetLight(Light{buttonPressed.Floor, buttonPressed.Button, true})
	localQueue[buttonPressed.Floor-1][buttonPressed.Button-1] = true
}

// Deletes requests from localQueue
func deleteLocalRequest(Direction bool, floor uint, localQueue [][]bool){
	if Direction {
		SetLights(Light{floor, Up, false})
		localQueue[floor-1][Up-1] = false
		//TODO: Tell network
	}else{
		SetLights(Light{floor, Down, false})
		localQueue[floor-1][Down-1] = false
		//TODO: Tell network
	}
}

// Returns next floor ordered from the local queue or 0 if empty
func GetOrder(floorOrder chan uint, status chan FloorStatus, localQueue [][]bool){
	status := <- status
	currentFloor := status.floor
	currentIndex = int(currentFloor-1)
	
	if status.direction {
		if next := checkUp(currentIndex , 3, localQueue); next && currentIndex !=3 {
			floorOrder <- next
		} else if next := checkDown(3, 0, localQueue); next {
			floorOrder <- next
		} else {
			floorOrder <- checkUp(0, currentIndex, localQueue)
			}
	} else {
		if next := checkDown(currentIndex , 0, localQueue); next && currentIndex !=3 {
			floorOrder <- next
		} else if next := checkUp(0, 3, localQueue); next{
			floorOrder <- next
		} else {
			floorOrder <- checkDown(3, currentIndex, localQueue)
		}
	}
}

// Returns next floor ordered above current in UP queue or 0 if empty
func checkUp(start int, stop int, lockalQueue [][]bool) int {
	for i := start; i <= stop ; i++ {
		if localQueue[i][Up-1] || localQueue[i][Command-1] {
			return i+1
		}
	} 
	return 0
}

// Returns next floor ordered below current in DOWN queue or 0 if empty
func checkDown(start int, stop int, lockalQueue [][]bool) int {
	for i := floor; i >= stop; i--{
		if localQueue[i][Down-1] || localQueue[i][Command-1]{
			return i+1
		}
	} 
	return 0
}

func writeQueueToFile(localQueue [][]bool, filename string){
	queue, _ := json.Marshal(localQueue [][]bool)
	file, err := os.Create("localQueue")
	if err != nil {
		log.Println(err)
	} 
	n, err := file.Write(queue)
	if err != nil {
		log.Println(err)
	}
	
}

func readQueueToFile(filename string) [][]bool{
	
}
