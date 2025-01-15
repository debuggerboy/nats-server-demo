package main

/*******************************************************
Bringing up a NATS server embedded into the application.
The NATS server runs as an internal process rather than
listening on the network.
*******************************************************/

import (
	"log"
	"time"

	"github.com/nats-io/nats-server/v2/server"
	"github.com/nats-io/nats.go"
)

func main() {

	initiateNatsServer()

}

func initiateNatsServer() {
	nc, ns := initNatsServer()

	subject := "example.subject"
	message := "Hello, Embedded NATS!"

	nc.Publish(subject, []byte(message))

	log.Printf("Published message to subject %s: %s\n", subject, message)

	nc.Subscribe(subject, func(msg *nats.Msg) {
		log.Printf("Received message on subject %s: %s\n", msg.Subject, string(msg.Data))
	})

	ns.WaitForShutdown()
}

func initNatsServer() (*nats.Conn, *server.Server) {
	opts := &server.Options{
		DontListen: true,
	}

	ns, err := server.NewServer(opts)
	if err != nil {
		log.Fatalln("Error while creating NATS server:", err)
	}

	ns.ConfigureLogger()

	go ns.Start()

	if !ns.ReadyForConnections(5 * time.Second) {
		log.Fatalln("Error during starting NATS server:", err)
	}

	clientOpts := []nats.Option{}

	clientOpts = append(clientOpts, nats.InProcessServer(ns))

	nc, err := nats.Connect(ns.ClientURL(), clientOpts...)
	if err != nil {
		log.Fatalln("Error while trying to connect to NATS server:", err)
	}

	return nc, ns
}
