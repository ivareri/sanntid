package main

import (
	"fmt"
	"os"
	"flag"
	"time"
	"os/exec"
	"os/signal"
	"syscall"
	"log"
)

var backupflag = flag.Bool("backup", false, "Start as backup process")

func main() {
	flag.Parse();
	fmt.Println("Backupflag is set to:", *backupflag)
	if *backupflag {
		backup()
	} else {
		spawnBackup()
		mainLoop()
	}
}

func mainLoop() {
	count()
}
func spawnBackup() {
	fmt.Println("Spawning backup process")
	cmd := exec.Command("./test", "--backup")
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	err := cmd.Start()
	if err != nil{
		log.Fatal(err)
	}
}

func backup() {
	fmt.Println("Backup started");
	ch := make(chan os.Signal, 1)
        signal.Notify(ch, syscall.SIGUSR1)
loop:
	for {
		select {
		case signal := <-ch:
			log.Println("Got signal, ", signal)
			break loop
		default:
			log.Println("No signal, sleeping")
			time.Sleep(200 * time.Millisecond)
		}
	}
	fmt.Println("Main proccess died");
	spawnBackup()
	mainLoop()
}

func count() {
	for i := 0; i <= 5; i++ {
		fmt.Println(i)
		time.Sleep(200 * time.Millisecond)
	}
}
