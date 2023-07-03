package main

import (
	"bufio"
	"fmt"
	"log"
	"os/exec"
	"strings"
)

const maxRetries = 5

func startPipe(args []string) (chan string, chan error) {
	outChan := make(chan string)
	errChan := make(chan error)

	go pipe(args, outChan, errChan)
	return outChan, errChan
}

func pipe(args []string, outChan chan string, errChan chan error) {

	retries := 0
	for {
		retries++
		if retries > maxRetries {
			errChan <- fmt.Errorf("broken pipe: max retries %d reached", maxRetries)
			return
		}

		scanner, cmd, err := runCommand(args)
		if err != nil {
			log.Printf("error executing command '%s' : %v", strings.Join(args, " "), err)
			continue
		}
		for scanner.Scan() {
			outChan <- scanner.Text()
		}

		if err := scanner.Err(); err != nil {
			log.Printf("error while scanning content: '%s' : %v", strings.Join(args, " "), err)
			continue
		}

		err = cmd.Wait()
		if err != nil {
			log.Printf("error waiting for command: '%s' : %v", strings.Join(args, " "), err)
			continue
		}
	}
}

func runCommand(args []string) (*bufio.Scanner, *exec.Cmd, error) {
	var extraArgs []string
	if len(args) > 1 {
		extraArgs = args[1:]
	}
	cmd := exec.Command(args[0], extraArgs...)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, nil, fmt.Errorf("cannot create pipe: %v (args: %v)", err, args)
	}
	scanner := bufio.NewScanner(stdout)

	err = cmd.Start()
	if err != nil {
		return nil, nil, fmt.Errorf("cannot start command: %v (args: %v)", err, args)
	}
	return scanner, cmd, nil
}

func execute(args []string) (string, error) {
	scanner, cmd, err := runCommand(args)

	if err != nil {
		return "", err
	}

	var lines []string
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return "", err
	}

	err = cmd.Wait()
	if err != nil {
		return "", err
	}
	return strings.Join(lines, ""), nil
}
