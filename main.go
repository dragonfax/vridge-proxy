package main

import "flag"

// these are still terrible names

// public IP of the VM
// for the proxy to connect to
const SERVER_PUBLIC_IP = "64.62.255.87"

const SERVER_VRIDGE_IP = "192.168.10.34"

// IP the proxy sends traffic to Vridge from
const SERVER_PROXY_SOURCE_IP = "192.168.10.35"

// IP the proxy listens to for the proxy connection
const SERVER_PROXY_LISTEN_IP = "0.0.0.0"

// on the macbook, the proxy binds to this. to listen for vridge traffic
const CLIENT_VRIDGE_LISTEN_IP = "0.0.0.0"

// TODO not really necessary, I think.
const TCP_PROXY_PORT = 23432

const PROXY_PORT = 3278

// the IP of thep phone.
const CLIENT_IP = "192.168.0.100"

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
