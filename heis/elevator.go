package elevatorControl

import (
	"./liftio"
	"./localQueue"
	"log"
	"os"
	"time"
)

// Does excatly what it says
func elevatorControl() {
	
	floorOrder := make(chan uint)     		// floor orders to io
	buttonPress := make(chan liftio.Button) // button presses from io
	status := make(chan liftio.FloorStatus) // the lifts status
	setLight := make(chan liftio.Light)		 
	toNetwork := make(chan liftnet.Message)   
	fromNetwork := make(chan liftnet.Message) 
	
	ReadQueueFromFile() // if no prev queue: "queue.txt doesn't exitst" entered in log
	init(&floorOrder, &status)
	MultiCastInit(&toNetwork, &fromNetwork)
	go runIO(&buttonPress)

	var orderedFloor uint = 1 // Not sure if this is the way to go
	for {
		select {
		case button := <-buttonPress:
			if button.Button == Up {
				log.Println("Request button %v pressed.", button.Button)
				addMessage(button.Floor, true)
			} else if button.Button == Down {
				log.Println("Request button %v pressed.", button.Button)
				addMessage(button.Floor, true)
			} else if button.Button == Command {
				log.Println("Command button %v pressed.", button.Floor)
				addLockalCommand(button, localQueue.Queue)
			} else if buttonPressed.button == Stop {
				log.Println("Stop button pressed")
				// action optional
			} else if buttonPressed.button == Obstuction {
				log.Println("Obstruction")
				// action optional
			}
		case floorReached := <-status:
			orderedFloor = GetOrder(floorReached.Floor, floorReached.Direction)
			// could lead to faulty behaviour, (if so: check more frequent)
			floorOrder <- orderedFloor
			if floorReached.Floor == orderedFloor {
				setLight <- Light{0, Door, true}
				time.Sleep(time.Second * 3)
				setLight <- Light{0, Door, false}
				DeleteLocalOrder(floorReached.Floor, FloorReached.Direction)
				delMessage(floorReached.Floor, floorReached.Direction)
				orderedFloor = GetOrder(floorReached.Floor, floorReached.Direction)
				floorOrder <-orderedFloor
			}
		case message := <-fromNetwork:
			if message.Status == liftnet.Done {
				if message.Direction {
					setLight <- Light{message.Floor, Up, false}
				} else {
					setLight <- Light{message.Floor, Down, false}
				}
			} else if message.Status == liftnet.New {
				if message.Direction {
					setLight <- Light{message.Floor, Up, true}
				} else {
					setLight <- Light{message.Floor, Down, true}
				}
			} else if message.Status == liftnet.Accepted {

			}
			messageManager(message)
		default:
			log.Println("No Action at all.")
			// teit å skrive  ut hele tiden?
		}
	}

}
