package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	mqtt "github.com/eclipse/paho.mqtt.golang"

	"github.com/lucasmelloec/mqtt-os-exec/internal/config"
	"github.com/lucasmelloec/mqtt-os-exec/internal/osExec"
)

var cfg config.Config
var topics []string

var messagePubHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	log.Printf("Received message: %s from topic: %s\n", msg.Payload(), msg.Topic())

	command := config.GetCommand(cfg.Topics, msg.Topic(), string(msg.Payload()))

	go osExec.Execute(command)
}

var connectHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
	log.Println("Connected")
	go subscribeTopics(client, topics)
}

var connectLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
	log.Printf("Connection lost: %v", err)
}

func subscribeTopics(client mqtt.Client, topics []string) {
	for _, topic := range topics {
		token := client.Subscribe(topic, 1, nil)
		token.Wait()
		log.Printf("Subscribed to topic %s\n", topic)
	}
}

func main() {
	cfg = config.GetConfig()

	opts := mqtt.NewClientOptions()
	opts.AddBroker(cfg.Broker)
	if cfg.Username != "" {
		opts.SetUsername(cfg.Username)
	}
	if cfg.Password != "" {
		opts.SetPassword(cfg.Password)
	}

	opts.SetAutoReconnect(true)

	opts.SetDefaultPublishHandler(messagePubHandler)
	opts.OnConnect = connectHandler
	opts.OnConnectionLost = connectLostHandler

	topics = config.HandleTopics(cfg.Topics)

	client := mqtt.NewClient(opts)

	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	// Reacting to signals (interrupt)
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt, syscall.SIGTERM)
	<-sigs

	client.Disconnect(250)
}
