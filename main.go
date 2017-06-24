package main

import (
	"fmt"
	"log"
)

var UDP_PORTS = [...]int{18481} // , 38217}

const TCP_PORTS_LOW = 38216
const TCP_PORTS_HIGH = 38230

var num_tcp_ports = TCP_PORTS_HIGH - TCP_PORTS_LOW + 1
var TCP_PORTS = make([]int, num_tcp_ports)

var serverIP = "64.62.255.87"
var localBindIP = "0.0.0.0"

func initTCPPorts() {
	for i := 0; i < num_tcp_ports; i++ {
		TCP_PORTS[i] = TCP_PORTS_LOW + i
	}
}

func main() {

	initTCPPorts()

	for _, udp_port := range UDP_PORTS {
		CreateUDPPort(localBindIP, udp_port)
	}

	for _, tcp_port := range TCP_PORTS {
		localBindAddr := fmt.Sprintf("%s:%d", localBindIP, tcp_port)
		serverAddr := fmt.Sprintf("%s:%d", serverIP, tcp_port)
		log.Println("binding ", localBindAddr, " to ", serverAddr)
		proxy := NewProxy(localBindAddr, serverAddr)
		err := proxy.Start()
		if err != nil {
			log.Fatal(err)
		}
	}

	for {
	}
}
