package command

import (
	"bytes"
	"fmt"

	"github.com/yidane/rid/context"
	"github.com/yidane/rid/log"
)

type LoadCommand struct {
}

func (LoadCommand) Name() string {
	return "load"
}

func (LoadCommand) Exec(ridContext *context.RidContext, args ...string) {
	fmt.Println(args)

	//load all database
	if len(args) == 0 {
		loadAllDBList(ridContext)
	} else {
		// load all tables of some database
		loadTables(ridContext, args[0])
	}
}

func loadAllDBList(ridContext *context.RidContext) {
	dbArr, err := ridContext.LoadDataBase()
	if err != nil {
		log.Error(err)
		return
	}
	log.Succeed(fmt.Sprintf("load %d databases as follows:", len(dbArr)))
	outPut(dbArr)
}

func loadTables(ridContext *context.RidContext, dbName string) {
	tables, err := ridContext.LoadTables(dbName)
	if err != nil {
		log.Error(err)
		return
	}
	log.Succeed(fmt.Sprintf("load %d tables as follows:", len(tables)))
	outPut(tables)
}

func outPut(strArr []string) {
	outStr := bytes.Buffer{}
	l := 0
	sub := len(strArr) - l*3
	for sub > 0 {
		switch {
		case sub > 3:
			outStr.WriteString("	|%-50s|%-50s|%-50s|\n")
		case sub == 2:
			outStr.WriteString("	|%-50s|%-50s|\n")
		case sub == 1:
			outStr.WriteString("	|%-50s|\n")
		}
		l++
		sub = len(strArr) - l*3
	}

	ins := []interface{}{}
	for i := 0; i < len(strArr); i++ {
		ins = append(ins, fmt.Sprintf("`%s`", strArr[i]))
	}

	fmt.Printf(outStr.String(), ins...)
}

func (LoadCommand) Usage() string {
	return "load [database];load the databases or tables from remote server."
}

func init() {
	packageCommand(LoadCommand{})
}
