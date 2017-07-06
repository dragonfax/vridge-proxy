package main

import (
	"encoding/binary"
	"fmt"
	"log"
	"net"
	"reflect"
	"sync"
)

var udpProxyConn *net.UDPConn

var handshake = []byte{'h', 'e', 'l', 'l', '0'}

func connectToServerProxy() {
	udpProxyAddr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", serverIP, PROXY_UDP_PORT))
	if err != nil {
		panic(err)
	}

	udpProxyConn, err = net.DialUDP("udp", nil, udpProxyAddr)
	if err != nil {
		log.Fatal(err)
	}

	// write inaugrual packet
	// get response.
	writeToProxy(0, handshake)

	buf := make([]byte, len(handshake))
	readFromProxy(buf)
	if !reflect.DeepEqual(buf, handshake) {
		log.Fatal("didn't get expected response from proxy server connection")
	}

	log.Println("client is connected to server, handshake complete")
}

func listenAsServerProxy() {
	udpProxyAddr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", "0.0.0.0", PROXY_UDP_PORT))
	if err != nil {
		panic(err)
	}

	udpProxyConn, err = net.ListenUDP("udp", udpProxyAddr)
	if err != nil {
		panic(err)
	}

	buf := make([]byte, len(handshake))
	readFromProxy(buf)
	if !reflect.DeepEqual(buf, handshake) {
		panic("wrong hand shake")
	}

	writeToProxy(0, handshake)

	log.Println("server is connected to client")
}

var udpWriteLock sync.Mutex = sync.Mutex{}

func writeToProxy(port int, payload []byte) {

	udpWriteLock.Lock()
	defer udpWriteLock.Unlock()

	// port
	buf := make([]byte, 2+len(payload))
	binary.LittleEndian.PutUint16(buf, uint16(port))

	copy(buf[2:], payload)

	// body
	n, err := udpProxyConn.Write(buf)
	if err != nil {
		log.Fatal("emitter: ", err)
	}
	if n != len(buf) {
		log.Fatal("emitter: didn't write the full buffer")
	}

}

var udpReadLock sync.Mutex = sync.Mutex{}

func readFromProxy(payload []byte) (pl int, port int) {

	udpReadLock.Lock()
	defer udpReadLock.Unlock()

	buf := make([]byte, len(payload)+2)
	n, err := udpProxyConn.Read(buf)
	if err != nil {
		panic(err)
	}

	if n == 1 {
		panic("port not read from proxy")
	}

	if n == 0 {
		return 0, 0
	}

	buf = buf[:2]
	port = int(binary.LittleEndian.Uint16(buf))

	buf = buf[2 : n-2]
	copy(payload, buf)
	return n - 2, port
}
