package command

import "github.com/yidane/rid/context"

type DownloadCommand struct {
}

func (DownloadCommand) Name() string {
	return "download"
}

func (DownloadCommand) Exec(ricContext *context.RidContext, args ...string) {
	panic("implement me")
}

func(DownloadCommand) Usage() string{
	return ""
}
