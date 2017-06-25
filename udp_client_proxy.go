package main

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"log"
	"net"
)

func client() {
	initClientTCPPorts()

	proxy = connectToServerProxy()
	proxyReader = bufio.NewReader(proxy)
	startEmiter(clientIP)

	CreateUDPPort(localBindIP, UDP_PORT, proxy)
}

var udpConn *net.UDPConn

func CreateUDPPort(bindIP string, port int, tcpProxy net.Conn) {

	bindAddr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", bindIP, port))
	if err != nil {
		log.Fatal(err)
	}
	log.Println("binding to UDP ", bindAddr)

	udpConn, err = net.ListenUDP("udp", bindAddr)
	if err != nil {
		log.Fatal(err)
	}

	// Catch UDP, tunnel into TCP
	go func() {

		// needs to be global, not goroutine local
		buffer := make([]byte, 1024*10)

		log.Println("upp port is bound, and processing")

		for {
			buffer = buffer[:cap(buffer)]
			n, sAddr, err := udpConn.ReadFrom(buffer)
			if err != nil {
				panic(err)
			}
			buffer = buffer[:n]

			sourceAddr := sAddr.(*net.UDPAddr)

			if sourceAddr == nil {
				// log.Println("empty packet")
				continue
			}

			writeToProxy(buffer)
		}
	}()
}

func writeToProxy(buf []byte) {

	bs := make([]byte, 2)
	binary.LittleEndian.PutUint16(bs, uint16(len(buf)))
	n, err := proxy.Write(bs)
	if err != nil {
		log.Fatal("emitter: ", err)
	}
	if n != 2 {
		log.Fatal("emitter: didn't write the full buffer length header")
	}

	n, err = proxy.Write(buf)
	if err != nil {
		log.Fatal("emitter: ", err)
	}
	if n != len(buf) {
		log.Fatal("emitter: didn't write the full buffer")
	}

}
