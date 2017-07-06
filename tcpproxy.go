package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"sync"
)

func initClientTCPPorts() {
	initTCPPorts("0.0.0.0", "", SERVER_PUBLIC_IP, 0, -1000)
}

func initServerTCPPorts() {
	initTCPPorts(SERVER_PROXY_LISTEN_IP, SERVER_PROXY_SOURCE_IP, SERVER_VRIDGE_IP, -1000, 0)
}

func initTCPPorts(from, fromto, to string, port_adjust_local int, port_adjust_remote int) {
	for i := 0; i < num_tcp_ports; i++ {
		TCP_PORTS[i] = TCP_PORTS_LOW + i
	}

	for _, tcp_port := range TCP_PORTS {
		localBindAddr := fmt.Sprintf("%s:%d", from, tcp_port+port_adjust_local)
		fromto := fmt.Sprintf("%s:0", fromto)
		remoteBindAddr := fmt.Sprintf("%s:%d", to, tcp_port+port_adjust_remote)

		log.Println("binding local ", localBindAddr, " through ", fromto, " to ", remoteBindAddr)
		p := NewProxy(localBindAddr, fromto, remoteBindAddr)
		err := p.Start()
		if err != nil {
			log.Fatal(err)
		}
	}

	log.Println("tcp ports are listening")
}

type Proxy struct {
	from   string
	fromto string
	to     string
	done   chan struct{}
}

func NewProxy(from, fromto, to string) *Proxy {
	p := &Proxy{
		from:   from,
		fromto: fromto,
		to:     to,
		done:   make(chan struct{}),
	}
	return p
}

func (p *Proxy) Start() error {
	listener, err := net.Listen("tcp", p.from)
	if err != nil {
		return err
	}
	go p.run(listener)
	return nil
}

func (p *Proxy) Stop() {
	if p.done == nil {
		return
	}
	close(p.done)
	p.done = nil
}

func (p *Proxy) run(listener net.Listener) {
	for {
		select {
		case <-p.done:
			return
		default:
			connection, err := listener.Accept()
			if err == nil {
				go p.handle(connection)
			} else {
				log.Println("tcp: ", err)
			}
		}
	}
}

func (p *Proxy) handle(connection net.Conn) {
	log.Println("Handling", connection)
	defer log.Println("Done handling", connection)
	defer connection.Close()

	to, err := net.ResolveTCPAddr("tcp", p.to)
	if err != nil {
		log.Fatal(err)
	}

	var fromto *net.TCPAddr
	if p.fromto != "" {
		fromto, err = net.ResolveTCPAddr("tcp", p.fromto)
		if err != nil {
			log.Fatal(err)
		}
	}

	remote, err := net.DialTCP("tcp", fromto, to)
	if err != nil {
		log.Println("tcp: ", err)
		return
	}
	defer remote.Close()
	wg := &sync.WaitGroup{}
	wg.Add(2)
	go p.copy(remote, connection, wg)
	go p.copy(connection, remote, wg)
	wg.Wait()
}

func (p *Proxy) copy(from, to net.Conn, wg *sync.WaitGroup) {
	defer wg.Done()
	select {
	case <-p.done:
		return
	default:
		if _, err := io.Copy(to, from); err != nil {
			log.Println("tcp: ", err)
			p.Stop()
			return
		}
	}
}
