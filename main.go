package main

import (
	"flag"
	"fmt"
	"net"
	"strconv"
)

// these are still terrible names

// Same on both sides.
const PROXY_BIND_IP = "192.168.10.34"
const PARSEC_BIND_IP = "192.168.10.35"

// Port the proxy uses its private comms
const PROXY_UDP_PORT = 3278

type portSlice []int

func (i *portSlice) String() string {
	return fmt.Sprintf("%d", *i)
}

func (i *portSlice) Set(value string) error {
	tmp, err := strconv.Atoi(value)
	if err != nil {
		return err
	}

	*i = append(*i, tmp)
	return nil
}

var udpPorts portSlice

var tcpPorts portSlice

// var clientIP string

var serverIP string

func main() {

	flag.Var(&udpPorts, "udp", "List of udp ports")
	flag.Var(&udpPorts, "tcp", "List of tcp ports")
	// flag.StringVar(&clientIP, "client", "Public Client IP")
	flag.StringVar(&serverIP, "server", "", "Public Server IP")

	flag.Parse()

	if flag.NFlag() == 0 {
		flag.PrintDefaults()
	}

	udpTargetAddrs = make([]*net.UDPAddr, len(udpPorts))

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
