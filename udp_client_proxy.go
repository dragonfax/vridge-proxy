package main

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"log"
	"net"
)

func client() {
	proxy = connectToServerProxy()
	proxyReader = bufio.NewReader(proxy)

	startEmiter(clientIP)
	initTCPPorts()
	CreateUDPPort(localBindIP, UDP_PORT, proxy, clientIP)
}

var udpConn net.UDPConn

func CreateUDPPort(bindIP string, port int, tcpProxy net.Conn, udpIP string) {

	bindAddr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", bindIP, port))
	if err != nil {
		log.Fatal(err)
	}
	log.Println("binding to UDP ", bindAddr)

	udpConn, err := net.ListenUDP("udp", bindAddr)
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
				log.Fatal(err)
			}
			buffer = buffer[:n]

			sourceAddr := sAddr.(*net.UDPAddr)

			if sourceAddr == nil {
				// log.Println("empty packet")
				continue
			}

			writeToProxy(buffer, sourceAddr.Port, 0)
		}
		log.Println("closing udp proxy")
	}()
}

func writeToProxy(buf []byte, port int, id int) {

	bs := make([]byte, 2)
	binary.LittleEndian.PutUint16(bs, uint16(len(buf)))
	n, err := proxy.Write(bs)
	if err != nil {
		log.Fatal(err)
	}
	if n != 2 {
		log.Fatal("didn't write the full buffer length header")
	}

	binary.LittleEndian.PutUint16(bs, uint16(port))
	n, err = proxy.Write(bs)
	if err != nil {
		log.Fatal(err)
	}
	if n != 2 {
		log.Fatal("didn't write the full port address header")
	}

	binary.LittleEndian.PutUint16(bs, uint16(id))
	n, err = proxy.Write(bs)
	if err != nil {
		log.Fatal(err)
	}
	if n != 2 {
		log.Fatal("didn't write the full id header")
	}

	n, err = proxy.Write(buf)
	if err != nil {
		log.Fatal(err)
	}
	if n != len(buf) {
		log.Fatal("didn't write the full buffer")
	}

}
