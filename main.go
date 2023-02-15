package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

type RtlSdrMessage map[string]interface{}

func main() {

	if len(os.Args) != 2 {
		log.Fatalf("Usage: %s <config>\n", os.Args[0])
	}

	config, err := readConfigFile(os.Args[1])
	if err != nil {
		log.Fatalf("can't read config file: %v", err)
	}

	mqttClient, err := connectMqtt(config)
	if err != nil {
		log.Fatalf("can't connect to mqtt: %v", err)
	}

	args, ok := getConfigSlice(config, "command")
	if !ok {
		log.Fatalf("no command configured")
	}
	inChan, errChan := startPipe(args)

	for {
		select {
		case line := <-inChan:
			data := RtlSdrMessage{}
			err := json.Unmarshal([]byte(line), &data)
			if err != nil {
				fmt.Println("error parsing json:", err)
				continue
			}
			err = dispatch(mqttClient, config, data)
			if err != nil {
				log.Printf("error while dispatching: %v\n", err)
			}
		case err := <-errChan:
			mqttClient.Disconnect(250)
			log.Fatalf("command stopped: %v", err)
		}
	}
}
