package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"sync"
)

const TCP_PORTS_LOW = 38216
const TCP_PORTS_HIGH = 38230

var num_tcp_ports = TCP_PORTS_HIGH - TCP_PORTS_LOW + 1
var TCP_PORTS = make([]int, num_tcp_ports)

func initClientTCPPorts() {
	initTCPPorts(0, -1000)
}

func initServerTCPPorts() {
	initTCPPorts(-1000, 0)
}

func initTCPPorts(port_adjust_local int, port_adjust_remote int) {
	for i := 0; i < num_tcp_ports; i++ {
		TCP_PORTS[i] = TCP_PORTS_LOW + i
	}

	for _, tcp_port := range TCP_PORTS {
		localBindAddr := fmt.Sprintf("%s:%d", localBindIP, tcp_port+port_adjust_local)
		log.Println("binding ", localBindAddr)
		remoteBindAddr := fmt.Sprintf("%s:%d", serverIP, tcp_port+port_adjust_remote)
		p := NewProxy(localBindAddr, remoteBindAddr)
		err := p.Start()
		if err != nil {
			log.Fatal(err)
		}
	}

	log.Println("tcp ports are listening")
}

type Proxy struct {
	from string
	to   string
	done chan struct{}
}

func NewProxy(from, to string) *Proxy {
	return &Proxy{
		from: from,
		to:   to,
		done: make(chan struct{}),
	}
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
				log.Print("tcp: ", err)
			}
		}
	}
}

func (p *Proxy) handle(connection net.Conn) {
	log.Println("Handling", connection)
	defer log.Println("Done handling", connection)
	defer connection.Close()
	remote, err := net.Dial("tcp", p.to)
	if err != nil {
		log.Print("tcp: ", err)
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
			log.Print("tcp: ", err)
			p.Stop()
			return
		}
	}
}
