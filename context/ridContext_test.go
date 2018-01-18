package context

import "testing"

var testUserInfo = UserInfo{UserID: `zcm99`, Password: `dbn002385`}

func init() {

}

func Test_LoginFaild(t *testing.T) {
	ridClient := &RidContext{}
	err := ridClient.Login(&testUserInfo)
	if err == nil {
		t.Error("there should have error message if login faild")
	}
	t.Log(err)
	if ridClient.HttpContext.HasLogin {
		t.Error("login status should be false")
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
