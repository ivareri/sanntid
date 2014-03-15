package globalqueue

import (
	"liftnet/"
	"time"
)

const acceptedTimeout = 400
const newTimeout = 40

var globalQueue = make(map[int]liftnet.Message)

func checkTimeout() {
		for key, message := range globalQueue {
			if message.Status == liftnet.Done {
				toNetwork<-message
				delete(globalQueue[key])
			} else if message.Status == liftnet.New {
				timestamp = time.Now()-time.Duration(newTimeout*time.Millisecond) 
				if message.TimeRecv <= timestamp {

				} else if message.TimeRecv <= 2*timestamp {
					log.Println("2x timeout")
				} else if message.TimeRecv <= 3*timestamp {
					log.Println("3x timeout")
				}
				log.Println("New order timed out: ", message)
				newOrderTimeout(&message)
			} else if message.Status == liftnet.Accpeted && message.TimeRecv <= time.Now()-time.Duration(acceptedTimeout*time.Millisecond) {
				log.Println("Accepted order timed out: ", message)
				acceptedOrderTimeout(&message)
			}
		}
}

func newOrderTimeout(message *liftnet.Message, critical int ) {
	switch critical {
	case 3:
		takeOrder(message)
	case 2:
		if isIdle {
			takeOrder(message)
		} else if figureOfSuitability(message liftnet.Message, status FloorStatus) > 1 {
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
	*message.Status = liftnet.Accepted
	toNetwork<-*message
	//TODO: add to local queue
}

func acceptedOrderTimout(message *liftnet.Message) {
	log.Println("Some elevator didn't do as promised")
}
