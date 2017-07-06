package main

import (
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"net"
	"sync"
)

var udpServingConns map[int]*net.UDPConn = make(map[int]*net.UDPCon)


func CreateUDPServingPorts() {
	for _, port := range UDP_PORTS {
		udpServingConns[port] = CreateUDPPort(port)
	}
}

var udpTargetAddrs = make([]*net.UDPAddr)

func CreateUDPServingPort(port int) *net.UDPConn {

	udpListenAddr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", PROXY_BIND_IP, port))
	if err != nil {
		log.Fatal(err)
	}
	log.Println("binding to UDP ", udpListenAddr)

	udpTargetAddr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", PARSEC_BIND_IP, port))
	if err != nil {
		panic(err)
	}
	udpTargetAddrs[port] = udpTargetAddr

	udpConn, err = net.ListenUDP("udp", udpListenAddr)
	if err != nil {
		log.Fatal(err)
	}

	for range "1234" {

		// Catch UDP, tunnel into TCP
		go func() {

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

				writeToProxy(port, buffer)
			}
		}()

	}

	return udpConn
}



func startEmiter(udpTargetIP string) {

	/*
	udpTargetAddr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", udpTargetIP, UDP_PORT))
	if err != nil {
		log.Fatal(err)
	}
	*/

	// create 4 of them.
	for range "1234" {

		// Read from TCP and emit packets into UDP or TCP
		go func() {
			buf := make([]byte, 1024*10)

			log.Println("emitter has started")

			for {
				n, port := readFromProxy(buf)

				if n == 0 {
					log.Println("zero length packet from proxy")
					continue
				}

				buf = buf[:n]

				// log.Println("emitting packet, size ", n)

				udpTargetAddr := targetAddrs[port]

				// send UDP
				n, err := udpConns[port].WriteTo(buf, udpTargetAddr)
				if err != nil {
					// log.Println(buf)
					log.Println("packet length ",n)
					log.Println("packet target ",udpTargetAddr)
					panic(err)
				}
				if n != len(buf) {
					panic("udp: wrong length")
				}
			}
		}()

	}
}

