package main

import (
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"net"
)

func connectToServerProxy() *net.TCPConn {
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", SERVER_PUBLIC_IP, 3278))
	if err != nil {
		log.Fatal(err)
	}

	log.Println("client is connected to server")

	return conn.(*net.TCPConn)
}

func startEmiter(udpTargetIP string) {

	udpTargetAddr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", udpTargetIP, UDP_PORT))
	if err != nil {
		log.Fatal(err)
	}

	// Read from TCP and emit packets into UDP or TCP
	go func() {
		buf := make([]byte, 1024*10)

		log.Println("emitter has started")

		for {
			n := readFromProxy(buf)

			if n == 0 {
				panic("zero length packet from proxy")
			}

			buf = buf[:n]

			// log.Println("emitting packet, size ", n)

			// send UDP
			n, err := udpConn.WriteTo(buf, udpTargetAddr)
			if err != nil {
				log.Println(buf)
				log.Println(udpTargetAddr)
				panic(err)
			}
			if n != len(buf) {
				panic("udp: wrong length")
			}
		}
	}()
}

func readFromProxy(buf []byte) int {

	buf = buf[:2]
	n, err := io.ReadFull(proxyReader, buf)
	if err != nil {
		log.Fatal("emitter: ", err)
	}
	if n != 2 {
		log.Fatal("emitter: wrong length")
	}
	pl := int(binary.LittleEndian.Uint16(buf))

	if pl == 0 {
		panic("zero length packet")
	}

	buf = buf[:pl]
	n, err = io.ReadFull(proxyReader, buf)
	if err != nil {
		log.Fatal("emitter: ", err)
	}
	if n != pl {
		log.Fatal("emitter: wrong length")
	}
	buf = buf[:n]

	return pl
}
