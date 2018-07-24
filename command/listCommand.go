package command

import (
	"github.com/labstack/gommon/log"
	"github.com/yidane/rid/context"
)

type ListCommand struct {
}

func (ListCommand) Name() string {
	return "list"
}

func (ListCommand) Exec(ridContext *context.RidContext, args ...string) {
	if len(args) > 0 {
		log.Error("command list do not need argument")
		return
	}

	var tables = ridContext.SelectedTables()
	if len(tables) == 0 {
		log.Info("nothing is selected")
		return
	}
	log.Info("such tables have added:\r", tables)
}

func (ListCommand) Usage() string {
	return "list;list all the tables in cache"
}

func init() {
	packageCommand(ListCommand{})
}
