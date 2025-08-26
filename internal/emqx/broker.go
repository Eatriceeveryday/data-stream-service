package emqx

import (
	"fmt"
	"time"

	"github.com/Eatriceeveryday/data-stream-service/internal/config"
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

var messagePubHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	fmt.Printf("Received message: %s from topic: %s\n", msg.Payload(), msg.Topic())
}

var connectHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
	fmt.Println("Connected")
}

var connectLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
	fmt.Printf("Connect lost: %v", err)
}

func ConnectToClient(cfg *config.Config) (mqtt.Client, error) {
	clientOpt := mqtt.NewClientOptions()
	clientOpt.AddBroker(fmt.Sprintf("tcp://%s:%s", cfg.EMQX_host, cfg.EMQX_port))
	clientOpt.SetClientID(fmt.Sprintf("client-%d", time.Now().Unix()))
	clientOpt.SetUsername(cfg.ApiKey)
	clientOpt.SetPassword(cfg.ApiKey)
	clientOpt.SetDefaultPublishHandler(messagePubHandler)
	clientOpt.OnConnect = connectHandler
	clientOpt.OnConnectionLost = connectLostHandler
	client := mqtt.NewClient(clientOpt)

	if token := client.Connect(); token.Wait() && token.Error() != nil {
		return nil, token.Error()
	}

	return client, nil
}
