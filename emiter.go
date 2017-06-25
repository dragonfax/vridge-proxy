package main

import (
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"net"
)

func connectToServerProxy() *net.TCPConn {
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", serverIP, 3278))
	if err != nil {
		log.Fatal(err)
	}

	log.Println("client is connected to server")

	return conn.(*net.TCPConn)
}

func startEmiter(udpIP string) {

	udpAddr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", udpIP, UDP_PORT))
	if err != nil {
		log.Fatal(err)
	}

	// Read from TCP and emit packets into UDP or TCP
	go func() {
		buf := make([]byte, 1024*10)

		log.Println("emitter has started")

		for {
			n := readFromProxy(buf)
			buf = buf[:n]

			// send UDP
			n, err := udpConn.WriteTo(buf, udpAddr)
			if err != nil {
				log.Fatal(err)
			}
			if n != len(buf) {
				log.Fatal("wrong length")
			}
		}
	}()
}

func readFromProxy(buf []byte) int {

	buf = buf[:2]
	n, err := io.ReadFull(proxyReader, buf)
	if err != nil {
		log.Fatal(err)
	}
	if n != 2 {
		log.Fatal("wrong length")
	}
	pl := int(binary.LittleEndian.Uint16(buf))

	if pl > 0 {

		buf = buf[:pl]
		n, err = io.ReadFull(proxyReader, buf)
		if err != nil {
			log.Fatal(err)
		}
		if n != pl {
			log.Fatal("wrong length")
		}
		buf = buf[:n]
	}

	return pl
}
