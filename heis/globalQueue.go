package globalqueue

import (
	"liftnet/"
	"time"
)

const acceptedTimeout = 400
const newTimeout = 40

var globalQueue = make(map[int]liftnet.Message)

func newMessage(message liftnet.Message) {
	key := message.Floor
	if message.Direction {
		key += 10
	}
	if val, ok := globalQueue[key]; !ok {
		globalQueue[key] = message
	} else {
		switch message.Status {
		case liftnet.Done:
			delete(globalQueue[key])
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
	for key, message := range globalQueue {
		if message.Status == liftnet.Done {
			toNetwork <- message
			delete(globalQueue[key])
		} else if message.Status == liftnet.New {
			timestamp = time.Now() - time.Duration(newTimeout*time.Millisecond)
			if message.TimeRecv <= 3*timestamp {
				log.Println("23x timeout")
				newOrderTimeout(&message, 3)
			} else if message.TimeRecv <= 2*timestamp {
				log.Println("2x timeout")
				newOrderTimeout(&message, 2)
			} else if message.TimeRecv <= timestamp {
				log.Println("1x timeout")
				newOrderTimeout(&message, 1)
			}
			log.Println("New order timed out: ", message)
		} else if message.Status == liftnet.Accpeted {
			timestamp = time.Now() - time.Duration(acceptedTimeout*time.Millisecond)
			log.Println("Accepted order timed out: ", message)
			acceptedOrderTimeout(&message)
		}
	}
}

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

func takeOrder(message *liftnet.Message) {
	log.Println("Accepted order", message)
	*message.Id = myID
	*message.Status = liftnet.Accepted
	//TODO: add to local queue
	toNetwork <- *message
}

func acceptedOrderTimout(message *liftnet.Message) {
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
