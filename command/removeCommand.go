package command

import "github.com/yidane/rid/context"

type RemoveCommand struct {
}

func (RemoveCommand) Name() string {
	return "rm"
}

func (RemoveCommand) Exec(ricContext *context.RidContext, args ...string) {

}

func (RemoveCommand) Usage() string {
	return ""
}
