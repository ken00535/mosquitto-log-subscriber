package main

import (
	"fmt"
	"mosquitto/log/pkg/client"
	"strconv"
	"strings"
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
		payload := string(msg.Payload())
		segments := strings.Split(payload, ":")
		i, err := strconv.ParseInt(segments[0], 10, 64)
		if err != nil {
			panic(err)
		}
		tm := time.Unix(i, 0)
		fmt.Println(tm.Format("2006-01-02 15:04:05  ") + payload)
	})
	if token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
	}
	for {
		time.Sleep(time.Second)
	}
}
