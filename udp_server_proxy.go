package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
)

var UDP_PORT = 18481

func listenAsServerProxy() *net.TCPConn {
	ln, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", 3278))
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

	startEmiter(serverIP)

	CreateUDPPort(localBindIP, 23432, proxy)
}
