package main

import (
	"flag"
	zmq "github.com/pebbe/zmq4"
	"log"
	"syscall"
)

var frontendLocation = flag.String("remote", "tcp://1.3.3.7:1234", "Where we're proxying to.")
var backendLocation = flag.String("bind", "tcp://0.0.0.0:1234", "Where we're binding to.")

func main() {
	flag.Parse()

	log.Println("Starting...")
	frontend, err := zmq.NewSocket(zmq.XSUB)
	if err != nil {
		log.Fatalln("Could not open socket to remote system:", err)
	}
	defer frontend.Close()
	frontend.Connect(*frontendLocation)

	backend, err := zmq.NewSocket(zmq.XPUB)
	if err != nil {
		log.Fatalln("Could not bind to socket:", err)
	}
	defer backend.Close()
	backend.Bind(*backendLocation)

	log.Println("Entering proxy...")
	for {
		err = zmq.Proxy(frontend, backend, nil)
		if err != nil {
			if zmq.AsErrno(err) == zmq.Errno(syscall.EINTR) {
				log.Println("Got interrupted:", err)
			} else {
				log.Fatalln("Proxy error:", err)
			}
		}
	}
}
