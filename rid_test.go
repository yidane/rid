package main

import "testing"
import "fmt"

func Test_Login(t *testing.T) {
	uid := "yibihao"
	pwd := "yibihao"
	ridClient := &RidClient{}
	err := ridClient.Login(uid, pwd)
	if err != nil {
		fmt.Println(err.Error())
	}

	err = ridClient.LoadDataBase()
	if err != nil {
		fmt.Println(err.Error())
	}

	fmt.Println(ridClient.DBString())
}
