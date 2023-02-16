package main

import (
	"encoding/json"
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"time"
)

func connectMqtt(config ConfigMap) (mqtt.Client, error) {
	// Connect to MQTT broker
	c, exist := getConfig(config, "mqtt")
	if !exist {
		return nil, fmt.Errorf("no mqtt configuration")
	}
	brokerConfig, _ := c.(ConfigMap)

	connect, ok := getConfigString(brokerConfig, "broker")
	if !ok {
		return nil, fmt.Errorf("no mqtt connect string given")
	}
	opts := mqtt.NewClientOptions().AddBroker(connect)

	if user, ok := getConfigString(brokerConfig, "user"); ok {
		opts.SetUsername(user)
		if password, ok := getConfigString(brokerConfig, "password"); ok {
			opts.SetPassword(password)
		}
	}
	mqttClient := mqtt.NewClient(opts)
	if token := mqttClient.Connect(); token.Wait() && token.Error() != nil {
		return nil, token.Error()
	}

	return mqttClient, nil
}

func publish(client mqtt.Client, topic string, payload RtlSdrMessage) error {
	// serialize the map to JSON
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	// Publish message to MQTT broker
	debug("Sent: %s", string(jsonData))
	token := client.Publish(topic, 0, false, jsonData)
	token.WaitTimeout(10 * time.Second)
	return token.Error()
}
