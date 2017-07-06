package main

func client() {

	initClientTCPPorts()

	connectToServerProxy()

	startEmiter()

	CreateUDPServingPorts()
}
