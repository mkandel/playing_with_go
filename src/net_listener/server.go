package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"strconv"
)

//const listen_port string = ":7777"
const sep string = "=====================================================================\n"

func main() {
	var port = flag.Int("p", 7777, "Listen port (default: 7777)")
	flag.Parse()
	listen_port := ":" + strconv.Itoa(*port)
	fmt.Printf("Port: %d\tStr: %s\n", *port, listen_port)
	l, err := net.Listen("tcp4", listen_port)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Print(sep)
	fmt.Printf("server waiting for client connection on port '%s'\n",
		listen_port)
	fmt.Print(sep)
	defer l.Close()
	for {
		conn, err := l.Accept()
		defer conn.Close()
		if err != nil {
			log.Fatal(err)
		}
		// Handle the connection in a new goroutine.
		// The loop then returns to accepting, so that
		// multiple connections may be served concurrently.
		go func(c net.Conn) {
			// Echo all incoming data.
			buf, err := ioutil.ReadAll(bufio.NewReader(conn))
			if err != nil {
				log.Println("Error trapped!")
				log.Fatal(err)
			}
			message := string(buf)
			fmt.Print(message)
			fmt.Printf("Message size(bytes): %d\n", len(message)-1)
			fmt.Print(sep)
			// Shut down the connection.
			c.Close()
		}(conn)
	}
}
