package command

import "github.com/yidane/rid/context"

type RemoveCommand struct {
}

func (RemoveCommand) Name() string {
	return "rm"
}

func (RemoveCommand) Exec(ridContext *context.RidContext, args ...string) {
	ridContext.RemoveFromCache(args...)
}

func (RemoveCommand) Usage() string {
	return "rm [table];remove table from cache"
}
