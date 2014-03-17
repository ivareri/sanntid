package elevatorControl

import (
	"../liftnet"
	"../localQueue"
	"log"
	"time"
)

const acceptedTimeoutBase = 3
const newTimeoutBase = 50

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
		Id:        myID,
		Floor:     floor,
		Direction: direction,
		Status:    liftnet.New,
		TimeSent:  time.Now(),
		TimeRecv:  time.Now()}

	if _, ok := globalQueue[key]; ok {
		log.Println("Order already in queue")
		return
	} else if isIdle {
		message.Status = liftnet.Accepted
		localQueue.AddLocalRequest(floor, direction)
	}
	globalQueue[key] = message
	orderLight(message)
	toNetwork <- message
}
func delMessage(floor uint, direction bool) {
	key := generateKey(floor, direction)
	if val, ok := globalQueue[key]; ok {
		val.Status = liftnet.Done
		toNetwork <- val
		delete(globalQueue, key)
	}
}

func newMessage(message liftnet.Message) {
	key := generateKey(message.Floor, message.Direction)
	if val, ok := globalQueue[key]; !ok {
		if val.Status == liftnet.Done {
			return
		}
		globalQueue[key] = message
	} else {
		switch message.Status {
		case liftnet.Done:
			delete(globalQueue, key)
		case liftnet.Accepted:
			if val.Status != liftnet.Accepted {
				globalQueue[key] = message
			} else {
				log.Println("Got new accept message for already accepted order.")
				log.Println("Check timings on elevators: ", val.Id, message.Id)
			}
		case liftnet.New:
			log.Println("Recived new order, but order already in queue")
		default:
			log.Println("Unknown status recived")
		}
	}
}

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
			if timediff > ((3 * acceptedTimeout) * time.Second) {
				log.Println("3x accepted timeout")
				acceptedOrderTimeout(key, 3)
			} else if timediff > ((2 * acceptedTimeout) * time.Second) {
				log.Println("2x accepted timeout")
				acceptedOrderTimeout(key, 2)
			} else if timediff > ((1 * acceptedTimeout) * time.Second) {
				log.Println("1x accepted timeout")
				acceptedOrderTimeout(key, 1)
			}
		}
	}
}

//called form checktimout
func newOrderTimeout(key, critical uint) {
	switch critical {
	case 3:
		takeOrder(key)
	case 2:
		if isIdle {
			takeOrder(key)
		} else if figureOfSuitability(globalQueue[key], true, 1) > 1 {
			takeOrder(key)
		}
	case 1:
		if isIdle {
			takeOrder(key)
		}
	}
}

//called from checkTimout
func acceptedOrderTimeout(key uint, critical uint) {
	log.Println("Some elevator didn't do as promised")
	switch critical {
	case 3:
		log.Println("Something went horribly wrong")
	case 2:
		takeOrder(key)
	case 1:
		if isIdle {
			takeOrder(key)
		}
	}
}

// called from timeouts
func takeOrder(key uint) {
	log.Println("Accepted order", globalQueue[key])
	msg := globalQueue[key] // TODO: Make pretty
	msg.Id = myID
	msg.Status = liftnet.Accepted
	msg.TimeRecv = time.Now()
	localQueue.AddLocalRequest(globalQueue[key].Floor, globalQueue[key].Direction)
	globalQueue[key] = msg
	toNetwork <- globalQueue[key]
}

// Nearest Car algorithm, returns Figure of Suitability
// Lift with largest FS should accept the request
func figureOfSuitability(request liftnet.Message, statDir bool, statFlr uint) int {
	MAXFLOOR := 4 // TODO: make pretty
	reqDir := request.Direction
	reqFlr := request.Floor
	if reqDir == statDir {
		// lift moving towards the requested floor and the request is in the same direction
		if (statDir && reqFlr > statFlr) || (!statDir && reqFlr < statFlr) {
			return MAXFLOOR + 1 - diff(reqFlr, statFlr)
		}
	} else {
		// lift moving towards the requested floor, but the request is in oposite direction
		if (statDir && reqFlr > statFlr) || (!statDir && reqFlr < statFlr) {
			return MAXFLOOR - diff(reqFlr, statFlr)
		} else {
			return 1
		}
	}
	return 0
}

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
