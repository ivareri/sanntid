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
	"net"
	"strconv"
)

var backupflag = flag.Bool("backup", false, "Start as backup process")
var countflag = flag.Int("count", 0, "Where to start couting")

func main() {
	flag.Parse();
	if *backupflag {
		backup(*countflag)
	} else {
		spawnBackup(0)
		mainLoop(0)
	}
}

func mainLoop(count int) {
	number := make(chan int)
	time.Sleep(200 * time.Millisecond)
	go netSend(number)
	counter(number, count)
}

func spawnBackup(count int) {
	log.Println("Spawning new backup process")
	cntr := "-count=" + strconv.Itoa(count)
	cmd := exec.Command("./test", "-backup", cntr)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	err := cmd.Start()
	if err != nil{
		log.Fatal(err)
	}
}

func netSend(number chan int) {
	serverAddr, err := net.ResolveUDPAddr("udp", "localhost:9999")
	if err != nil {
		log.Fatal("Error: ", err)
	}
	con, err := net.DialUDP("udp", nil, serverAddr)
	if err != nil {
		log.Fatal("Error: ", err)
	}
	for {
		num :=  <-number
		buf := []byte(strconv.Itoa(num))
		_, err := con.Write(buf)
		if err != nil {
			log.Fatal("error writing data", err)
		}
	}
}

func netBackup(number chan int, quit chan bool) {
	udpAddr, err := net.ResolveUDPAddr("udp", "localhost:9999")
	if err != nil {
		log.Fatal("Error finding address: ", err)
	}
	netlisten, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		log.Fatal("Error opening port:", err);
	}
	buffer := make([]byte, 1024)
	oob := make([]byte, 1024)
	defer netlisten.Close()
	for {
		select{
		case <-quit:
			log.Println("Listener shutting down")
			return
		default:
			netlisten.SetDeadline(time.Now().Add(1e9))
			n, _, _, _, err := netlisten.ReadMsgUDP(buffer, oob)
			if err != nil {
				if opErr, ok := err.(*net.OpError); ok && opErr.Timeout() {
					continue
				} else {
					log.Fatal("Error:", err)
				}
			} else {
				num, err := strconv.Atoi(string(buffer[:n]))
				if err != nil {
					log.Fatal("Error casting buffer to int: ", err)
				}
				number <- num
			}
		}
	}
}

func backup(count int) {
	number := make(chan int)
	quit := make(chan bool)
	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, syscall.SIGINT)
	go netBackup(number, quit)
	missed_packet := 0
	var num int
loop:
	for {
		select {
		case bla := <-number:
			missed_packet = 0
			num = bla
			time.Sleep(200 * time.Millisecond)
		case <-sigint:
			log.Println("Backup caught SIGINT, ignore")
		default:
			if missed_packet == 3 {
				quit <- true
				spawnBackup(num+1)
				mainLoop(num+1)
				break loop
			}
			missed_packet += 1
			time.Sleep(200 * time.Millisecond)
		}
	}
}

func counter(number chan int, count int) {
	i := count
	for  {
		fmt.Println(i)
		number <-i
		time.Sleep(200 * time.Millisecond)
		i++
	}
}
