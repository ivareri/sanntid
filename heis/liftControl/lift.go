package liftControl

import (
	"../liftio"
	"../liftnet"
	"../localQueue"
	"log"
	"math/rand"
	"time"
)

var myID int
var myPenalty int
var isIdle = true
var lastOrder = uint(0)
var toNetwork = make(chan liftnet.Message, 10)
var fromNetwork = make(chan liftnet.Message, 10)
var floorOrder = make(chan uint, 5) // floor orders to io
var setLight = make(chan liftio.Light, 5)
var liftStatus liftio.LiftStatus
var maxFloor = liftio.MAXFLOOR
var quit *chan bool
func RunLift(quit *chan bool) {

	buttonPress := make(chan liftio.Button, 5) // button presses from io
	status := make(chan liftio.LiftStatus, 5)  // the lifts status
	rand.Seed(time.Now().Unix())
	myPenalty = rand.Intn(100)
	myID = liftnet.NetInit(&toNetwork, &fromNetwork, quit)
	liftio.IOInit(&floorOrder, &setLight, &status, &buttonPress)
	restoreBackup()
	liftStatus = <-status
	ticker1 := time.NewTicker(10 * time.Millisecond).C
	ticker2 := time.NewTicker(5 * time.Millisecond).C
	log.Println("Up and running. My is is: ", myID, "Penalty time is: ", myPenalty)
	for {
		select {
		case button := <-buttonPress:
			newKeypress(button)
		case liftStatus = <-status:
			runQueue()
		case message := <-fromNetwork:
			newMessage(message)
			orderLight(message)
		case <-ticker1:
			checkTimeout()
		case <-ticker2:
			runQueue()
		case <-*quit:
			return
		}
	}
}

// Called byLiftStatus RunLift
func newKeypress(button liftio.Button) {
	switch button.Button {
	case liftio.Up:
		log.Println("Request up button pressed.", button.Floor)
		addMessage(button.Floor, true)
		setOrderLight(button.Floor, true, true)
	case liftio.Down:
		log.Println("Request down button ressed.", button.Floor)
		addMessage(button.Floor, false)
		setOrderLight(button.Floor, false, true)
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

}

// Called by RunLift
func runQueue() {
	floor := liftStatus.Floor
	if liftStatus.Running {
		if liftStatus.Direction {
			floor++
		} else {
			floor--
		}
	}
	order, direction := localQueue.GetOrder(floor, liftStatus.Direction)
	if liftStatus.Floor == order && liftStatus.Door {
		removeFromQueue(order, direction) 
		lastOrder = 0
		liftStatus.Door = true
		time.Sleep(20 * time.Millisecond)
	} else if order == 0 && !liftStatus.Door {
		isIdle = true
	} else if order != 0 {
		isIdle = false
		if lastOrder != order && !liftStatus.Door {
			lastOrder = order
			floorOrder <- order
		}
	}
}

// Called by runQueue
func removeFromQueue(floor uint, direction bool) {
	log.Println("Removing from queue", floor, direction)
	localQueue.DeleteLocalOrder(floor, direction)
	delMessage(floor, direction)
	setLight <- liftio.Light{floor, liftio.Command, false}
	setOrderLight(floor, direction, false)

}

// Called by RunLift
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

// Called by orderLight
func setOrderLight(floor uint, direction bool, on bool) {
	if direction {
		setLight <- liftio.Light{floor, liftio.Up, on}
	} else {
		setLight <- liftio.Light{floor, liftio.Down, on}
	}
}

// Called by RunLift and ReadQueuFromFile
func addCommand(floor uint) {
	localQueue.AddLocalCommand(floor)
	setLight <- liftio.Light{floor, liftio.Command, true}
}

// Called by RunLift
func restoreBackup() {
	for i, val := range localQueue.ReadQueueFromFile() {
		if val {
			addCommand(uint(i + 1))
		}
	}
}
