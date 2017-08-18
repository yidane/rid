package command

import (
	"github.com/yidane/rid/context"
	"github.com/yidane/rid/log"
)

type OutCommand struct {
}

func (command OutCommand) Name() string {
	return "output"
}

func (command OutCommand) Exec(ridContext *context.RidContext, args ...string) {
	if len(args) != 1 {
		log.Error("command 'output' need just one argument which is a folder")
		return
	}

	if ridContext.CurrentDB == nil {
		log.Error("you should set current database with command 'use'")
		return
	}

	err := ridContext.SetOutput(args[0])
	if err != nil {
		log.Error(err)
	}
}

func (OutCommand) Usage() string {
	return ""
}
