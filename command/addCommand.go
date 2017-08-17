package command

import (
	"github.com/yidane/rid/context"
	"github.com/labstack/gommon/log"
)

type AddCommand struct {
}

func (AddCommand) Name() string {
	return "add"
}

func (AddCommand) Exec(ridContext *context.RidContext, args ...string) {
	if ridContext.CurrentDB==nil{
		log.Error("you should set current database with command 'use'")
	}
}

func(AddCommand) Usage() string{
	return ""
}

