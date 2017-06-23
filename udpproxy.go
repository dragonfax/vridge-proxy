package main

import (
	"fmt"
	"log"
	"net"
)

const MAX_UDP_PACKET_LENGTH = 65535

func CreateUDPPort(bindIP string, port int) {

	bindAddr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", bindIP, port))
	if err != nil {
		log.Fatal(err)
	}
	log.Println("binding to UDP ", bindAddr)

	pc, err := net.ListenUDP("udp", bindAddr)
	if err != nil {
		log.Fatal(err)
	}

	serverAddr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", serverIP, port))
	if err != nil {
		log.Fatal(err)
	}

	// lets start 4 of these up.
	for _, _ = range "1234" {

		go func() {

			// needs to be global, not goroutine local
			var clientAddr *net.UDPAddr
			foundClientAddr := false

			buffer := make([]byte, MAX_UDP_PACKET_LENGTH)
			oobuf := make([]byte, 1024)

			for {
				buffer = buffer[:cap(buffer)]
				oobuf = oobuf[:cap(oobuf)]
				n, oobn, _, sourceAddr, err := pc.ReadMsgUDP(buffer, oobuf)
				if err != nil {
					log.Fatal(err)
				}
				buffer = buffer[:n]
				oobuf = oobuf[:oobn]

				if sourceAddr == nil && len(buffer) == 0 {
					// log.Println("empty packet")
					continue
				}

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

					n2, _, err := pc.WriteMsgUDP(buffer, oobuf, serverAddr)
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
}
