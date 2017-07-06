func server() {

	initServerTCPPorts()

	proxy = listenAsServerProxy()
	proxyReader = bufio.NewReader(proxy)

	startEmiter(SERVER_VRIDGE_IP)

	CreateUDPPort(SERVER_PROXY_SOURCE_IP, UDP_PORT, proxy)
}
