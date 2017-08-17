package command

import (
	"github.com/yidane/rid/log"
	"github.com/yidane/rid/context"
)

type UseCommand struct {
}

func (UseCommand) Name() string {
	return "use"
}

func (UseCommand) Exec(ricContext *context.RidContext, args ...string) {
	if len(args) != 1 {
		log.Error("command 'use' need just one argument of database name")
		return
	}

	dbName := args[0]
	err := ricContext.SetCurrentDatabase(dbName)
	if err != nil {
		log.Error(err)
	}

	log.Succeed("load database ", dbName, " success")
}

func (UseCommand) Usage() string {
	return ""
}
