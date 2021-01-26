package main

import (
	"fmt"
	"mosquitto/log/pkg/client"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

// Subscribe message
func Subscribe(host client.Host) {
	client, err := client.NewClient(host)
	if err != nil {
		panic(err)
	}
	token := client.Subscribe("$SYS/broker/log/#", 0, func(_ mqtt.Client, msg mqtt.Message) {
		fmt.Println(msg.Topic())
		fmt.Println(string(msg.Payload()))
	})
	if token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
	}
	for {
		time.Sleep(time.Second)
	}
}
