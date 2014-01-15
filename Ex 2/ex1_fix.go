package main

import (
    . "fmt"     // Using '.' to avoid prefixing functions with their package names
    . "runtime" //   This is probably not a good idea for large projects...
)

var i = 0
var lock = make(chan int, 1)
var masterlock = make(chan int, 2)
func adder() {
    for x := 0; x < 1000000; x++ {
	<-lock
        i++
	lock <- 1
    }
   masterlock <- 1
}

func subtract() {
     for x := 0; x < 1000000; x++ {
        <-lock
        i--
        lock <- 1
    }
   masterlock <- 1
}

func main() {
    GOMAXPROCS(NumCPU())        // I guess this is a hint to what GOMAXPROCS does...
    lock <- 1
    go adder()                  // This spawns adder() as a goroutine
    go subtract()
    for x := 0; x < 50; x++ {
        Println(i)
    }
    <-masterlock
    <-masterlock
    // No way to wait for the completion of a goroutine (without additional syncronization)
    // We'll come back to using channels in Exercise 2. For now: Sleep
    Println("Done:", i);
}
