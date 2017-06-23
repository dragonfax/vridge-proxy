package main

import (
	"fmt"
	"log"
	"net"
)

const MAX_UDP_PACKET_LENGTH = 65535

func CreateUDPPort(bindIP string, port int) {

	bindAddr := fmt.Sprintf("%s:%d", bindIP, port)
	log.Println("binding to UDP ", bindAddr)

	pc, err := net.ListenPacket("udp", bindAddr)
	if err != nil {
		log.Fatal(err)
	}

	var clientAddr *net.UDPAddr
	foundClientAddr := false
	serverAddr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", serverIP, port))
	if err != nil {
		log.Fatal(err)
	}

	buffer := make([]byte, MAX_UDP_PACKET_LENGTH)

	// deadline := time.Now().Add(time.Hour * 999)

	// pc.SetDeadline(deadline)

	go func() {

		for {
			buffer = buffer[:cap(buffer)]
			n, sAddr, err := pc.ReadFrom(buffer)
			if err != nil {
				log.Fatal(err)
			}
			buffer = buffer[:n]

			if sAddr == nil && len(buffer) == 0 {
				// log.Println("empty packet")
				continue
			}

			sourceAddr := sAddr.(*net.UDPAddr)

			if sourceAddr.String() == serverAddr.String() {
				log.Println("received server packet")
				pc.WriteTo(buffer, clientAddr)
			} else {

				// first time save the address
				// after that, verify the address is expected.
				if !foundClientAddr {
					clientAddr = sourceAddr
					foundClientAddr = true
					log.Println("found client of ", clientAddr)
				}

				if sourceAddr.String() != clientAddr.String() {
					log.Fatal("unknown source address")
				}

				// log.Println("received client packet")

				n2, err := pc.WriteTo(buffer, serverAddr)
				if err != nil {
					log.Fatal(err)
				}

				if n != n2 {
					log.Fatal("failed to write all of the packet")
				}
			}

		}
	}()
}
