package main

import (
	"fmt"
	"io"
	"log"
	"net"
)

const TCP_PORTS_LOW = 38216
const TCP_PORTS_HIGH = 38230

var num_tcp_ports = TCP_PORTS_HIGH - TCP_PORTS_LOW + 1
var TCP_PORTS = make([]int, num_tcp_ports)

func initTCPPorts() {
	for i := 0; i < num_tcp_ports; i++ {
		TCP_PORTS[i] = TCP_PORTS_LOW + i
	}

	for _, tcp_port := range TCP_PORTS {
		localBindAddr := fmt.Sprintf("%s:%d", localBindIP, tcp_port)
		log.Println("binding ", localBindAddr)
		err := Start(localBindAddr)
		if err != nil {
			log.Fatal(err)
		}
	}

	log.Println("tcp ports are listening")
}

func Start(from string) error {
	addr, err := net.ResolveTCPAddr("tcp", from)
	if err != nil {
		log.Fatal(err)
	}

	listener, err := net.Listen("tcp", from)
	if err != nil {
		return err
	}
	go run(listener, addr.Port)
	return nil
}

// TODO not a synchronous map
var connections map[string]*net.TCPConn = make(map[string]*net.TCPConn)

var id_seq int = 0

func run(listener net.Listener, port int) {
	for {
		connection, err := listener.Accept()
		if err != nil {
			log.Print(err)
		} else {

			// track the connections
			id_seq += 1
			id := id_seq
			key := fmt.Sprintf("%d:%d", port, id)
			_, ok := connections[key]
			if ok {
				log.Fatal("creating duplicate port id connection")
			}
			connections[key] = connection.(*net.TCPConn)

			// announce the connection to the server
			writeToProxy([]byte{}, port, id)
			log.Println("announced new tcp connection for port ", port, " id ", id)

			go handle(connection, port, id)
		}
	}
}

func handle(connection net.Conn, port int, id int) {
	log.Println("Handling tcp port ", port, " id ", id)
	defer log.Println("Done handling tcp port ", port, " id ", id)
	defer connection.Close()

	buf := make([]byte, 1024*32)

	for {
		n, err := connection.Read(buf)
		if err != nil {
			if err == io.EOF {
				log.Println("connection closed, port ", port, " id ", id)
				break
			} else {
				log.Println("ignoring: ", err)
				return
			}
		}

		if n == 0 {
			log.Println("got a 0 length read from TCP")
			continue
		}
		writeToProxy(buf[0:n], port, id)
	}

	writeToProxy([]byte{}, port, id)
	log.Println("announced closing tcp connection for port ", port, " id ", id)
}
