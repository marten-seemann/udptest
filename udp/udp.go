package main

import (
	"log"
	"net"
	"time"
)

const (
	rnet       = "udp"
	addrString = "127.0.0.1:4242"
	dataSize   = 100
	nSends     = 1000
)

var (
	conn *net.UDPConn
)

func main() {
	addr, err := net.ResolveUDPAddr(rnet, addrString)
	if err != nil {
		panic(err)
	}
	conn, err = net.DialUDP(rnet, nil, addr)
	if err != nil {
		panic(err)
	}

	// go sendPackets()
	for i := 0; true; i++ {
		log.Printf("Attempt #%d\n", i)
		runOnce()
	}
}

func runOnce() {
	c := make(chan struct{})
	addr, err := net.ResolveUDPAddr(rnet, addrString)
	if err != nil {
		panic(err)
	}
	ln, err := net.ListenUDP(rnet, addr)
	if err != nil {
		panic(err)
	}
	defer ln.Close()

	go func() {
		data := make([]byte, dataSize)
		log.Println("ReadFrom called")
		_, _, err := ln.ReadFrom(data)
		if err != nil {
			log.Printf("ReadFrom errored: %s\n", err.Error())
			panic(err)
		}
		close(c)
	}()

	conn, err := net.DialUDP(rnet, nil, addr)
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	for i := 0; i < nSends; i++ {
		data := make([]byte, dataSize)
		_, err = conn.Write(data)
		if err != nil {
			log.Printf("Sending errored: %s\n", err.Error())
			panic(err)
		}
		log.Println("Sent packet")
		time.Sleep(time.Millisecond)
		select {
		case <-c:
			return
		default:
		}
	}
	log.Println("FAIL")
	panic("FAIL")
}

// func runOnce2() {
// 	addr, err := net.ResolveUDPAddr(rnet, addrString)
// 	if err != nil {
// 		panic(err)
// 	}
// 	ln, err := net.ListenUDP(rnet, addr)
// 	if err != nil {
// 		panic(err)
// 	}
// 	defer ln.Close()

// 	data := make([]byte, dataSize)
// 	_, _, err = ln.ReadFrom(data)
// 	if err != nil {
// 		panic(err)
// 	}
// }

// func sendPackets() {
// 	addr, err := net.ResolveUDPAddr(rnet, addrString)
// 	if err != nil {
// 		panic(err)
// 	}
// 	conn, err := net.DialUDP(rnet, nil, addr)
// 	if err != nil {
// 		panic(err)
// 	}
// 	data := make([]byte, dataSize)
// 	for {
// 		_, err = conn.Write(data)
// 		// if err != nil {
// 		// panic(err)
// 		// }
// 		// log.Println("Sent packet")
// 		// time.Sleep(time.Microsecond)
// 	}
// }
