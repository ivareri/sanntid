package elevatorControl

import (
	"../liftnet"
	"../localQueue"
	"log"
	"time"
)

const acceptedTimeoutBase = 3
const newTimeoutBase = 500

var globalQueue = make(map[uint]liftnet.Message)

// TODO: Find appropriate channel size
var timeoutch = make(chan uint, 100)

func generateKey(floor uint, direction bool) uint {
	if direction {
		floor += 10
	}
	return floor
}

func addMessage(floor uint, direction bool) {
	key := generateKey(floor, direction)
	message := liftnet.Message{
		LiftId:        myID,
		Floor:     floor,
		Direction: direction,
		Status:    liftnet.New,
		Weigth:	figureOfSuitability(floor, direction),
		TimeRecv:  time.Now()}

	if _, ok := globalQueue[key]; ok {
		log.Println("Order already in queue")
		return
	}
	globalQueue[key] = message
	toNetwork <- message
	log.Println("Sent message, ", message)
}
func delMessage(floor uint, direction bool) {
	key := generateKey(floor, direction)
	if val, ok := globalQueue[key]; ok {
		val.Status = liftnet.Done
		toNetwork <- val
		delete(globalQueue, key)
	}
}

// Called by RunElevator
func newMessage(message liftnet.Message) {
	log.Println("Recv new message", message)
	key := generateKey(message.Floor, message.Direction)
	val, inQueue := globalQueue[key]
	if inQueue {
		switch message.Status {
		case liftnet.Done:
			delete(globalQueue, key)
		case liftnet.Accepted:
			if val.Status != liftnet.Accepted {
				globalQueue[key] = message
			} else {
				log.Println("Got new accept message for already accepted order.")
				log.Println("Check timings on elevators: ", val.LiftId, message.LiftId)
			}
		case liftnet.New:
			if  val.Weigth <  message.Weigth {
				globalQueue[key] = message
			}
		default:
			log.Println("Unknown status recived: ", message.Status, ". Ignoring message")
		}
	} else {
		switch message.Status {
		case liftnet.Done:
			// Promptly ignore?
		case liftnet.Accepted:
			log.Println("Old message acceppted by lift: ", message.LiftId)
				// Might be wise to recalculate and find better lift
			globalQueue[key] = message
		case liftnet.New:
			fs := figureOfSuitability(message.Floor, message.Direction)
			if fs >= message.Weigth {
				message.Weigth = fs
				message.LiftId = myID
				globalQueue[key] = message
				toNetwork <- message
				log.Println("I'm best so far", fs)
			} else {
				globalQueue[key] = message
			}
			log.Println("My fs:", fs, " best fs", message.Weigth)
		default:
			log.Println("Unknown status recived: ", message.Status, ". Ignoring message")
		}
	}
}

// Called by Run Elevator
func checkTimeout() {
	newTimeout := time.Duration(newTimeoutBase + (myID%10)*10)
	acceptedTimeout := time.Duration(acceptedTimeoutBase + (myID / 10))
	for key, val := range globalQueue {
		if val.Status == liftnet.New {
			timediff := time.Now().Sub(val.TimeRecv)
			if timediff > ((3 * newTimeout) * time.Millisecond) {
				log.Println("3x timeout")
				newOrderTimeout(key, 3)
			} else if timediff > ((2 * newTimeout) * time.Millisecond) {
				log.Println("2x timeout")
				newOrderTimeout(key, 2)
			} else if timediff > ((1 * newTimeout) * time.Millisecond) {
				log.Println("1x timeout")
				newOrderTimeout(key, 1)
			}
		} else if val.Status == liftnet.Accepted {
			timediff := time.Now().Sub(val.TimeRecv)
			if timediff > ((4 * acceptedTimeout) * time.Second) {
				log.Println("2x accepted timeout")
				acceptedOrderTimeout(key, 3)
			} else if timediff > ((3 * acceptedTimeout) * time.Second) {
				log.Println("1.5x accepted timeout")
				acceptedOrderTimeout(key, 2)
			} else if timediff > ((2 * acceptedTimeout) * time.Second) {
				log.Println("1x accepted timeout")
				acceptedOrderTimeout(key, 1)
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
	log.Println("Some elevator didn't do as promised")
	switch critical {
	case 3:
		log.Println("Something went horribly wrong")
		takeOrder(key)
	case 2:
		takeOrder(key)
	case 1:
		// TODO: isIdle not working??
		log.Println("Idle", isIdle)
		if isIdle {
			takeOrder(key)
		}
	}
}

// Called by timeout functions
func takeOrder(key uint) {
	log.Println("Accepted order", globalQueue[key])
	msg := globalQueue[key] // TODO: Make pretty
	msg.LiftId = myID
	msg.Status = liftnet.Accepted
	msg.TimeRecv = time.Now()
	localQueue.AddLocalRequest(globalQueue[key].Floor, globalQueue[key].Direction)
	globalQueue[key] = msg
	toNetwork <- globalQueue[key]
}

// Called by NewMessage, addMessage and newOrderTimeout
// Nearest Car algorithm, returns Figure of Suitability
func figureOfSuitability(reqFlr uint, reqDir bool) int {
	statFlr := liftStatus.Floor
	statDir := liftStatus.Direction
	if reqDir == statDir {
		// lift moving towards the requested floor and the request is in the same direction
		if (statDir && reqFlr > statFlr) || (!statDir && reqFlr < statFlr) {
			return maxFloor + 1 - diff(reqFlr, statFlr)
		}
	} else {
		// lift moving towards the requested floor, but the request is in oposite direction
		if (statDir && reqFlr > statFlr) || (!statDir && reqFlr < statFlr) {
			return maxFloor - diff(reqFlr, statFlr)
		} else {
			return 1
		}
	}
	return 0
}

// Called from figureOfSuitabillity
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
