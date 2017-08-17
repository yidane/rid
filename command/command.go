package command

import (
	"github.com/yidane/rid/context"
	"sync"
	"fmt"
	"github.com/yidane/rid/log"
	"strings"
	"regexp"
)

type Command interface {
	Name() string
	Exec(ricContext *context.RidContext, args ...string)
	Usage() string
}

var cMap = make(map[string]*Command)
var once sync.Once = sync.Once{}

func init() {
	once.Do(func() {
		packageCommand(AddCommand{})
		packageCommand(ClearCommand{})
		packageCommand(DownloadCommand{})
		packageCommand(ListCommand{})
		packageCommand(LoadCommand{})
		packageCommand(OutCommand{})
		packageCommand(RemoveCommand{})
		packageCommand(UseCommand{})
	})
}

func packageCommand(command Command) {
	if command == nil {
		panic("argument command can not be nil")
	}

	if _, ok := cMap[command.Name()]; ok {
		panic(fmt.Sprint("Command ", command.Name(), " has loaded"))
	}

	cMap[command.Name()] = &command
}

func getCommand(cName string) *Command {
	if c, ok := cMap[cName]; ok {
		return c
	}
	return nil
}

func Exec(ridContext *context.RidContext, cName string, args ...string) {
	c := getCommand(cName)
	if c == nil {
		log.Error("no such command like ", cName)
		return
	}
	(*c).Exec(ridContext, args...)
}

func Help(commandName ...[]string) {
	//output all usage
}

//采用%作为通配符
//%ne 表示以ne结尾的匹配
//yi% 表示以yi
//func isMatch(reg, str string) bool {
//	reg = strings.Trim(reg, "")
//	l := len(reg)
//	if l == 0 {
//		return false
//	}
//
//	newReg := ""
//	switch l {
//	case 1:
//		if reg[0] == '%' {
//			return false
//		}
//	case 2:
//		if reg[0] == '%' {
//			newReg += "(.*)"
//		} else {
//			newReg += string(reg[0])
//		}
//		if reg[1] == '%' {
//			newReg += "(.*)"
//		} else {
//			newReg += string(reg[0])
//		}
//	case l > 2:
//		for i := 1; i < len(reg)-1; i++ {
//			c := reg[i]
//			if c != '%' {
//				newReg += string(c)
//				continue
//			}
//
//			if i == 0 {
//
//			} else if i == len(reg) {
//
//			} else {
//
//			}
//		}
//	}
//
//	regexp,err:= regexp.Compile(newReg)
//	if err!=nil{
//		return false
//	}
//}
