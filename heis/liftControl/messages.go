package liftControl

import (
	"../liftnet"
	"../localQueue"
	"log"
	"time"
)

const acceptedTimeoutBase = 4
const newTimeoutBase = 500
const reassignTimeoutBase = 500 

var globalQueue = make(map[uint]liftnet.Message)

func generateKey(floor uint, direction bool) uint {
	if direction {
		floor += 10
	}
	return floor
}

// Called by newkeyPress
func addMessage(floor uint, direction bool) {
	key := generateKey(floor, direction)
	message := liftnet.Message{
		LiftId:    myID,
		Floor:     floor,
		Direction: direction,
		Status:    liftnet.New,
		Weigth:    figureOfSuitability(floor, direction),
		TimeRecv:  time.Now()}

	if _, ok := globalQueue[key]; ok {
		return
	}
	globalQueue[key] = message
	toNetwork <- message
}

// Called by removeFromQueue
func delMessage(floor uint, direction bool) {
	key := generateKey(floor, direction)
	if val, ok := globalQueue[key]; ok {
		val.Status = liftnet.Done
		toNetwork <- val
		delete(globalQueue, key)
	}
}

// Called by RunLift
func newMessage(message liftnet.Message) {
	key := generateKey(message.Floor, message.Direction)
	val, inQueue := globalQueue[key]
	if inQueue {
		switch message.Status {
		case liftnet.Done:
			delete(globalQueue, key)
		case liftnet.Accepted:
			globalQueue[key] = message
		case liftnet.New:
			if val.Weigth <= message.Weigth {
				globalQueue[key] = message
			}
		case liftnet.Reassign:
			if message.ReassId != myID{
				if val.Weigth <= message.Weigth {
					globalQueue[key] = message
				}
			}
		default:
			log.Println("Unknown status recived: ", message.Status, ". Ignoring message")
		}
	} else {
		switch message.Status {
		case liftnet.Done:
			// TODO: Promptly ignore? (write something more sensible here)
		case liftnet.Accepted:
			log.Println("Old message accepted by lift: ", message.LiftId)
			if val.Status == liftnet.Reassign && val.ReassId == myID { 
				localQueue.DeleteLocalRequest(message.Floor, message.Direction)
			}
			globalQueue[key] = message
		case liftnet.Reassign:
			fs := figureOfSuitability(message.Floor, message.Direction)
			if fs > message.Weigth {
				message.Weigth = fs
				message.LiftId = myID
				globalQueue[key] = message
				toNetwork <- message
				log.Println("Reassign, my fs i better now ", fs)
			} else {
				globalQueue[key] = message
				log.Println("I'm not fast enough, My fs: ", fs, " best fs: ", message.Weigth)
			}
		case liftnet.New:
			fs := figureOfSuitability(message.Floor, message.Direction)
			if fs > message.Weigth {
				message.Weigth = fs
				message.LiftId = myID
				globalQueue[key] = message
				toNetwork <- message
				log.Println("I'm best so far ", fs)
			} else {
				globalQueue[key] = message
			}
			log.Println("My fs:", fs, " best fs: ", message.Weigth)
		default:
			log.Println("Unknown status recived: ", message.Status, ". Ignoring message")
		}
	}
}

