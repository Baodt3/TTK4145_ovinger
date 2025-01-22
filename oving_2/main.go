package main

import (
	"fmt"
	"log"
	"net"
)

func udp() {
	s, err := net.ResolveUDPAddr("udp", ":30000")
	ln, err := net.ListenUDP("udp", s)

	if err != nil {
		log.Fatal(err)
	}

	defer ln.Close()
	buffer := make([]byte, 1024)
	n, remoteAddr, err := ln.ReadFromUDP(buffer)
	log.Printf("Received %d bytes from %v: %s", n, remoteAddr, buffer[:n])

	con, _ := net.Dial("udp", "10.100.23.204:20003")
	con.Write([]byte("HELOooo"))
	buffer2 := make([]byte, 1024)

	defer con.Close()

	s2, err2 := net.ResolveUDPAddr("udp", ":20003")

	if err2 != nil {
		log.Fatal(err2)
	}

	ln2, _ := net.ListenUDP("udp", s2)
	n, remoteAddr, _ = ln2.ReadFromUDP(buffer2)
	log.Printf("Received %d bytes from %v: %s", n, remoteAddr, buffer2[:n])
}

func main() {
	
	conn, err := net.Dial("tcp", "10.100.23.204:33546")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	buf := make([]byte, 1024)
	n, _ := conn.Read(buf)

	fmt.Println(string(buf[:n]))

	ln, _ := net.Listen("tcp", ":33546")

	data := []byte("Connect to: 10.100.23.13:33546\x00")
	_, err1 := conn.Write(data)

	conn2, _ := ln.Accept()

	if err1 != nil {
		log.Fatal(err1)
	}

	n, _ = conn2.Read(buf)

	fmt.Println(string(buf[:n]))
}
