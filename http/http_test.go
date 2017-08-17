package http

import (
	"fmt"
	"testing"
)


func Test_GetResponse(t *testing.T){
	url:="http://rid.nxin.com/login.html"
	request,err:=NewRequest("GET",url,nil)
	if err!=nil{
		t.Error(err)
	}
	data,cookies, err:= GetResponseContent(request)
	if err!=nil{
		t.Error(err)
	}
	if cookies!=nil&&len(cookies)>0{
		for i:=0;i<len(cookies) ;i++  {
			fmt.Println(cookies[i].Name,":",cookies[i].Value)
		}
	}

	if data==nil{
		t.Error("empty response")
	}

	fmt.Println(string(data))
}
