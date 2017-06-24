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
			n, port, id := readFromProxy(buf)
			buf = buf[:n]
			key := fmt.Sprintf("%d:%d", port, id)

			if port == UDP_PORT {
				// send UDP
				n, err := udpConn.WriteTo(buf, udpAddr)
				if err != nil {
					log.Fatal(err)
				}
				if n != len(buf) {
					log.Fatal("wrong length")
				}
			} else {
				// send TCP

				if len(buf) == 0 {
					// A zero length packet from TCP means to initiate the connection, or close it

					log.Println("received zero length packet for port ", port, " id ", id)

					_, ok := connections[key]
					if ok {
						log.Println("exists, closing connection")
						connections[key].Close()
						delete(connections, key)
					} else {
						log.Println("doesn't exist, creating")

						serverAddr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:%d", serverIP, port))
						if err != nil {
							panic(err)
						}

						conn, err := net.DialTCP("tcp", nil, serverAddr)
						if err != nil {
							panic(err)
						}

						// track the connection
						connections[key] = conn

						log.Println("created server side connection to port ", port, " id ", id)

						go handle(conn, port, id)
					}
				} else {

					connection, ok := connections[key]
					if !ok {
						log.Fatal("no connection for port ", port, " id ", id)
					}
					n, err := connection.Write(buf)
					if err != nil {
						log.Fatal(err)
					}
					if n != len(buf) {
						log.Fatal("failed to write all of packet")
					}
				}
			}
		}
	}()
}

func readFromProxy(buf []byte) (int, int, int) {

	buf = buf[:2]
	n, err := io.ReadFull(proxyReader, buf)
	if err != nil {
		log.Fatal(err)
	}
	if n != 2 {
		log.Fatal("wrong length")
	}
	pl := int(binary.LittleEndian.Uint16(buf))

	n, err = io.ReadFull(proxyReader, buf)
	if err != nil {
		log.Fatal(err)
	}
	if n != 2 {
		log.Fatal("wrong length")
	}
	port := int(binary.LittleEndian.Uint16(buf))

	n, err = io.ReadFull(proxyReader, buf)
	if err != nil {
		log.Fatal(err)
	}
	if n != 2 {
		log.Fatal("wrong length")
	}
	id := int(binary.LittleEndian.Uint16(buf))

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

	return pl, port, id
}
