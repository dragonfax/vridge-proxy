package main

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"log"
	"net"
	"sync"
)



var udpConn *net.UDPConn

func CreateUDPPort(udpListenIP string, port int, tcpProxy net.Conn) {

	udpListenAddr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", udpListenIP, port))
	if err != nil {
		log.Fatal(err)
	}
	log.Println("binding to UDP ", udpListenAddr)

	udpConn, err = net.ListenUDP("udp", udpListenAddr)
	if err != nil {
		log.Fatal(err)
	}

	for range "1234" {

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
}

var udpWriteLock sync.Mutex = sync.Mutex{}

func writeToProxy(buf []byte) {

	udpWriteLock.Lock()
	defer udpWriteLock.Unlock()

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

	// create 4 of them.
	for range "1234" {

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
}

var udpReadLock sync.Mutex = sync.Mutex{}

func readFromProxy(buf []byte) int {

	udpReadLock.Lock()
	defer udpReadLock.Unlock()

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
