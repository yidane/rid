package command

import (
	"fmt"

	"github.com/yidane/rid/context"
	"github.com/yidane/rid/log"
)

type LoadCommand struct {
}

func (LoadCommand) Name() string {
	return "load"
}

func (LoadCommand) Exec(ricContext *context.RidContext, args ...string) {
	fmt.Println(args)

	//load all database
	if len(args) == 0 {
		dbArr, err := ricContext.LoadDataBase()
		if err != nil {
			log.Error(err)
			return
		}

		log.Succeed(dbArr)
		return
	}
	// load all tables of some database
	ricContext.LoadTables(args[0])
}

func (LoadCommand) Usage() string {
	return ""
}
