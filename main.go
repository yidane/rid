package main

import (
	"bufio"
	"os"
	"strings"

	"errors"
	"flag"
	"time"

	"github.com/yidane/rid/command"
	"github.com/yidane/rid/context"
	"github.com/yidane/rid/log"
)

var ridContext *context.RidContext

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
		input := strings.ToLower(strings.Trim(string(data), ""))
		inputs := strings.Split(input, " ")
		cName := inputs[0]
		cArgs := inputs[1:]

		switch cName {
		case "help":
			command.Help(cArgs)
		case "exit":
			running = false
			log.Succeed("rid is being exit")
		default:
			command.Exec(ridContext, cName, cArgs...)
		}
	}

	log.Succeed("finish ", time.Now().String())
}

func login() error {
	uid := flag.String("uid", "", "What is your rid userid?")
	pwd := flag.String("pwd", "", "What is your rid password?")

	flag.Parse()

	if *uid == "" {
		flag.PrintDefaults()
		return errors.New("please use -uid to set rid userid")
	}
	if *pwd == "" {
		flag.PrintDefaults()
		return errors.New("please use -uid to set rid userid")
	}

	ridContext = context.NewRidContext()
	userInfo := context.UserInfo{UserID: *uid, Password: *pwd}
	return ridContext.Login(&userInfo)
}
