package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func main() {

	type initialSetup struct {
		host     string
		clientID string
		username string
		password string
		topic    string
	}

	connectionDetails := initialSetup{host: "tcp://localhost:1883", clientID: "myDesktopClient", username: "3", password: "880b5c97-6a55-4dae-8f98-5b0cd74aac5a", topic: "channels/1/messages"}
	client := connectToMQTTServer(connectionDetails.host, connectionDetails.clientID, connectionDetails.username, connectionDetails.password)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "UP")
	})

	http.HandleFunc("/send_message", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Sending Message")
		publishToTopic(client, connectionDetails.topic, "hello from go")
	})

	http.HandleFunc("/command/on", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Sending Message")
		publishToTopic(client, connectionDetails.topic, "1")
	})

	http.HandleFunc("/command/off", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Sending Message")
		publishToTopic(client, connectionDetails.topic, "0")
	})

	subscribeToTopic(client, connectionDetails.topic)
	log.Fatal(http.ListenAndServe(":8081", nil))
}

/******** Connect to MQTT Server on Mainflux *********/
func connectToMQTTServer(host string, clientID string, username string, password string) mqtt.Client {
	opts := mqtt.NewClientOptions().AddBroker(host).SetClientID(clientID)
	opts.SetUsername(username)
	opts.SetPassword(password)

	c := mqtt.NewClient(opts)
	if token := c.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
	return c
}

/**** Subscribe/Unsubscribe to Topic and Message handler ****/
var messageHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	fmt.Println("GO -----> :", string(msg.Payload()))
}

func subscribeToTopic(c mqtt.Client, topic string) {
	if token := c.Subscribe(topic, 1, messageHandler); token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
		os.Exit(1)
	}
}

func unsubscribeFromTopic(c mqtt.Client, topic string) {
	if token := c.Unsubscribe(topic); token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
		os.Exit(1)
	}
}

/******** Publish to Topic *********/
func publishToTopic(c mqtt.Client, topic string, message string) {
	token := c.Publish("channels/1/messages", 1, false, message)
	token.Wait()
}
