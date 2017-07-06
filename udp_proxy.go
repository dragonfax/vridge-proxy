var udpProxyConn *net.UDPConn

const HANDSHAKE = "hello"

func connectToServerProxy() {
	udpProxyConn, err := net.DialUDP("udp", fmt.Sprintf("%s:%d", SERVER_PUBLIC_IP, PROXY_UDP_PORT))
	if err != nil {
		log.Fatal(err)
	}

	// write inaugrual packet
	// get response.
	buf := []byte{HANDSHAKE}
	writeToProxy(buf)
	n := readFromProxy(buf)
	if len(buf) != len(HANDSHAKE) {
		log.Fatal("wrong length of handshake responnse")
	}
	if buf != HANDSHAKE {
		log.Fatal("didn't get expected response from proxy server connection")
	}

	log.Println("client is connected to server, handshake complete")
}

func listenAsServerProxy() {
	udpProxyAddr, err := net.ResolveUDPAddr("udp",fmt.Sprintf("%s:%d","0.0.0.0", PROXY_UDP_PORT)
	if err != nil {
		panic(err)
	}

	udpProxyConn, err := net.ListenUDP("udp",updProxyAddr)
	if err != nil {
		panic(err)
	}

	buf := make([]byte,len(HANDSHAKE),len(HANDSHAKE))
	n := readFromProxy(buf)
	if n != len(HANDSHAKE) {
		panic("wrong hand shake length")
	}
	if buf != HANDSHAKE {
		panic("wrong hand shake")
	}

	log.Println("server is connected to client")
}

var udpWriteLock sync.Mutex = sync.Mutex{}

func writeToProxy(port int, payload []byte) {

	udpWriteLock.Lock()
	defer udpWriteLock.Unlock()

	// port
	buf := make([]byte, 2 + len(payload))
	binary.LittleEndian.PutUint16(buf, uint16(port))

	copy(buf[2:],payload)

	// body
	n, err = udpProxyConn.Write(buf)
	if err != nil {
		log.Fatal("emitter: ", err)
	}
	if n != len(buf) {
		log.Fatal("emitter: didn't write the full buffer")
	}

}

var udpReadLock sync.Mutex = sync.Mutex{}

func readFromProxy(payload []byte) (n int, port int) {

	udpReadLock.Lock()
	defer udpReadLock.Unlock()

	buf := make([]byte,len(payload) + 2)
	n := udpProxyConn.Read(buf)

	if n == 1 {
		panic("port not read from proxy")
	}

	if n == 0 {
		return 0, 0
	}

	buf = buf[:2]
	port := int(binary.LittleEndian.Uint16(buf))

	buf = buf[2:n-2]
	copy(payload,buf)
	return n-2, port
}
