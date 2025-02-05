package main

import (
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"os/exec"
	"strconv"
	"time"
)

func backup() {
	fmt.Println("Starting backup")
	s, err := net.ResolveUDPAddr("udp", "127.0.0.1:20003")
	ln, err := net.ListenUDP("udp", s)

	if err != nil {
		log.Fatal(err)
	}

	defer ln.Close()

	buffer := make([]byte, 1024)
	primaryAlive := true
	var lastNumber int

	for primaryAlive {
		ln.SetReadDeadline(time.Now().Add(3 * time.Second))
		n, _, err := ln.ReadFromUDP(buffer)

		if err != nil {
			//log.Fatal(err)
			if errors.Is(err, os.ErrDeadlineExceeded) {
				primaryAlive = false
			}
		} else {
			lastNumber, err = strconv.Atoi(string(buffer[:n]))
		}

	}

	ln.Close()
	transformToPrimary(lastNumber)
}

func transformToPrimary(lastNumber int) {
	exec.Command("gnome-terminal", "--", "go", "run", "main.go").Run()
	fmt.Println("I am promoted:)")

	con, _ := net.Dial("udp", "127.0.0.1:20003")
	defer con.Close()
	//go backup()
	count := lastNumber + 1
	for {
		log.Println(count)
		con.Write([]byte(strconv.Itoa(count)))
		count++
		time.Sleep(time.Second)
	}

}

func main() {

	backup()
}
