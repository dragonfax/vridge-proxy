package main

func client() {

	initClientTCPPorts()

	connectToServerProxy()

	startEmiter(CLIENT_IP)

	CreateUDPPorts(CLIENT_VRIDGE_LISTEN_IP, UDP_PORT, proxy)
}
