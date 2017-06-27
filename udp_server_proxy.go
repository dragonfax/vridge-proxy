package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
)

const UDP_PORT = 18481

func listenAsServerProxy() *net.TCPConn {
	ln, err := net.Listen("tcp", fmt.Sprintf("%s:%d", SERVER_PROXY_LISTEN_IP, PROXY_PORT))
	if err != nil {
		log.Fatal(err)
	}

	conn, err := ln.Accept()
	if err != nil {
		log.Fatal(err)
	}

	log.Println("server is connected to client")

	return conn.(*net.TCPConn)
}

var proxy *net.TCPConn
var proxyReader *bufio.Reader

func server() {

	initServerTCPPorts()

	proxy = listenAsServerProxy()
	proxyReader = bufio.NewReader(proxy)

	startEmiter(SERVER_VRIDGE_IP)

	CreateUDPPort(SERVER_PROXY_SOURCE_IP, UDP_PORT, proxy)
}
