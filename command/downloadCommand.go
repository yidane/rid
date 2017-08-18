package command

import (
	"github.com/yidane/rid/context"
	"github.com/yidane/rid/log"
)

type DownloadCommand struct {
}

func (DownloadCommand) Name() string {
	return "download"
}

func (DownloadCommand) Exec(ridContext *context.RidContext, args ...string) {
	if len(args) > 0 {
		log.Error("command download does not need any argument")
		return
	}

	err := ridContext.DownloadAll()
	if err != nil {
		log.Error(err)
		return
	}

	log.Succeed("download scripts succeed")
}

func (DownloadCommand) Usage() string {
	return ""
}
