package command

import "testing"

func Test_Help(t *testing.T) {
	Help()
}

func Test_GetCommand(t *testing.T){
	c:=getCommand("add")
	if c==nil{
		t.Error("can not get command 'add'")
	}

	t.Log((*c).Name())
}