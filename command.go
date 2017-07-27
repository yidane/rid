package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/yidane/rid/log"
)

func output(args []string) {
	if len(args) != 1 {
		log.Error("command 'output' need just one argument which is a folder")
		return
	}
	ridClient.SetOutput(args[0])
}

func load(args []string) {
	if len(args) > 1 {
		log.Error("command 'load' need at most one argument which is a folder")
	}

	if len(args) == 0 { //load all database
		ridClient.LoadDataBase()
	} else { // load all tables of some database
		if ridClient.CurrentDB == nil {
			log.Warn("you need shoose database first by using command 'use'")
			return
		}

		ridClient.LoadTables(args[0])
	}
}

func list(args []string) {
	if len(args) != 1 {
		log.Error("command 'list' does not need any argument which is a folder")
		return
	}

	result := ridClient.SelectedTables()
	log.Info(result)
}

func use(args []string) {
	if len(args) != 1 {
		log.Error("command 'use' need just one argument which is a folder")
		return
	}

	ridClient.SetCurrentDatabase(args[0])
}

func clear(args []string) {
	if len(args) != 1 {
		log.Error("command 'use' does not need any argument which is a folder")
		return
	}
}

func add(args []string) {
	if len(args) != 1 {
		log.Error("command 'use' need just one argument which is a folder")
		return
	}
}

func download(args []string) {

}

func rm(args []string) {

}

func errorCommand(command string) {
	log.Error("unkonwn command '", command, "' please input 'help' to show all command")
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

	ridClient = &RidClient{}
	return ridClient.Login(*uid, *pwd)
}

func help() {
	fmt.Println("show all command")
}

func exit() {
	ridClient = nil
	log.Info("exit successfully!")
	time.Sleep(time.Second)
	os.Exit(0)
}
