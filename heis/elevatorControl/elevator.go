package elevatorControl

import (
	"../liftnet"
	"../liftio"
	"../localQueue"
	"log"
	"time"
)

var myID int
var isIdle bool
var toNetwork = make(chan liftnet.Message, 10)
var fromNetwork = make(chan liftnet.Message, 10)

// Does excatly what it says
func Run() {

	floorOrder := make(chan uint, 5)           // floor orders to io
	buttonPress := make(chan liftio.Button, 5) // button presses from io
	status := make(chan liftio.FloorStatus, 5) // the lifts status
	setLight := make(chan liftio.Light, 5)

	localQueue.ReadQueueFromFile() // if no prev queue: "queue.txt doesn't exitst" entered in log
	liftio.Init(&floorOrder, &setLight, &status, &buttonPress)
	addr, iface, err := liftnet.FindIP()
	if err != nil {
		log.Println("Error finding interface", err)
	}
	go liftnet.MulticastInit(toNetwork, fromNetwork, iface)
	myID = liftnet.FindID(addr)
	var orderedFloor uint = 1 // Not sure if this is the way to go
	floorReached := <-status
	for {
		select {
		case button := <-buttonPress:
			if button.Button == liftio.Up {
				log.Println("Request up button pressed.", button.Floor)
				addMessage(button.Floor, true)
			} else if button.Button == liftio.Down {
				log.Println("Request down button ressed.", button.Floor)
				addMessage(button.Floor, true)
			} else if button.Button == liftio.Command {
				log.Println("Command button %v pressed.", button.Floor)
				localQueue.AddLocalCommand(button)
			} else if button.Button == liftio.Stop {
				log.Println("Stop button pressed")
				// action optional
			} else if button.Button == liftio.Obstruction {
				log.Println("Obstruction")
				// action optional
			}
		case floorReached = <-status:
			log.Println("Passing floor: ", floorReached.Floor)
		default:
			order := localQueue.GetOrder(floorReached.Floor, floorReached.Direction)
			if floorReached.Floor != 0 && floorReached.Floor == order {
				setLight <- liftio.Light{0, liftio.Door, true}
				time.Sleep(time.Second * 3)
				setLight <- liftio.Light{0, liftio.Door, false}
				localQueue.DeleteLocalOrder(floorReached.Floor, floorReached.Direction)
				delMessage(floorReached.Floor, floorReached.Direction)
				orderedFloor = localQueue.GetOrder(floorReached.Floor, floorReached.Direction)
				if orderedFloor != 0 {
					isIdle = false
					floorOrder <- orderedFloor
				} else {
					isIdle =true
				}
			} else if order != 0 {
				floorOrder <-order
			}
			time.Sleep(10 * time.Millisecond)
		case message := <-fromNetwork:
			if message.Status == liftnet.Done {
				if message.Direction {
					setLight <- liftio.Light{message.Floor, liftio.Up, false}
				} else {
					setLight <- liftio.Light{message.Floor, liftio.Down, false}
				}
			} else if message.Status == liftnet.New {
				if message.Direction {
					setLight <- liftio.Light{message.Floor, liftio.Up, true}
				} else {
					setLight <- liftio.Light{message.Floor, liftio.Down, true}
				}
			} else if message.Status == liftnet.Accepted {

			}
			messageManager(message)
		}
	}

}
