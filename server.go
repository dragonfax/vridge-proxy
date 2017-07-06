package main

func server() {

	initServerTCPPorts()

	listenAsServerProxy()

	startEmiter()

	CreateUDPServingPorts()
}
