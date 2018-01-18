package command

import (
	"github.com/yidane/rid/context"
	"github.com/yidane/rid/log"
)

type FindCommand struct {
}

func (FindCommand) Name() string {
	return "find"
}

func (FindCommand) Exec(ridContext *context.RidContext, args ...string) {
	if ridContext.CurrentDB == nil {
		log.Error("you should set current database with command 'use'")
		return
	}

	if len(args) != 1 {
		log.Error("argument should be only one")
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

	fTables := []string{}
	for i := 0; i < len(allTables); i++ {
		if reg.MatchString(allTables[i]) {
			fTables = append(fTables, allTables[i])
		}
	}

	if len(fTables) == 0 {
		log.Warn("no such table")
		return
	}

	log.Succeed("find some tables :\r", fTables)
}

func (FindCommand) Usage() string {
	return "find [table];find some table"
}
