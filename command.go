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
		log.Error("command 'load' need at most one argument,you can input 'load' to show all databases,or you can input 'load dbname' to show all tables of dbname")
	}

	if len(args) == 0 { //load all database
		dbArr, err := ridClient.LoadDataBase()
		if err != nil {
			log.Error(err)
		}

		log.Succeed(dbArr)
	} else { // load all tables of some database
		fmt.Println(args[0])
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
	if len(args) > 1 {
		log.Error("command 'download' do not need any argument")
		return
	}
	if ridClient.CurrentDB == nil {
		log.Error("you must choose one database first by using command 'use dbname'")
		return
	}
	if ridClient.selectedTables == nil || len(ridClient.selectedTables) == 0 {
		log.Error("you must select some table of the choosed database  by using command 'add table'")
		return
	}
	if ridClient.Output == "" {
		log.Error("you must set output fold before download by using command 'add table'")
		return
	}

	ridClient.DownloadAll()
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
