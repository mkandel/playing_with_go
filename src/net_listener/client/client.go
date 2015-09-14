package main

import (
    "flag"
    "fmt"
    "io/ioutil"
    "log"
    "net"
)

const port string = ":7777"

var content string

func main() {
    defer log.Print("Starting\n")
    var file = flag.String("f", "none", "file with contents to send")
    var debug = flag.Bool("d", false, "enable debug output")
    flag.Parse()
    if *debug == true {
        log.Printf("Debug enabled\n")
    }
    conn, err := net.Dial("tcp", "127.0.0.1"+port)
    if err != nil {
        log.Fatal(err)
    }
    content := "hello world (from Go interpreted)"
    if *file != "none" {
        buf, err := ioutil.ReadFile(*file)
        if err != nil {
            log.Fatal(err)
        }
        content = string(buf)
    }
    log.Printf("Filename: %s\n", *file)
    log.Printf("Contents: %s\n", content)
    fmt.Printf("Message size(bytes): %d\n", len(content))
    fmt.Fprintf(conn, content)
    //fmt.Fprintf(conn, content+"\n")
}
