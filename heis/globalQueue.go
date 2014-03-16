package elevatorControl

import (
	"./liftnet"
	"./localQueue"
	"log"
	"time"
)

const acceptedTimeout = 400
const newTimeout = 40

var globalQueue = make(map[int]liftnet.Message)
var timeoutch = make(chan int, 100)

func generateKey(floor int, direction bool) int {
	if direction {
		floor += 10
	}
	return floor
}

func timeout(key int, duration time.Duration) {
	timer := time.NewTimer(duration)
	<-timer.C
	log.Println("Timer expired, key: ", key)
	timeoutch <-key
}
func addMessage(floor int, direction bool) {
	key := generateKey(floor, direction)
	message := liftnet.Message{
		Id:        myID,
		Floor:     floor,
		Direction: direction,
		Status:    liftnet.New,
		TimeSent:  time.Now(),
		TimeRecv:  time.Now()}

	if val, ok := globalQueue[key]; !ok {
		globalQueue[key] = message
		toNetwork <- message
		go timeout(key, newTimeout * time.Millisecond)
	} else {
		log.Println("Order already in queue")
	}
}
func delMessage(floor int, direction bool) {
	key := generateKey(floor, direction)
	if val, ok := globalQueue[key]; !ok {
		log.Println("Trying to remove nonexsiting message from queue")
	} else {
		val.Status = liftnet.Done
		toNetwork<-val
		delete(globalQueue, key)
	}
}

func messageManager(message liftnet.Message) {
	key := generateKey(message.Floor, message.Direction)
	if val, ok := globalQueue[key]; !ok {
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
	select {
	default:
		return
	case key := <-timeoutch:

		if val, ok := globalQueue[key]; !ok {
			return
		} else {
			if val.Status == liftnet.Done {
				toNetwork <-val
				delete(globalQueue, key)
			} else if val.Status == liftnet.New {
				log.Println("1x timeout")
				newOrderTimeout(&val, 1)
			} else if val.Status == liftnet.Accepted {
				log.Println("Accepted order timed out: ", val)
				acceptedOrderTimeout(&val)
			}
		}
	}
}
//called form checktimout
func newOrderTimeout(message *liftnet.Message, critical int) {
	switch critical {
	case 3:
		takeOrder(message)
	case 2:
		if isIdle {
			takeOrder(message)
	} else if figureOfSuitability(message, status) > 1 {
			takeOrder(message)
		}
	case 1:
		if isIdle {
			takeOrder(message)
		}
	}
}

//called from checkTimout
func acceptedOrderTimeout(message *liftnet.Message) {
	log.Println("Some elevator didn't do as promised")
	critical := 1
	switch critical {
	case 2:
		takeOrder(message)
	case 1:
		if isIdle {
			takeOrder(message)
		}
	}
}

// called from timeouts
func takeOrder(message *liftnet.Message) {
	log.Println("Accepted order", message)
	*message.Id = myID
	*message.Status = liftnet.Accepted
	addLocalRequest(*message.Floor, *message.Direction)
	toNetwork <- *message
}

// Nearest Car algorithm, returns Figure of Suitability
// Lift with largest FS should accept the request
func figureOfSuitability(request liftnet.Message, status FloorStatus) int {
	reqDir := request.Direction
	reqFlr := request.Floor
	statDir := status.Direction
	statFlr := status.Floor
	if reqDir == statDir {
		// lift moving towards the requested floor and the request is in the same direction
		if (statDir && reqFlr > statFlr) || (!statDir && reqFlr < statFlr) {
			fs := MAXFLOOR + 1 - diff(reqFlr, statFlr)
		}
	} else {
		// lift moving towards the requested floor, but the request is in oposite direction
		if (statDir && reqFlr > statFlr) || (!statDir && reqFlr < statFlr) {
			fs := MAXFLOOR - diff(reqFlr, statFlr)
		} else {
			fs := 1
		}
	}
	return fs
}
