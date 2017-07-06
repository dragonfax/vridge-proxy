func client() {
	initClientTCPPorts()

	proxy = connectToServerProxy()
	proxyReader = bufio.NewReader(proxy)
	startEmiter(CLIENT_IP)

	CreateUDPPort(CLIENT_VRIDGE_LISTEN_IP, UDP_PORT, proxy)
}