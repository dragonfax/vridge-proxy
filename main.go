package main

import (
	"flag"
	"fmt"
	"net"
	"os"
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

var serverFlagP *bool

func main() {

	flag.Var(&udpPorts, "udp", "List of udp ports")
	flag.Var(&tcpPorts, "tcp", "List of tcp ports")
	// flag.StringVar(&clientIP, "client", "Public Client IP")
	flag.StringVar(&serverIP, "serverIP", "", "Public Server IP")

	serverFlagP = flag.Bool("server", false, "run the server side of the proxy")
	flag.Parse()

	if flag.NFlag() == 0 {
		flag.PrintDefaults()
		os.Exit(0)
	}

	udpTargetAddrs = make(map[int]*net.UDPAddr)

	if *serverFlagP {
		server()
	} else {
		client()
	}

	for {
		// wait
	}
}
