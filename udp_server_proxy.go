package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
)


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