package sip

import (
	"bufio"
	"bytes"
	"log"
	"net"
	"sync"
	"time"
)

////////////////////Interface//////////////////////////////

type Provider interface {
	AddTransport(Transport)
	RemoveTransport(Transport)

	AddListener(Listener)
	RemoveListener(Listener)

	GetNewCallId() string

	GetNewClientTransaction(Request) ClientTransaction
	GetNewServerTransaction(Request) ServerTransaction

	SendRequest(Request) error
	SendResponse(Response) error
}

////////////////////Implementation////////////////////////

type provider struct {
	listeners    map[Listener]Listener
	transports   map[Transport]Transport
	transactions map[Transaction]Transaction

	forward chan Message
	join    chan Transaction
	leave   chan Transaction

	quit      chan bool
	waitGroup *sync.WaitGroup

	tracer Tracer
}

func newProvider(tracer Tracer) *provider {
	this := &provider{}

	this.listeners = make(map[Listener]Listener)
	this.transports = make(map[Transport]Transport)
	this.transactions = make(map[Transaction]Transaction)

	this.forward = make(chan Message)
	this.join = make(chan Transaction)
	this.leave = make(chan Transaction)

	this.quit = make(chan bool)
	this.waitGroup = &sync.WaitGroup{}

	this.tracer = tracer

	return this
}

func (this *provider) AddTransport(t Transport) {
	this.transports[t] = t
}

func (this *provider) RemoveTransport(t Transport) {
	delete(this.transports, t)
}

func (this *provider) AddListener(l Listener) {
	this.listeners[l] = l
}

func (this *provider) RemoveListener(l Listener) {
	delete(this.listeners, l)
}

func (this *provider) GetNewCallId() string {
	return ""
}

func (this *provider) GetNewClientTransaction(req Request) ClientTransaction {
	ct := newClientTransaction(req)
	this.join <- ct
	return ct
}
func (this *provider) GetNewServerTransaction(req Request) ServerTransaction {
	st := newServerTransaction(req)
	this.join <- st
	return st
}

func (this *provider) SendRequest(Request) error {
	return nil
}
func (this *provider) SendResponse(Response) error {
	return nil
}

func (this *provider) Run() {
	for _, t := range this.transports {
		if err := t.Listen(); err != nil {
			this.tracer.Printf("Listening %s://%s:%d Failed!!!\n", t.GetNetwork(), t.GetAddress(), t.GetPort())
		} else {
			this.tracer.Printf("Listening %s://%s:%d Runing...\n", t.GetNetwork(), t.GetAddress(), t.GetPort())
			this.waitGroup.Add(1)
			go this.ServeAccept(t.(*transport))
		}
	}

	//infinite loop run until ctrl+c
	for {
		select {
		case <-this.quit:
			this.tracer.Println("Provider Stopped!!!")
			return

		case s := <-this.join:
			this.transactions[s] = s

		case s := <-this.leave:
			delete(this.transactions, s)

		case msg := <-this.forward:
			var buffer bytes.Buffer
			if err := msg.StartLineWrite(&buffer); err != nil {
				log.Println(err)
			} else {
				log.Println("Received: ", buffer.String())
			}
		}
	}
}

func (this *provider) Stop() {
	close(this.quit)
	for _, s := range this.transactions {
		s.Close()
	}
	this.waitGroup.Wait()
}

func (this *provider) ServeAccept(t *transport) {
	defer this.waitGroup.Done()
	defer t.lner.Close()

	for {
		select {
		case <-this.quit:
			log.Printf("Listening %s://%s:%d Stoped!!!\n", t.GetNetwork(), t.GetAddress(), t.GetPort())
			return
		default:
			//can't delete default, otherwise blocking call
		}
		t.SetDeadline(time.Now().Add(1e9))
		conn, err := t.Accept()
		if err != nil {
			if opErr, ok := err.(*net.OpError); !(ok && opErr.Timeout()) {
				log.Println(err)
			}
			continue
		}
		this.waitGroup.Add(1)
		go this.ServeConn(conn)
	}
}

func (this *provider) ServeConn(conn net.Conn) {
	defer this.waitGroup.Done()
	defer conn.Close()

	for {
		select {
		case <-this.quit:
			log.Println("Disconnecting...", conn.RemoteAddr())
			return
		default:
			//can't delete default, otherwise blocking call
		}

		conn.SetDeadline(time.Now().Add(1e9)) //wait for 1 second
		if msg, err := ReadMessage(bufio.NewReader(conn)); err != nil {
			if opErr, ok := err.(*net.OpError); ok && opErr.Timeout() {
				continue
			} else {
				log.Println(err)
				return
			}
		} else {
			this.forward <- msg
		}
	}
}
