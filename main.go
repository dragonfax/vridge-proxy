package main

import "flag"

// these are still terrible names

// public IP of the VM
// for the proxy to connect to
const SERVER_PUBLIC_IP = "64.62.255.87"

const CLIENT_PUBLIC_IP = ""

// Same on both sides.
const PROXY_BIND_IP = "192.168.10.34"
const PARSEC_BIND_IP = "192.168.10.35"

// Port the proxy uses its private comms
const PROXY_UDP_PORT = 3278

var UDP_PORTS = []int{8000, 8001, 8002, 8003, 8004, 8005, 8006}

const TCP_PORTS_LOW = 38216
const TCP_PORTS_HIGH = 38230

var num_tcp_ports = TCP_PORTS_HIGH - TCP_PORTS_LOW + 1
var TCP_PORTS = make([]int, num_tcp_ports)

func main() {

	var serverFlagP = flag.Bool("server", false, "run the server side of the proxy")
	flag.Parse()

	if *serverFlagP {
		server()
	} else {
		client()
	}

	for {
		// wait
	}
}
