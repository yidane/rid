package command

import (
	"fmt"

	"github.com/yidane/rid/context"
	"github.com/yidane/rid/log"
)

type ShowCommand struct {
}

func (ShowCommand) Name() string {
	return "show"
}

func (ShowCommand) Exec(ridContext *context.RidContext, args ...string) {
	if ridContext.CurrentDB == nil {
		log.Error("current database is nil")
		return
	}

	log.Succeed(fmt.Sprintf("current database is `%s`", ridContext.CurrentDB.Name))
}

func (ShowCommand) Usage() string {
	return "show:show current database"
}

func init() {
	packageCommand(ShowCommand{})
}
