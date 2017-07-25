package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
)

const loginURL string = "http://rid.nxin.com/Login"
const indexURL string = "http://rid.nxin.com/index"
const databaseURL string = ""
const exportURL string = ""

//RidClient for http request with rid
type RidClient struct {
	Cookies []*http.Cookie
	DBList  []*DBInfo
}

//DBInfo 数据库对象信息
type DBInfo struct {
	ID   int    `json:"id"`
	IP   string `json:"ip"`
	Name string `json:"name"`
	Port string `json:"port"`
	Type int    `json:"type"`
}

//DBString 返回数据库列表
func (ridClient *RidClient) DBString() string {
	s, err := json.Marshal(ridClient.DBList)
	if err != nil {
		return err.Error()
	}

	return string(s)
}

//Login login rid
func (ridClient *RidClient) Login(uid, pwd string) error {
	body := strings.NewReader("userName=" + uid + "&password=" + pwd)

	client := &http.Client{}
	req, err := http.NewRequest("POST", loginURL, body)
	if err != nil {
		return err
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

	response, err := client.Do(req)
	if err != nil {
		return err
	}

	defer response.Body.Close()

	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
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

	ridClient.Cookies = response.Cookies()

	return nil
}

//LoadDataBase load all database that current user can view
func (ridClient *RidClient) LoadDataBase() error {
	client := &http.Client{}
	req, err := http.NewRequest("GET", indexURL, nil)

	if err != nil {
		return err
	}

	if ridClient.Cookies != nil {
		for i := 0; i < len(ridClient.Cookies); i++ {
			req.AddCookie(ridClient.Cookies[i])
		}
	}

	response, err := client.Do(req)
	if err != nil {
		return err
	}

	defer response.Body.Close()

	bytesData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}

	reg, err := regexp.Compile(`\[{\"id.*\]`)
	if err != nil {
		return err
	}

	if !reg.Match(bytesData) {
		return errors.New("未找到相应数据库信息")
	}

	ridClient.DBList = []*DBInfo{}
	err = json.Unmarshal(reg.Find(bytesData), &ridClient.DBList)
	if err != nil {
		return err
	}

	return nil
}

//LoadTables load all tables of database dbName
func (ridClient *RidClient) LoadTables(dbName string) error {
	if ridClient.DBList == nil || len(ridClient.DBList) == 0 {
		return errors.New("当前用户尚未初始化或无数据库权限")
	}

	hasDb := false
	for _, dbInfo := range ridClient.DBList {
		if dbInfo.Name == dbName {
			hasDb = true
			break
		}
	}

	if !hasDb {
		return errors.New("不存在数据库[" + dbName + "]")
	}

	return errors.New("")
}
