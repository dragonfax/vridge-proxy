package main

import "flag"

const serverIP = "64.62.255.87"
const localBindIP = "0.0.0.0"
const clientIP = "192.168.0.100"

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
