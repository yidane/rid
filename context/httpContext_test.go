package context

import "testing"

func Test_HttpLoginFail(t *testing.T){
	httpContext:=HttpContext{}
	err:= httpContext.Login("yibihao1","yibihao")

	if err==nil{
		t.Error("there should have error message when login faild")
	}
	t.Log(err)
}

func Test_HttpLoginSuccess(t *testing.T){
	httpContext:=HttpContext{}
	err:= httpContext.Login("yibihao","yibihao")

	if err!=nil{
		t.Error(err)
	}
	t.Log("login status is ",httpContext.HasLogin)
}

func Test_LoadDatabase(t *testing.T){
	httpContext:=HttpContext{}
	err:= httpContext.Login("yibihao","yibihao")

	if err!=nil{
		t.Error(err)
	}

	dblist,err:= httpContext.LoadDataBase()
	if err!=nil{
		t.Error(err)
	}
	if dblist==nil{
		t.Log("loaded nothing")
	}

	for i:=0;i<len(dblist) ;i++  {
		t.Log(*dblist[i])
	}
}

func Test_LoadTables(t *testing.T){
	httpContext:=HttpContext{}
	err:= httpContext.Login("yibihao","yibihao")

	if err!=nil{
		t.Error(err)
	}

	tables,err:= httpContext.LoadTables(184,"nx_crm")
	if err!=nil{
		t.Error(err)
	}
	if tables==nil{
		t.Error("loaded nothing")
	}

	for i:=0;i<len(tables) ;i++  {
		t.Log(tables[i])
	}
}