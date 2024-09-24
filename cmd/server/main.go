package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/doovel/bootdev/pubsub/learn-pub-sub-starter/internal/pubsub"
	"github.com/doovel/bootdev/pubsub/learn-pub-sub-starter/internal/routing"
	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	fmt.Println("Starting Peril server...")

	const connection = "amqp://guest:guest@localhost:5672/"

	conn, err := amqp.Dial(connection)
	if err != nil {
		fmt.Println("unable to connect RabbitMQ: ", err)
	}

	fmt.Println("Connection established to ${connection}")

	defer func() {
		fmt.Println("closing connection")
		conn.Close()
	}()

	ch, err := conn.Channel()
	if err != nil {
		fmt.Println("unable to open channel: ", err)
		os.Exit(1)
	}
	defer ch.Close()

	state := routing.PlayingState{isPaused: true}
	err = pubsub.PublishJSON(ch, routing.ExchangePerilDirect, routing.PauseKey, state)
	if err != nil {
		fmt.Println("Unable to publish JSON: ", err)
		os.Exit(1)
	}

	fmt.Println("Published paused message to exchange.")

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	sig := <-sigCh

	fmt.Println("received signal: ", sig)
}
