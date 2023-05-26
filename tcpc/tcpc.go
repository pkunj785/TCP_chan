package tcpc

import (
	"encoding/gob"
	"log"
	"net"
	"sync"
	"time"
)

type TCPC[T any] struct {
	listenAddr string
	remoteAddr string

	Sendch       chan T
	Recvch       chan T
	outboundConn net.Conn
	ln           net.Listener
	wg           sync.WaitGroup
}

func New[T any](listenAddr, remoteAddr string) (*TCPC[T], error) {
	tcpc := &TCPC[T]{
		listenAddr: listenAddr,
		remoteAddr: remoteAddr,
		Sendch:     make(chan T, 10),
		Recvch:     make(chan T, 10),
	}

	ln, err := net.Listen("tcp", listenAddr)

	if err != nil {
		return nil, err
	}
	tcpc.ln = ln

	go tcpc.loop()
	go tcpc.acceptLoop()
	go tcpc.dialRemoteAndRead()

	return tcpc, nil
}


func (t *TCPC[T]) loop() {
	t.wg.Wait()
	
	for {
		msg := <-t.Sendch
		log.Println("sending msg over wire:", msg)
		if err := gob.NewEncoder(t.outboundConn).Encode(&msg); err != nil {
			log.Println(err)
			continue
		}
		log.Println("Done")
	}
}

func (t *TCPC[T]) acceptLoop() {
	defer func() {
		t.ln.Close()
	}()

	for {
		conn, err := t.ln.Accept()

		if err != nil {
			log.Println("accept error:", err)
			return
		}

		log.Printf("sender connected %s", conn.RemoteAddr())

		go t.handleConn(conn)
	}
} ////// handle connection for accept() and close() endpoint

func (t *TCPC[T]) handleConn(conn net.Conn) {
	defer func() {
		conn.Close()
	}()

	for {
		var msg T
		if err := gob.NewDecoder(conn).Decode(&msg); err != nil {
			log.Println(err)
			return
		}
		t.Recvch <- msg
	}
}

func (t *TCPC[T]) dialRemoteAndRead() {
	t.wg.Add(1)
	
	conn, err := net.Dial("tcp", t.remoteAddr)

	if err != nil {
		log.Printf("dial error (%s): ", err)
		time.Sleep(3 * time.Second)
		t.dialRemoteAndRead()
	}

	t.outboundConn = conn

	t.wg.Done()

}