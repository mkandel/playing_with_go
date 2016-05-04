package main

import (
	"fmt"
	"log"
	"net"
	"strconv"
	"sync"
	"time"

	"github.com/hectane/go-nonblockingchan"
)

type serverStat struct {
	name string
	val  int
}

type APIServer struct {
	isRunning      bool
	processedCount int
	stats          []serverStat
}

func (APIServer) Version() string {
	return apiVers
}

func (APIServer) PrintVersion() {
	fmt.Println(apiVers)
}

func (srv *APIServer) GetStats() *[]serverStat {
	return &srv.stats
}

// Start(port int): Starts the listener on 'port'
func (srv *APIServer) Start(port int, wgMain sync.WaitGroup) bool {
	listen_port := ":" + strconv.Itoa(port)
	l, err := net.Listen("tcp4", listen_port)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("APIServer (API Version '%s') listening on port '%s'\n", srv.Version(), listen_port)
	defer func() {
		log.Println("Closing APIServer listen port")
		l.Close()
	}()

	srv.isRunning = true

	cmdChan := nbc.New()

	var wgConn sync.WaitGroup

	for {
		conn, err := l.Accept()
		defer conn.Close()
		if err != nil {
			log.Fatal(err)
		}
		// Handle the connection in a new goroutine.
		// The loop then returns to accepting, so that
		// multiple connections may be served concurrently.
		wgConn.Add(1)
		go handleConn(conn, cmdChan, wgConn)

		tm := time.Now()
		log.Println(tm.Format(time.RFC3339))

		select {
		case cmd := <-cmdChan.Recv:
			log.Println("APIServer: Received from channel: ", cmd)
			log.Println("APIServer: Calling WaitGroup.Done()")
			wgMain.Done()
			return true
		case <-time.After(1 * time.Nanosecond):
			log.Println("Timeout")
		}
		// Do this here so we don't count 'shutdown_server' messages
		srv.processedCount++
	}

	return true
}
