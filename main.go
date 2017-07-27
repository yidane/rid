package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/yidane/rid/log"
)

var ridClient *RidClient

func main() {
	err := login()
	if err != nil {
		log.Error(err)
		return
	}

	log.Info("login rid successful!")

	running := true
	reader := bufio.NewReader(os.Stdin)
	for running {
		data, _, _ := reader.ReadLine()
		if len(data) == 0 {
			continue
		}
		command := strings.ToLower(strings.Trim(string(data), ""))

		if len(command) < 2 {
			errorCommand(command)
			continue
		}

		if command == "help" {
			help()
		} else if command == "exit" {
			exit()
			running = false
		} else {
			handCommand(command)
		}
	}

	fmt.Println("rid")
}

func handCommand(command string) {
	args := strings.Split(command, " ")
	switch args[0] {
	case "output":
		output(args[0:])
	case "load":
		load(args[0:])
	case "list":
		list(args[0:])
	case "use":
		list(args[0:])
	case "clear":
		clear(args[0:])
	case "download":
		download(args[0:])
	case "add":
		add(args[0:])
	case "rm":
		rm(args[0:])
	default:
		errorCommand(command)
	}
}
