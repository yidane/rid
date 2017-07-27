package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const loginURL string = "http://rid.nxin.com/Login"
const loadDataBaseURL string = "http://rid.nxin.com/index"
const loadTablesURL string = "http://rid.nxin.com/dbquery/getTables/"
const exportURL string = "http://rid.nxin.com/dbquery/exportTableData/"

func (ridClient *RidClient) newRequest(method, url string) (*http.Request, error) {
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Host", "rid.nxin.com")
	req.Header.Add("Connection", "keep-alive")
	req.Header.Add("Origin", "http://rid.nxin.com")
	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 6.3; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/59.0.3071.115 Safari/537.36") // nothing
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")                                                            // must
	req.Header.Add("Accept", "application/json, text/javascript, */*; q=0.01")
	req.Header.Add("X-Requested-With", "XMLHttpRequest")
	req.Header.Add("Referer", "http://rid.nxin.com/login.html")
	req.Header.Add("Accept-Encoding", "gzip, deflate")
	req.Header.Add("Accept-Language", "zh-CN,zh;q=0.8")

	if ridClient.Cookies != nil {
		for i := 0; i < len(ridClient.Cookies); i++ {
			req.AddCookie(ridClient.Cookies[i])
		}
	}

	return req, nil
}

func (ridClient *RidClient) getResponseContent(req *http.Request) ([]byte, error) {
	client := &http.Client{}
	client.Timeout = time.Hour
	response, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	//reload cookie
	newCookies := response.Cookies()
	if len(newCookies) > 0 {
		ridClient.Cookies = response.Cookies()
	}
	return data, nil
}

//Login login rid
func (ridClient *RidClient) Login(uid, pwd string) error {
	req, err := ridClient.newRequest("POST", loginURL)
	if err != nil {
		return err
	}

	body := strings.NewReader("userName=" + uid + "&password=" + pwd)
	req.Body = ioutil.NopCloser(body)

	data, err := ridClient.getResponseContent(req)
	if err != nil {
		fmt.Println(err)
		return err
	}

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

	return nil
}

//LoadDataBase load all database that current user can view
func (ridClient *RidClient) LoadDataBase() error {
	if ridClient.Cookies == nil {
		return errors.New("请先登录后在加载数据库信息")
	}

	req, err := ridClient.newRequest("GET", loadDataBaseURL)
	if err != nil {
		return err
	}

	data, err := ridClient.getResponseContent(req)
	if err != nil {
		return err
	}

	reg, _ := regexp.Compile(`\[{\"id.*\]`)
	if !reg.Match(data) {
		return errors.New("未找到相应数据库信息")
	}

	dbs := []*DBInfo{}
	err = json.Unmarshal(reg.Find(data), &dbs)
	if err != nil {
		return err
	}

	ridClient.allDBs = make(map[string]*DBInfo)
	for _, db := range dbs {
		ridClient.allDBs[db.Name] = db
	}

	return nil
}

//LoadTables load all tables of the database named dbName
func (ridClient *RidClient) LoadTables(dbName string) error {
	if !ridClient.SetCurrentDatabase(dbName) {
		return errors.New("不存在数据库[" + dbName + "]")
	}

	req, err := ridClient.newRequest("GET", loadTablesURL+strconv.Itoa(ridClient.CurrentDB.ID))
	if err != nil {
		return err
	}

	data, err := ridClient.getResponseContent(req)
	if err != nil {
		return err
	}

	reg, _ := regexp.Compile(`\[.*\]`)
	if !reg.Match(data) {
		return errors.New("查询返回数据结果集为空")
	}

	type tableInfo struct {
		DbName string `json:"dbName"`
	}
	tables := []tableInfo{}
	err = json.Unmarshal(reg.Find(data), &tables)
	if err != nil {
		return err
	}

	ridClient.selectedTables = make(map[string]int)
	if len(tables) > 0 {
		ts := []string{}
		for _, t := range tables {
			ts = append(ts, t.DbName)
		}
		ridClient.allTables[dbName] = ts
	}

	return nil
}

//Download download biz data
func (ridClient RidClient) Download(dbName, tableName string) (*TableScriptInfo, error) {
	now := time.Now()
	if !ridClient.SetCurrentDatabase(dbName) {
		return nil, errors.New("不存在该数据库")
	}

	req, err := ridClient.newRequest("GET", exportURL+strconv.Itoa(ridClient.CurrentDB.ID)+"/"+tableName)
	if err != nil {
		return nil, err
	}

	data, err := ridClient.getResponseContent(req)
	if err != nil {
		return nil, err
	}

	return &TableScriptInfo{
		DBName: dbName,
		Name:   tableName,
		Script: string(data),
		Tick:   time.Now().Sub(now).Seconds(),
		HTTP:   true,
	}, nil
}
