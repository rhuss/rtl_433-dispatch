package main

import (
	"encoding/json"
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"log"
	"strconv"
	"strings"
)

type DispatcherFunc func(mqtt.Client, ConfigMap, RtlSdrMessage) error

var dispatcherMap map[string]DispatcherFunc

func init() {
	dispatcherMap = map[string]DispatcherFunc{
		"wireless-mbus": processWmBus,
	}
}

func processWmBus(client mqtt.Client, config ConfigMap, message RtlSdrMessage) error {
	if config == nil {
		return fmt.Errorf("Wireless-MBus configuration required for wmbus extractions")
	}
	decodedMessage, err := decodeWmBusTelegram(config, message)
	if err != nil {
		return err
	}
	topic, exists := config["topic"].(string)
	if !exists {
		topic = "watermeter"
	}
	return publish(client, topic, decodedMessage)
}

func decodeWmBusTelegram(config ConfigMap, message RtlSdrMessage) (RtlSdrMessage, error) {
	data := RtlSdrMessage{}
	device, ok := getConfigString(config, "device")
	if !ok {
		return data, fmt.Errorf("no device configured for wmbus")
	}
	key, ok := getConfigString(config, "key")
	if !ok {
		return data, fmt.Errorf("no key configured for wmbus")
	}
	name, ok := getConfigString(config, "name")
	if !ok {
		name = "watermeter"
	}
	telegram, ok := message["data"]
	if !ok {
		return RtlSdrMessage{}, fmt.Errorf("no telegram data found in message: %v", message)
	}
	id, ok := message["id"]
	if !ok {
		return data, fmt.Errorf("no id found in message: %v", message)
	}
	idS := strconv.FormatFloat(id.(float64), 'f', 0, 64)
	args := []string{
		"./wmbusmeters",
		"--format=json",
		telegram.(string),
		name,
		device,
		idS,
		key}
	log.Printf("Exec: %v", args)
	telegramDecode, err := execute(args)

	if err != nil {
		return RtlSdrMessage{}, err
	}

	err = json.Unmarshal([]byte(telegramDecode), &data)
	if err != nil {
		return data, err
	}
	data["name"] = name
	data["id"] = idS
	data["meter"] = device

	return data, nil
}

func dispatch(client mqtt.Client, config ConfigMap, message RtlSdrMessage) error {

	model, ok := message["model"].(string)
	if !ok {
		return fmt.Errorf("unknown message received: %v", message)
	}
	model = strings.ToLower(model)

	var modelConfig ConfigMap = nil
	sensorCfg, exists := config["sensors"]
	if exists {
		c, ok := getConfig(sensorCfg.(ConfigMap), model)
		if ok {
			modelConfig = c.(ConfigMap)
		}
	}

	// Check for special handling for the given model
	dispFunc, exists := dispatcherMap[model]
	if exists {
		return dispFunc(client, modelConfig, message)
	}

	topic := model
	if modelConfig != nil && modelConfig["topic"] != nil {
		topic = modelConfig["topic"].(string)
	}
	return publish(client, topic, message)
}
