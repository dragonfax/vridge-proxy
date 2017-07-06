var udpProxyConn *net.UDPConn

const HANDSHAKE = "hello"

func connectToServerProxy() *net.UDPConn {
	conn, err := net.DialUDP("udp", fmt.Sprintf("%s:%d", SERVER_PUBLIC_IP, 3278))
	if err != nil {
		log.Fatal(err)
	}

	// write inaugrual packet
	// get response.
	buf := []byte{HANDSHAKE}
	writeToProxy(buf)
	readFromProxy(buf)
	if buf != "repsonse" {
		log.Fatal("didn't get expected response from proxy server connection")
	}

	log.Println("client is connected to server")

	return conn.(*net.UDPConn)
}

var udpWriteLock sync.Mutex = sync.Mutex{}

func writeToProxy(buf []byte) {

	udpWriteLock.Lock()
	defer udpWriteLock.Unlock()

	bs := make([]byte, 2)
	binary.LittleEndian.PutUint16(bs, uint16(len(buf)))
	n, err := proxy.Write(bs)
	if err != nil {
		log.Fatal("emitter: ", err)
	}
	if n != 2 {
		log.Fatal("emitter: didn't write the full buffer length header")
	}

	n, err = proxy.Write(buf)
	if err != nil {
		log.Fatal("emitter: ", err)
	}
	if n != len(buf) {
		log.Fatal("emitter: didn't write the full buffer")
	}

}

var udpReadLock sync.Mutex = sync.Mutex{}

func readFromProxy(buf []byte) int {

	udpReadLock.Lock()
	defer udpReadLock.Unlock()

	buf = buf[:2]
	n, err := io.ReadFull(proxyReader, buf)
	if err != nil {
		log.Fatal("emitter: ", err)
	}
	if n != 2 {
		log.Fatal("emitter: wrong length")
	}
	pl := int(binary.LittleEndian.Uint16(buf))

	if pl == 0 {
		panic("zero length packet")
	}

	buf = buf[:pl]
	n, err = io.ReadFull(proxyReader, buf)
	if err != nil {
		log.Fatal("emitter: ", err)
	}
	if n != pl {
		log.Fatal("emitter: wrong length")
	}
	buf = buf[:n]

	return pl
}
