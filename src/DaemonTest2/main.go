package main

import (
	"encoding/json"
	"flag"
	"io"
	"log"
	"os"
	"os/signal"
	"sync"
	"time"
)

const version = "0.0.1"
const apiVers = "v1"
const sep string = "=======================================\n"
const logFile = "/var/log/DaemonTest.log"
const passFile = "/var/airwave/secure/security"
const appName = "DaemonTest"

var port = 7777
var apiport = 7778

func Init() {
	// long and short options: --port|-p
	//                         --apiport|-a
	flag.IntVar(&port, "port", 7777, "Listen port")
	flag.IntVar(&port, "p", 7777, "Listen port")
	flag.IntVar(&apiport, "apiport", 7778, "Listen port")
	flag.IntVar(&apiport, "a", 7778, "Listen port")
}

func main() {
	startTime := time.Now()
	// Setup logging
	f, err := os.OpenFile(logFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	defer f.Close()
	log.SetOutput(io.Writer(f))
	log.Print(sep)
	log.Printf("%s: Starting", appName)

	var srv AppServer
	var api APIServer
	var wg sync.WaitGroup

	defer func() {
		log.Print(sep)
		elapsedTime := time.Now().Sub(startTime)
		log.Println("=", appName, "ran for", elapsedTime)
		var processed = srv.ProcessedCount()
		log.Println("=", "Processed", processed, "messages")
		log.Println("=", float64(processed)/elapsedTime.Seconds(), "per second")
		log.Printf("= %s: Exiting", appName)
		log.Print(sep)
	}()

	// Trap ctrl-c to exit nicely
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	go func() {
		for _ = range signalChan {
			log.Println("Received an interrupt, stopping services...")
			// Kill the waits for the 2 servers? ...
			wg.Done()
			wg.Done()
		}
	}()

	// Commandline args
	flag.Parse()

	/*
		    // For later
				rabbitpass := parseRabbitPass()
				log.Printf("Rabbit password: '%s'", rabbitpass)
	*/

	// and start our servers
	wg.Add(2)
	go srv.Start(port, wg)
	go api.Start(apiport, wg)
	println("Server:", srv.Version(), "   API:", api.Version())
	wg.Wait()
	log.Println("WaitGroup returned!!")
}

var rabbitPass struct {
	RabbitMQPassword string `json:"rabbitmq_password"`
}

func parseRabbitPass() string {
	f, err := os.Open(passFile)
	if err != nil {
		log.Println("opening config file", err.Error())
	}
	defer f.Close()

	jsonParser := json.NewDecoder(f)
	if err = jsonParser.Decode(&rabbitPass); err == io.EOF {
		//break
	} else if err != nil {
		log.Println("parsing config file", err.Error())
	}

	return rabbitPass.RabbitMQPassword
}
