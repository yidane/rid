package command

import (
	"bytes"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"sync"

	"github.com/yidane/rid/context"
	"github.com/yidane/rid/log"
)

type Command interface {
	Name() string
	Exec(ridContext *context.RidContext, args ...string)
	Usage() string
}

var cMap = make(map[string]*Command)
var once sync.Once = sync.Once{}

func init() {
	once.Do(func() {
		packageCommand(AddCommand{})
		packageCommand(ClearCommand{})
		packageCommand(DownloadCommand{})
		packageCommand(FindCommand{})
		packageCommand(ListCommand{})
		packageCommand(LoadCommand{})
		packageCommand(OutCommand{})
		packageCommand(RemoveCommand{})
		packageCommand(UseCommand{})
		packageCommand(ShowCommand{})
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

func Help(commandName ...string) {
	buf := bytes.Buffer{}
	errCmd := []string{}
	if len(commandName) > 0 {
		for _, c := range commandName {
			if cmd, ok := cMap[c]; ok {
				buf.WriteString(fmt.Sprintf(`	|%-20s|%s%s`, c, (*cmd).Usage(), "\n"))
			} else {
				errCmd = append(errCmd, c)
			}
		}
	} else {
		for _, cmd := range cMap {
			buf.WriteString(fmt.Sprintf(`	|%-20s|%s%s`, (*cmd).Name(), (*cmd).Usage(), "\n"))
		}
	}

	if buf.Len() > 0 {
		log.Succeed("command usage as follows：")
		fmt.Printf(buf.String())
	}
	if len(errCmd) > 0 {
		log.Error(fmt.Sprintf("no such commands:%s", strings.Join(errCmd, ",")))
	}
}

//采用%作为通配符
//%ne 表示以ne结尾的字符串
//yi% 表示以yi开始的字符串
//yi%ne表示以yi开始ne结尾的字符串
//y%d%ne表示以y开始中间包含d且以ne结尾的字符串
func getRegex(reg string) (*regexp.Regexp, error) {
	reg = strings.Trim(reg, "")
	l := len(reg)
	if l == 0 {
		return nil, errors.New("need argument")
	}

	newReg := ""
	for i := 0; i < l; i++ {
		c := reg[i]
		if c == '%' || c == '*' {
			newReg += "(.*)"
			continue
		}
		switch {
		case i == 0:
			newReg += "^" + string(c)
		case l-1 == i:
			newReg += string(c) + "$"
		case i > 0 && i < l-1:
			newReg += string(c)
		}
	}

	return regexp.Compile(newReg)
}
