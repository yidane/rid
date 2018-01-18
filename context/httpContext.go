package context

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"regexp"
	"strconv"
	"strings"
	"time"

	"net/http"

	ridHttp "github.com/yidane/rid/http"
)

const loginURL string = "http://rid.nxin.com/Login"
const loadDataBaseURL string = "http://rid.nxin.com/indexOld"
const loadTablesURL string = "http://rid.nxin.com/dbquery/getTables/"
const exportURL string = "http://rid.nxin.com/dbquery/exportTableData/"

//UserInfo 用户信息
type UserInfo struct {
	UserID   string
	Password string
}

//HttpContext http上下文
type HttpContext struct {
	CurrentUser *UserInfo
	Cookies     []*http.Cookie
	HasLogin    bool
}

//Login login rid
func (httpContext *HttpContext) Login() error {
	req, err := ridHttp.NewRequest("POST", loginURL, nil)
	if err != nil {
		return err
	}

	body := strings.NewReader("userName=" + httpContext.CurrentUser.UserID + "&password=" + httpContext.CurrentUser.Password)
	req.Body = ioutil.NopCloser(body)

	data, cookies, err := ridHttp.GetResponseContent(req)
	if err != nil {
		fmt.Println(err)
		return err
	}
	httpContext.Cookies = cookies

	type result struct {
		Code  int    `json:"code"`
		Data  int    `json:"data"`
		Error string `json:"error"`
	}

	r := &result{}
	err = json.Unmarshal(data, r)
	if err != nil {
		fmt.Println(string(data))
		return err
	}

	if r.Code != 0 {
		return errors.New(r.Error)
	}

	httpContext.HasLogin = true
	return nil
}

//LoadDataBase load all database that current user can view
func (httpContext *HttpContext) LoadDataBase() ([]*DBInfo, error) {
	if httpContext.Cookies == nil {
		return nil, errors.New("请先登录后在加载数据库信息")
	}

	req, err := ridHttp.NewRequest("GET", loadDataBaseURL, httpContext.Cookies)
	if err != nil {
		return nil, err
	}

	data, _, err := ridHttp.GetResponseContent(req)
	if err != nil {
		return nil, err
	}

	reg, _ := regexp.Compile(`\[{\"id.*\]`)
	if !reg.Match(data) {
		return nil, fmt.Errorf("未找到相应数据库信息:%s", string(data))
	}

	dbs := []*DBInfo{}
	err = json.Unmarshal(reg.Find(data), &dbs)
	if err != nil {
		return nil, err
	}

	return dbs, nil
}

//LoadTables load all tables of the database named dbName
func (httpContext *HttpContext) LoadTables(id int, dbName string) ([]string, error) {
	if httpContext.Cookies == nil {
		return nil, errors.New("请先登录后在加载数据库信息")
	}

	req, err := ridHttp.NewRequest("GET", loadTablesURL+strconv.Itoa(id), httpContext.Cookies)
	if err != nil {
		return nil, err
	}

	data, _, err := ridHttp.GetResponseContent(req)
	if err != nil {
		return nil, err
	}

	reg, _ := regexp.Compile(`\[.*\]`)
	if !reg.Match(data) {
		return nil, errors.New("查询返回数据结果集为空")
	}

	type tableInfo struct {
		DbName string `json:"dbName"`
	}

	tables := []tableInfo{}
	err = json.Unmarshal(reg.Find(data), &tables)
	if err != nil {
		return nil, err
	}

	if len(tables) > 0 {
		ts := []string{}
		for _, t := range tables {
			ts = append(ts, t.DbName)
		}

		return ts, nil
	}

	return nil, errors.New(fmt.Sprint("no table belong to database ", dbName))
}

//Download download biz data
func (httpContext *HttpContext) Download(id int, dbName, tableName string) (*TableScriptInfo, error) {
	now := time.Now()

	req, err := ridHttp.NewRequest("GET", exportURL+strconv.Itoa(id)+"/"+tableName, httpContext.Cookies)
	if err != nil {
		return nil, err
	}

	data, _, err := ridHttp.GetResponseContent(req)
	if err != nil {
		return nil, err
	}

	return &TableScriptInfo{
		Name:   tableName,
		Script: string(data),
		Tick:   time.Now().Sub(now).Seconds(),
		HTTP:   true,
	}, nil
}
