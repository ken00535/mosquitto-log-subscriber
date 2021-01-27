package client

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

// Subscribe message
func Subscribe(host Host) {
	_, err := newClient(host)
	if err != nil {
		panic(err)
	}
	for {
		time.Sleep(time.Second)
	}
}

func subscribeTopic(client mqtt.Client) {
	token := client.Subscribe("$SYS/broker/log/#", 0, func(_ mqtt.Client, msg mqtt.Message) {
		log.Infoln(msg.Topic())
		payload := string(msg.Payload())
		segments := strings.Split(payload, ":")
		i, err := strconv.ParseInt(segments[0], 10, 64)
		if err != nil {
			panic(err)
		}
		tm := time.Unix(i, 0)
		log.Infoln(tm.Format("2006-01-02 15:04:05  ") + payload)
	})
	if token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
	}
}
