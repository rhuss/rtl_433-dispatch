package main

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"log"
	"os"
)

type ConfigMap map[interface{}]interface{}

var Debug = false

func readConfigFile(filename string) (ConfigMap, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("Error reading file: %v\n", err)
	}

	var config ConfigMap
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("Error parsing YAML: %v\n", err)
	}

	if config["debug"] != nil && config["debug"] == true {
		Debug = true
	}
	return config, nil
}

// Get configuration by path
func getConfig(m ConfigMap, keys ...string) (interface{}, bool) {
	var ok bool
	var value interface{} = m

	for _, key := range keys {
		if value, ok = value.(ConfigMap)[key]; !ok {
			return nil, false
		}
	}

	return value, true
}

func getConfigString(m ConfigMap, key ...string) (string, bool) {
	if v, ok := getConfig(m, key...); ok {
		return v.(string), true
	}
	return "", false
}

func getConfigSlice(m ConfigMap, key ...string) ([]string, bool) {
	v, ok := getConfig(m, key...)
	if !ok {
		return nil, false
	}
	elements, ok := v.([]interface{})
	if !ok {
		return nil, false
	}
	strings := make([]string, len(elements))
	for i, v := range elements {
		var ok bool
		strings[i], ok = v.(string)
		if !ok {
			return nil, false
		}
	}
	return strings, true
}

func debug(args ...string) {
	if Debug {
		if len(args) > 1 {
			log.Printf(args[0], args[1:])
		} else {
			log.Println(args)
		}
	}
}
