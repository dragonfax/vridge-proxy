package main

func server() {

	initServerTCPPorts()

	listenAsServerProxy()

	startEmiter(SERVER_VRIDGE_IP)

	CreateUDPPort(SERVER_PROXY_SOURCE_IP, UDP_PORT, proxy)
}
