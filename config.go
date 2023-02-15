package main

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"os"
)

type ConfigMap map[interface{}]interface{}

func readConfigFile(filename string) (ConfigMap, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("Error reading file: %v\n", err)
	}

	var config ConfigMap
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("Error parsing YAML: %v\n", err)
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
