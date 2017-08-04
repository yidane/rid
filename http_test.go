package main

import (
	"fmt"
	"testing"
)

func initRidClient() error {
	uid := "yibihao"
	pwd := "yibihao"
	ridClient = &RidClient{}
	return ridClient.Login(uid, pwd)
}

func Test_Login(t *testing.T) {
	err := initRidClient()
	if err != nil {
		t.Error(err)
	}

	fmt.Println(ridClient.Cookies)
}

func Test_LoadDataBase(t *testing.T) {
	Test_Login(t)

	err := ridClient.LoadDataBase()
	if err != nil {
		fmt.Println(err.Error())
	}

	if len(ridClient.AllDataBases()) == 0 {
		t.Fatal("获取的数据库集合为空")
	}

	fmt.Println("有权限访问的数据库集合为：\n", ridClient.AllDataBases())
}

func Test_LoadTables(t *testing.T) {
	Test_Login(t)
	Test_LoadDataBase(t)
	err := ridClient.LoadTables("nx_crm")
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println("已加载的所有表集合\n:", ridClient.allTables)
}

// func Test_Login(t *testing.T) {
// 	uid := "yibihao"
// 	pwd := "yibihao"
// 	ridClient := &RidClient{}
// 	err := ridClient.Login(uid, pwd)
// 	if err != nil {
// 		fmt.Println(err.Error())
// 	}

// 	err = ridClient.LoadTables("nx_crm")
// 	if err != nil {
// 		fmt.Println(err.Error())
// 	}
// 	fmt.Println(ridClient.TablesString())

// 	_, err = ridClient.Download("nx_crm", "role_user")
// 	if err != nil {
// 		fmt.Println(err.Error())
// 	}
// }
