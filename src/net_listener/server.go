package main

import (
    "bufio"
    "fmt"
    "io/ioutil"
    "log"
    "net"
)

const listen_port string = ":7777"
const sep string = "=====================================================================\n"

func main() {

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
