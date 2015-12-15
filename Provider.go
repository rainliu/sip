package sip

import (
	"bufio"
	"log"
	"net"
	"sync"
	"time"
)

////////////////////Interface//////////////////////////////

type Provider interface {
	AddTransport(t Transport)
	RemoveTransport(t Transport)

	AddListener(l Listener)
	RemoveListener(l Listener)
}

////////////////////Implementation////////////////////////

type provider struct {
	listeners    map[Listener]Listener
	transports   map[Transport]Transport
	transactions map[Transaction]Transaction

	//forward   chan Message
	join  chan Transaction
	leave chan Transaction

	quit      chan bool
	waitGroup *sync.WaitGroup

	tracer Tracer
}

func newProvider(tracer Tracer) *provider {
	this := &provider{}

	this.listeners = make(map[Listener]Listener)
	this.transports = make(map[Transport]Transport)
	this.transactions = make(map[Transaction]Transaction)

	//this.forward 	= make(chan Message)
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
			this.tracer.Println("Provider Quit")
			return

		case s := <-this.join:
			this.transactions[s] = s

		case s := <-this.leave:
			delete(this.transactions, s)

			//		case msg := <-this.forward:
			//			for _, s := range this.transactions {
			//				if err := s.Forward(msg); err != nil {
			//					this.tracer.Println(err)
			//					for _, l := range this.listeners {
			//						l.ProcessIOException(newEventIOException(s, s.conn.RemoteAddr()))
			//					}
			//				}
			//			}
		}
	}
}

func (this *provider) Stop() {
	for _, s := range this.transactions {
		s.Close()
	}
	for _, t := range this.transports {
		t.Close()
	}
	close(this.quit)
	this.waitGroup.Wait()
}

func (this *provider) ServeAccept(t *transport) {
	defer this.waitGroup.Done()
	defer t.lner.Close()

	for {
		select {
		case <-t.quit:
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

	s := newServerTransaction()
	this.join <- s

	for {
		select {
		case <-s.quit:
			log.Println("Disconnecting", conn.RemoteAddr())
			//for _, l := range this.listeners {
			//	l.ProcessSessionTerminated(newEventSessionTerminated(s, s.Error(), s.Will()))
			//}
			this.leave <- s
			return
		default:
			//can't delete default, otherwise blocking call
		}

		if req, err := this.readRequest(conn); err != nil {
			if opErr, ok := err.(*net.OpError); ok && opErr.Timeout() {
				/*s.keepAliveAccumulated += 1 //add 1 second
				if s.keepAlive != 0 && s.keepAliveAccumulated >= (s.keepAlive*3)/2 {
					log.Println("Timeout", conn.RemoteAddr())
					for _, l := range this.listeners {
						//l.ProcessTimeout()
					}
					s.Close()
				} else {
					continue
				}*/
			} else {
				log.Println(err)
				//for _, l := range this.listeners {
				//	l.ProcessIOException(newEventIOException(s, conn.RemoteAddr()))
				//}
				s.Close()
			}
		} else {
			log.Println(req.GetMethod())

			/*s.keepAliveAccumulated = 0
			if evt := s.Process(buf); evt != nil {
				switch evt.GetEventType() {
				case EVENT_CONNECT:
					for _, l := range this.listeners {
						l.ProcessConnect(evt.(EventConnect))
					}
				case EVENT_PUBLISH:
					for _, l := range this.listeners {
						l.ProcessPublish(evt.(EventPublish))
					}
				case EVENT_SUBSCRIBE:
					for _, l := range this.listeners {
						l.ProcessSubscribe(evt.(EventSubscribe))
					}
				case EVENT_UNSUBSCRIBE:
					for _, l := range this.listeners {
						l.ProcessUnsubscribe(evt.(EventUnsubscribe))
					}
				case EVENT_IOEXCEPTION:
					for _, l := range this.listeners {
						l.ProcessIOException(evt.(EventIOException))
					}
					s.Terminate(errors.New(s.Error()))
				default:
					s.Terminate(errors.New(s.Error()))
				}
			}*/
		}
	}
}

func (this *provider) readRequest(conn net.Conn) (Request, error) {
	conn.SetDeadline(time.Now().Add(1e9)) //wait for 1 second
	if req, err := ReadRequest(bufio.NewReader(conn)); err != nil {
		return nil, err
	} else {
		return req, nil
	}
}
