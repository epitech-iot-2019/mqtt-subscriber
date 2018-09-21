package main

import (
	"fmt"
	"os"
	"os/signal"

	"github.com/yosssi/gmq/mqtt"
	"github.com/yosssi/gmq/mqtt/client"
)

func main() {
	// Set up channel on which to send signal notifications.
	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, os.Interrupt, os.Kill)

	// Create an MQTT Client.
	cli := client.New(&client.Options{
		// Define the processing of the error handler.
		ErrorHandler: func(err error) {
			fmt.Println(err)
		},
	})

	// Terminate the Client.
	defer cli.Terminate()

	// Connect to the MQTT Server.
	err := cli.Connect(&client.ConnectOptions{
		Network:  "tcp",
		Address:  "broker.hivemq.com:1883",
		ClientID: []byte("ESP8266_Oui-subscriber"),
	})
	if err != nil {
		panic(err)
	}

	fmt.Println("Client init")

	// Subscribe to topics.
	err = cli.Subscribe(&client.SubscribeOptions{
		SubReqs: []*client.SubReq{
			&client.SubReq{
				TopicFilter: []byte("/ESP8266_Oui/value"),
				QoS:         mqtt.QoS1,
				Handler: func(topicName, message []byte) {
					fmt.Println("Got a new message !")
					fmt.Println(string(topicName), string(message))
				},
			},
		},
	})
	if err != nil {
		panic(err)
	}

	fmt.Println("Topic subscribed")

	// // Unsubscribe from topics.
	// err = cli.Unsubscribe(&client.UnsubscribeOptions{
	// 	TopicFilters: [][]byte{
	// 		[]byte("foo"),
	// 	},
	// })
	// if err != nil {
	// 	panic(err)
	// }

	// Wait for receiving a signal.
	<-sigc

	// Disconnect the Network Connection.
	if err := cli.Disconnect(); err != nil {
		panic(err)
	}
	fmt.Println("Client disconnected")
}
