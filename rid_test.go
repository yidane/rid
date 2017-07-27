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

	err = ridClient.LoadTables("nx_crm")
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(ridClient.TablesString())

	_, err = ridClient.Download("nx_crm", "role_user")
	if err != nil {
		fmt.Println(err.Error())
	}
}

func Test_NewLogin(t *testing.T) {
	uid := "yibihao"
	pwd := "yibihao"
	ridClient := &RidClient{}
	err := ridClient.Login(uid, pwd)
	if err != nil {
		fmt.Println(err.Error())
	}
}
func Test_DownloadAll(t *testing.T) {
	uid := "yibihao"
	pwd := "yibihao"
	ridClient := &RidClient{}
	err := ridClient.Login(uid, pwd)
	if err != nil {
		fmt.Println(err.Error())
	}

	ridClient.DownloadAll("nx_crm")
}
