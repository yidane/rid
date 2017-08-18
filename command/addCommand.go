package command

import (
	"fmt"
	"github.com/yidane/rid/context"
	"github.com/yidane/rid/log"
)

type AddCommand struct {
}

func (AddCommand) Name() string {
	return "add"
}

func (AddCommand) Exec(ridContext *context.RidContext, args ...string) {
	if ridContext.CurrentDB == nil {
		log.Error("you should set current database with command 'use'")
		return
	}
	if len(args) == 0 {
		log.Error("arguments should be at least one")
		return
	}

	reg, err := getRegex(args[0])
	if err != nil {
		log.Error(err)
		return
	}

	allTables, err := ridContext.LoadTables(ridContext.CurrentDB.Name)
	if err != nil {
		log.Error(err)
		return
	}

	fTable := []string{}
	for i := 0; i < len(allTables); i++ {
		if reg.MatchString(allTables[i]) {
			fTable = append(fTable, allTables[i])
		}
	}

	if len(fTable) == 0 {
		log.Warn(fmt.Sprint("no table like [", args[0], "]"))
		return
	}

	ridContext.AddToCache(fTable...)
	log.Succeed("such tables have added to cache :\r", fTable)
}

func (AddCommand) Usage() string {
	return ""
}
