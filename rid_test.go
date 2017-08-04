package main

import "testing"
import "fmt"

func Test_NewLogin(t *testing.T) {
	uid := "yibihao"
	pwd := "yibihao"
	ridClient := &RidClient{}
	err := ridClient.Login(uid, pwd)
	if err != nil {
		fmt.Println(err.Error())
	}
}

// func Test_DownloadAll(t *testing.T) {
// 	uid := "yibihao"
// 	pwd := "yibihao"
// 	ridClient := &RidClient{}
// 	err := ridClient.Login(uid, pwd)
// 	if err != nil {
// 		fmt.Println(err.Error())
// 	}

// 	ridClient.DownloadAll("nx_crm")
// }
