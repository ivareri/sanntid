package elevatorControl

import (
	"../liftio"
	"../liftnet"
	"../localQueue"
	"log"
	"time"
	"encoding/json"
	"os"
)

var myID int
var isIdle bool
var lastOrder = uint(0)
var toNetwork = make(chan liftnet.Message, 10)
var fromNetwork = make(chan liftnet.Message, 10)
var floorOrder = make(chan uint, 5)           // floor orders to io
var setLight = make(chan liftio.Light, 5)
// Does excatly what it says
func Run() {

	buttonPress := make(chan liftio.Button, 5) // button presses from io
	status := make(chan liftio.FloorStatus, 5) // the lifts status

	myID = liftnet.Init(&toNetwork, &fromNetwork)
	liftio.Init(&floorOrder, &setLight, &status, &buttonPress)
	readQueueFromFile() // if no prev queue: "queue.txt doesn't exitst" entered in log
	floorReached := <-status
	ticker1 := time.NewTicker(10 * time.Millisecond).C
	ticker2 := time.NewTicker(5 * time.Millisecond).C
	for {
		select {
		case button := <-buttonPress:
			switch button.Button {
			case liftio.Up:
				log.Println("Request up button pressed.", button.Floor)
				addMessage(button.Floor, true)
			case liftio.Down:
				log.Println("Request down button ressed.", button.Floor)
				addMessage(button.Floor, false)
			case liftio.Command:
				log.Println("Command button pressed.", button.Floor)
				addCommand(button.Floor)
			case liftio.Stop:
				log.Println("Stop button pressed")
				// action optional
			case liftio.Obstruction:
				log.Println("Obstruction")
				// action optional
			}
		case floorReached = <-status:
			log.Println("Passing floor: ", floorReached.Floor)
		case message := <-fromNetwork:
			newMessage(message)
			orderLight(message)
		case <-ticker1:
			checkTimeout()
		case <-ticker2:
			runQueue(floorReached)
		}
	}
}

// Called from run loop
func runQueue(floorReached liftio.FloorStatus) {
	floor := floorReached.Floor
	if floorReached.Running {
		if floorReached.Direction {
			floor++
		} else {
			floor--
		}
	}
	order, direction := localQueue.GetOrder(floor, floorReached.Direction)
	if order == 0  {
		return
	}
	if floorReached.Floor == order && !floorReached.Running {
		removeFromQueue(order, direction)
		lastOrder = 0
		go openDoor()
		isIdle = true
	} else {
		isIdle = false
		if lastOrder != order {
			lastOrder = order
			floorOrder <- order
		}
	}
}
func removeFromQueue(floor uint, direction bool) {
	log.Println("Removing from queue", floor, direction)
	localQueue.DeleteLocalOrder(floor, direction)
	delMessage(floor, direction)
	setLight <- liftio.Light{floor, liftio.Command, false}
	setOrderLight(floor, direction, false)

}
// Called from runQueue
// io wrapper makes sure lift is stationary when door open
func openDoor() {
	log.Println("open door")
	setLight <-liftio.Light{0, liftio.Door, true}
	time.Sleep(time.Second * 3)
	log.Println("close door")
	setLight <-liftio.Light{0, liftio.Door, false}
}

// called from run loop and netsomething
func orderLight(message liftnet.Message) {
	switch message.Status {
	case liftnet.Done:
		setOrderLight(message.Floor, message.Direction, false)
	case liftnet.New:
		setOrderLight(message.Floor, message.Direction, true)
	case liftnet.Accepted:
		setOrderLight(message.Floor, message.Direction, true)
	}
}

// called from orderLight
func setOrderLight(floor uint, direction bool, on bool) {
	if direction {
		setLight <- liftio.Light{floor, liftio.Up, on}
	} else {
		setLight <- liftio.Light{floor, liftio.Down, on}
	}
}

// called by run and ReadQueuFromFile
func addCommand(floor uint) {
	localQueue.AddLocalCommand(floor)
	setLight <- liftio.Light{floor, liftio.Command, true}
}

// Called by run
func readQueueFromFile() {
        input, err := os.Open(localQueue.BackupFile)
        if err != nil {
                log.Println("Error in opening file: ", err)
                return
        }
        defer input.Close()
        byt := make([]byte, 23)
        dat, err := input.Read(byt)
        if err != nil {
                log.Println("Error in reading file: ", err)
                return
        }
        log.Println("Read ", dat, " bytes from file ")
	var cmd []bool
        if err := json.Unmarshal(byt, &cmd); err != nil {
                log.Println(err)
        }
	for i, val := range(cmd) {
		if val {
			addCommand(uint(i+1))
		}
	}
}