// Called by RunLift
func checkTimeout() {
	newTimeout := time.Duration(newTimeoutBase + myPenalty)
	acceptedTimeout := time.Duration(acceptedTimeoutBase)
	reassignTimeout := time.Duration(reassignTimeoutBase + myPenalty)
	for key, val := range globalQueue {
		// TODO : make func for all if else statements
		if val.Status == liftnet.New {
			timediff := time.Now().Sub(val.TimeRecv)
			if timediff > ((3 * newTimeout) * time.Millisecond) {
				newOrderTimeout(key, 3)
			} else if timediff > ((2 * newTimeout) * time.Millisecond) {
				newOrderTimeout(key, 2)
			} else if timediff > ((1 * newTimeout) * time.Millisecond) {
				newOrderTimeout(key, 1)
			}
		} else if val.Status == liftnet.Accepted && val.LiftId != myID {
			timediff := time.Now().Sub(val.TimeRecv)
			if timediff > ((4 * acceptedTimeout) * time.Second) {
				acceptedOrderTimeout(key, 3)
			} else if timediff > ((3 * acceptedTimeout) * time.Second) {
				acceptedOrderTimeout(key, 2)
			} else if timediff > ((2 * acceptedTimeout) * time.Second) {
				acceptedOrderTimeout(key, 1)
			}
		} else if val.Status == liftnet.Accepted && val.LiftId == myID {
			timediff := time.Now().Sub(val.TimeRecv)
			if timediff > (acceptedTimeout * time.Second) {
				val.Weigth = figureOfSuitability(val.Floor, val.Direction)
				val.TimeRecv = time.Now()
				globalQueue[key] = val
				toNetwork <- globalQueue[key]
			}
		} else if val.Status == liftnet.Reassign{
			timediff := time.Now().Sub(val.TimeRecv)
			if timediff > ((3 * reassignTimeout)* time.Millisecond){
				newOrderTimeout(key, 3)
			} else if timediff > ((2 * newTimeout) * time.Millisecond) {
				newOrderTimeout(key, 2)
			} else if timediff > ((1 * newTimeout) * time.Millisecond) {
				newOrderTimeout(key, 1)
			}
		}
	}
}

// Called by checkTimeout
func newOrderTimeout(key, critical uint) {
	switch critical {
	case 3:
		takeOrder(key)
	case 2:
		if isIdle {
			log.Println("Lift is idle, timout 2x")
			takeOrder(key)
		} else if figureOfSuitability(globalQueue[key].Floor, globalQueue[key].Direction) > globalQueue[key].Weigth {
			takeOrder(key)
		}
	case 1:
		if globalQueue[key].LiftId == myID {
			takeOrder(key)
		}
	}
}

// Called by checkTimeout
func acceptedOrderTimeout(key uint, critical uint) {
	log.Println("Some lift didn't do as promised")
	switch critical {
	case 3:
		log.Println("Something went horribly wrong")
		takeOrder(key)
	case 2:
		takeOrder(key)
	case 1:
		reassignOrder(key)
	}
}


// Called by timeout functions
func takeOrder(key uint) {
	if val, ok := globalQueue[key]; !ok {
		log.Println("Trying to accept order not in queue")
	} else {
		log.Println("Accepted order", globalQueue[key])
		val.LiftId = myID
		val.Status = liftnet.Accepted
		val.TimeRecv = time.Now()
		localQueue.AddLocalRequest(val.Floor, val.Direction)
		globalQueue[key] = val
		toNetwork <- globalQueue[key]
	}
}

// TODO: should be called from somewhere suitable
func reassignOrder(key uint){
	if val, ok := globalQueue[key]; !ok{
		log.Println("trying to reassign order not in queue")
	} else {
		log.Println("Reassigning order", globalQueue[key])
		val.Status = liftnet.Reassign
		val.ReassId = myID
		val.Weigth = figureOfSuitability(val.Floor, val.Direction) //necessary to recalculate fs here?
		val.TimeRecv = time.Now()
		globalQueue[key] = val
		toNetwork <- globalQueue[key]
	}
}

// Called by NewMessage, addMessage and newOrderTimeout
// Nearest Car algorithm, returns Figure of Suitability
func figureOfSuitability(reqFlr uint, reqDir bool) int {
	statFlr := liftStatus.Floor
	statDir := liftStatus.Direction
	if isIdle {
		if reqFlr == statFlr {
			return 6
		} else {
			return maxFloor + 1 - diff(reqFlr, statFlr)
		}
	} else if reqDir == statDir {
		// lift moving towards the requested floor and the request is in the same direction
		if (statDir && reqFlr > statFlr) || (!statDir && reqFlr < statFlr) {
			return maxFloor + 1 - diff(reqFlr, statFlr)
		}
	} else {
		// lift moving towards the requested floor, but the request is in oposite direction
		if (statDir && reqFlr > statFlr) || (!statDir && reqFlr < statFlr) {
			return maxFloor - diff(reqFlr, statFlr)
		}
	}
	return 1
}

// Called by figureOfSuitabillity
func diff(a, b uint) int {
	x := int(a)
	y := int(b)
	c := x - y
	if c < 0 {
		return c * -1
	} else {
		return c
	}
}
