package main

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"strings"

	"github.com/streadway/amqp"
)

var (
	debug = flag.Bool("d", false, "enable debug output")
)

func main() {
	var (
		queue  = flag.String("queue", "AMP_amon_decoder", "Queue to connect to")
		host   = flag.String("host", "localhost", "host (default: localhost)")
		port   = flag.String("port", ":5672", "port (default: 5672)")
		user   = flag.String("user", "some_user", "user")
		pass   = flag.String("pass", "some_password", "password")
		unsafe = flag.Bool("unsafe", false, "skip queue name pre-check (starts a hair faster)")

		list = flag.Bool("list", false, "list queues and exit")

		help = flag.Bool("h", false, "priint help message and exit")
	)
	flag.Parse()

	me := filepath.Base(os.Args[0])
	if *debug {
		fmt.Printf("%v: Starting\n", me)
		defer fmt.Printf("%v: Completed\n", me)
	}

	if *help {
		usage(me)
	} else if *list {
		queues := listQueues()
		for _, v := range queues {
			fmt.Println(v)
		}
	} else {
		hostUrl := *user + ":" + *pass + "@" + *host + *port
		url := "amqp://" + hostUrl
		if *debug {
			fmt.Println(hostUrl)
		}
		conn, err := amqp.Dial(url)
		if err != nil {
			log.Fatal(err)
		}
		defer conn.Close()

		channel, err := conn.Channel()
		if err != nil {
			log.Fatal(err)
		}
		defer channel.Close()

		// Check to see if the queue exists before attempting to connect
		// Not sure if unsafe buys me anything, TODO - time the connect
		// UPDATE: unsafe is ~25% faster:
		//        s/iter   safe unsafe
		// safe     2.57     --   -19%
		// unsafe   2.07    24%     --

		if !(*unsafe) {
			queues := listQueues()
			if queues == nil {
				log.Fatal("Got empty queues ...")
			}

			found := false
			for _, v := range queues {
				// EqualFold checks exact equality in a UTF-8 world
				if strings.EqualFold(v, *queue) {
					found = true
				}
			}
			if found != true {
				log.Fatalf("Queue %v not found!", *queue)
			}
		}
		// hard-coding AMP_exchange as the exchange for now
		err = channel.QueueBind(*queue, *queue, "AMP_exchange", false, nil)
		if err != nil {
			log.Fatal(err)
		}
		msgs, err := channel.Consume(
			*queue,            // queue
			"rabbbitmq_snoop", // consumer name
			true,              // auto-ack
			false,             // exclusive
			false,             // no-local
			false,             // no-wait
			nil,               // args
		)
		if err != nil {
			log.Fatalf("Failed to register: \n%s", err)
		}

		forever := make(chan bool)

		go func() {
			for d := range msgs {
				log.Printf("%s\n", d.Body)
			}
		}()

		fmt.Printf("Snooping on %v:\n", *queue)
		<-forever
	}

}

func listQueues() []string {
	var ret []string

	u, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	cmd := exec.Command("/opt/airwave/sbin/rabbitmqctl", "list_queues")
	if u.Username != "root" {
		cmd = exec.Command("sudo", "/opt/airwave/sbin/rabbitmqctl", "list_queues")
	}

	var out bytes.Buffer
	cmd.Stdout = &out
	err = cmd.Run()
	if err != nil {
		log.Fatal(err)
	}

	queues := strings.Split(out.String(), "\n")
	for _, v := range queues {
		if strings.Contains(v, "Listing queues ...") {
			continue
		}
		fields := strings.Fields(v)
		if len(fields) == 0 {
			continue
		}
		if *debug {
			fmt.Printf("Adding: %s\n", fields[0])
		}
		ret = append(ret, fields[0])
	}
	return ret
}

func usage(name string) {
	fmt.Printf("Usage: %s [OPTION]\n", name)
	fmt.Printf("Options:\n")
	fmt.Println("    --queue    queue name")
	fmt.Println("    --host     host to connect to (default: localhost)")
	fmt.Println("    --port     port to connect to (default: \":5672\")")
	fmt.Println("    --user     user to connect to RabbitMQ")
	fmt.Println("    --pass     password to connect to RabbitMQ")
	fmt.Println("")
	fmt.Println("    --unsafe   disable queue name check, possible connect speed increase")
	fmt.Println("")
	fmt.Println("    --list     list queues and exit")
	fmt.Println("")
	fmt.Println("    -d         enable possible extra debug output")
	fmt.Println("    -h         Print this help message and exit")
	return
}
