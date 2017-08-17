package command

import (
	"github.com/yidane/rid/context"
	"github.com/labstack/gommon/log"
)

type ListCommand struct {
}

func (ListCommand) Name() string {
	return "list"
}

func (ListCommand) Exec(ridContext *context.RidContext, args ...string) {
	if len(args)>0{
		log.Error("command list do not need argument")
	}

	var tables=ridContext.SelectedTables()
	if len(tables)==0{
		log.Info("nothing be selected")
	}
	log.Info(tables)
}

func(ListCommand) Usage() string{
	return ""
}