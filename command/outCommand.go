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

func (command OutCommand) Exec(ricContext *context.RidContext, args ...string) {
	if len(args) != 1 {
		log.Error("command 'output' need just one argument which is a folder")
		return
	}
	ricContext.SetOutput(args[0])
}

func(OutCommand) Usage() string{
	return ""
}