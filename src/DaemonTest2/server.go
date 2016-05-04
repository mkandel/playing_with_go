package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	// Non-blocking channel library
	"github.com/hectane/go-nonblockingchan"
)

// stat: type for gathering statistics
type stat struct {
	name string
	val  int
}

// AppServer:
type AppServer struct {
	isRunning      bool
	processedCount int
	stats          []stat
}

// Version(): returns a string representation of the application version
func (AppServer) Version() string {
	return version
}

// PrintVersion(): prints a string representation of the application version to STDOUT
func (srv *AppServer) PrintVersion() {
	fmt.Println(srv.Version())
}

// GetStats(): Returns an array of 'stat's
func (srv AppServer) GetStats() *[]stat {
	return &srv.stats
}

// Start(port int): Starts the listener on 'port'
func (srv *AppServer) Start(port int, wgMain sync.WaitGroup) bool {
	listen_port := ":" + strconv.Itoa(port)
	l, err := net.Listen("tcp4", listen_port)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Server (Version '%s') listening on port '%s'\n", srv.Version(), listen_port)
	defer func() {
		log.Println("Closing Server listen port")
		l.Close()
	}()

	srv.isRunning = true

	cmdChan := nbc.New()

	var wgConn sync.WaitGroup

	for {
		conn, err := l.Accept()
		defer conn.Close()
		srv.processedCount++
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
			log.Println("Server: Received from channel: ", cmd)
			wgMain.Done()
			log.Println("Server: Returning ...")
			return true
		case <-time.After(1 * time.Nanosecond):
			log.Println("Timeout")
		}
	}

	return true
}

func handleConn(c net.Conn, cmdChan *nbc.NonBlockingChan, wg sync.WaitGroup) {
	buf, err := ioutil.ReadAll(bufio.NewReader(c))
	if err != nil {
		log.Println("ioutil.ReadAll() error trapped:")
		log.Println(err)
		return
	}
	message := strings.TrimSpace(string(buf))
	shutdownMsg := "^shutdown_server$"
	log.Print("Message: ", message)
	log.Printf("Message size(bytes): %d\n", len(message))

	matched, err := regexp.MatchString(shutdownMsg, message)
	if err != nil {
		log.Println("Error received: ", err)
		return
	}
	if matched {
		log.Println("Sending shutdown message")
		cmdChan.Send <- message
		wg.Done()
	}

	// Shut down the connection.
	log.Println("Closing connection: ", c.RemoteAddr())
	log.Print(sep)
	c.Close()
	return
}

func (srv *AppServer) ProcessedCount() int {
	return srv.processedCount
}
