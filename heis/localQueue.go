package localQueue
// Returns next floor ordered from the local queue or 0 if empty
func getOrder(status chan struct,localQueue [][]bool) int {
	status := <-status
	currentFloor := status.floor
	currentIndex = int(currentFloor-1)
	
	if status.direction {
		if next := checkUp(currentIndex , 3, localQueue); next && currentIndex !=3 {
			return next
		}else if next := checkDown(3, 0, localQueue); next{
			return next
		}else{
			return checkUp(0, currentIndex, localQueue)
			}
		}
	}else{
		if next := checkDown(currentIndex , 0, localQueue); next && currentIndex !=3 {
			return next
		}else if next := checkUp(0, 3, localQueue); next{
			return next
		}else{
			return checkDown(3, currentIndex, localQueue)
			}
		}
	}
}

// Returns next floor ordered above current in UP queue or 0 if empty
func checkUp(start int, stop int, lockalQueue [][]bool) int {
	for i := floor; i <= stop ; i++ {
		if localQueue[i][1] || localQueue[i][2]{
			return i+1
			// TODO: Delete from localqueue and tell queuemanager
		}
	} return 0
}

// Returns next floor ordered below current in DOWN queue or 0 if empty
func checkDown(start int, stop int, lockalQueue [][]bool) int {
	for i := floor; i >= stop; i--{
		if localQueue[i][0] || localQueue[i][2]{
			return i+1
			// TODO: Delete from localqueue and tell queuemanager
		}
	} return 0
}
