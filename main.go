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
	commands := strings.Split(command, " ")
	commandName := commands[0]
	args := commands[1:]
	switch commandName {
	case "output":
		output(args)
	case "load":
		load(args)
	case "list":
		list(args)
	case "use":
		list(args)
	case "clear":
		clear(args)
	case "download":
		download(args)
	case "add":
		add(args)
	case "rm":
		rm(args)
	default:
		errorCommand(command)
	}
}
