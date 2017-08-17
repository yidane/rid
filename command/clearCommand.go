package command

import (
	"github.com/yidane/rid/context"
	"github.com/yidane/rid/log"
)

type ClearCommand struct {
}

func (ClearCommand) Name() string {
	return "clear"
}

func (ClearCommand) Exec(ridContext *context.RidContext, args ...string) {
	if len(args) > 0 {
		log.Error("command clear do not need any arguments")
	}

	ridContext.ClearCache()
	log.Succeed("clear cache succeed")
}

func (ClearCommand) Usage() string {
	return ""
}
