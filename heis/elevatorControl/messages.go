package elevatorControl

import (
	"../liftnet"
	"../localQueue"
	"log"
	"time"
)

const acceptedTimeout = 1000
const newTimeout = 100

var globalQueue = make(map[uint]liftnet.Message)

// TODO: Find appropriate channel size
var timeoutch = make(chan uint, 100)

func generateKey(floor uint, direction bool) uint {
	if direction {
		floor += 10
	}
	return floor
}

func timeout(key uint, duration time.Duration) {
	timer := time.NewTimer(duration)
	<-timer.C
	log.Println("Timer expired, key: ", key)
	timeoutch <- key
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

	if _, ok := globalQueue[key]; !ok {
		globalQueue[key] = message
		orderLight(message)
		toNetwork <- message
		go timeout(key, newTimeout*time.Millisecond)
	} else {
		log.Println("Order already in queue")
	}
}
func delMessage(floor uint, direction bool) {
	key := generateKey(floor, direction)
	if val, ok := globalQueue[key]; !ok {
		log.Println("Trying to remove nonexsiting message from queue")
	} else {
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
		} else if val.Status == liftnet.New {
			go timeout(key, newTimeout *time.Millisecond)
		} else if val.Status == liftnet.Accepted {
			go timeout(key, acceptedTimeout *time.Millisecond)
		}
		globalQueue[key] = message
	} else {
		switch message.Status {
		case liftnet.Done:
			delete(globalQueue, key)
		case liftnet.Accepted:
			if val.Status != liftnet.Accepted {
				globalQueue[key] = message
				go timeout(key, acceptedTimeout * time.Millisecond)
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
	select {
	default:
		return
	case key := <-timeoutch:

		if val, ok := globalQueue[key]; !ok {
			return
		} else {
			if val.Status == liftnet.Done {
				toNetwork <- val
				delete(globalQueue, key)
			} else if val.Status == liftnet.New {
				log.Println("1x timeout")
				newOrderTimeout(key, 3)
			} else if val.Status == liftnet.Accepted {
				log.Println("Accepted order timed out: ", val)
				acceptedOrderTimeout(key)
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
func acceptedOrderTimeout(key uint) {
	log.Println("Some elevator didn't do as promised")
	critical := 1
	switch critical {
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
