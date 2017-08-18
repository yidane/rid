package command

import (
	"fmt"
	"github.com/yidane/rid/context"
	"testing"
)

func Test_Help(t *testing.T) {
	Help()
}

func Test_GetCommand(t *testing.T) {
	c := getCommand("add")
	if c == nil {
		t.Error("can not get command 'add'")
	}

	t.Log((*c).Name())
}

func Test_RegexMatch(t *testing.T) {
	reg, err := getRegex("y%e")
	if err != nil {
		t.Error(err)
	}

	strMap := map[string]bool{
		"yidane":   true,
		"yinsiwen": false,
		"nxin":     false,
		"yasde":    true,
		"yeyeyee":  true,
	}

	for k, v := range strMap {
		if reg.MatchString(k) != v {
			t.Error(fmt.Sprint(k, " match should be ", v))
		}
	}
}

func Test_Command(t *testing.T) {
	ridContext := context.NewRidContext()
	if err := ridContext.Login("yibihao", "yibihao"); err != nil {
		t.Error(err)
		return
	}

	UseCommand{}.Exec(ridContext, "nx_crm")
	OutCommand{}.Exec(ridContext, "E:\\rid")
	FindCommand{}.Exec(ridContext, "%a%")
	AddCommand{}.Exec(ridContext, "%a%")
	ListCommand{}.Exec(ridContext)
}
