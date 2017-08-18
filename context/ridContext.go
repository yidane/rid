package context

import (
	"errors"
	"fmt"
	"os"
	"runtime"
	"sync"

	"github.com/yidane/rid/log"
)

type RidContext struct {
	HttpContext    *HttpContext
	allDBs         map[string]*DBInfo
	CurrentDB      *DBInfo
	allTables      map[string][]string
	selectedTables map[string]bool
	Output         string
}

func NewRidContext() *RidContext {
	return &RidContext{
		allDBs:         make(map[string]*DBInfo),
		CurrentDB:      nil,
		allTables:      make(map[string][]string),
		selectedTables: make(map[string]bool),
		Output:         "",
	}
}

//DBInfo 数据库对象信息
type DBInfo struct {
	ID   int    `json:"id"`
	IP   string `json:"ip"`
	Name string `json:"name"`
	Port string `json:"port"`
	Type int    `json:"type"`
}

//TableScriptInfo 数据表对象信息
type TableScriptInfo struct {
	Name   string
	Script string
	Tick   float64
	HTTP   bool
	Path   string
	Error  *error
}

func (tableScriptInfo TableScriptInfo) Save() error {
	if tableScriptInfo.HTTP {
		f, err := os.Create(tableScriptInfo.Path)
		if err != nil {
			log.Error(err)
		}
		defer f.Close()
		f.WriteString(tableScriptInfo.Script)
	}
	log.Succeed(tableScriptInfo.Name+"	%vs\r\n", tableScriptInfo.Tick)
	return nil
}

func (tableScriptInfo TableScriptInfo) Exists() bool {
	_, err := os.Stat(tableScriptInfo.Path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}

func (ridContext *RidContext) createTableScriptInfo(name string) *TableScriptInfo {
	tableScriptInfo := new(TableScriptInfo)
	tableScriptInfo.Path = ridContext.Output + "\\" + name + ".sql"
	tableScriptInfo.Name = name

	return tableScriptInfo
}

//AllDataBases 返回数据库列表
func (ridContext *RidContext) AllDataBases() []string {
	result := []string{}
	for _, db := range ridContext.allDBs {
		result = append(result, db.Name)
	}

	return result
}

//SelectedTables 返回数据库表集合
func (ridContext *RidContext) SelectedTables() []string {
	result := []string{}
	for k, v := range ridContext.selectedTables {
		if v {
			result = append(result, k)
		}
	}
	return result
}

//SetCurrentDatabase choose database
func (ridContext *RidContext) SetCurrentDatabase(dbName string) error {
	if ridContext.allDBs == nil || len(ridContext.allDBs) == 0 {
		_, err := ridContext.LoadDataBase()
		if err != nil {
			return err
		}
	}

	if len(ridContext.allDBs) == 0 {
		return errors.New("您有权限访问的数据库列表为空")
	}

	if dbInfo, hasDb := ridContext.allDBs[dbName]; hasDb {
		ridContext.CurrentDB = dbInfo
		ridContext.selectedTables = make(map[string]bool)
		tables, err := ridContext.LoadTables(dbName)
		if err != nil {
			return nil
		}
		for i := 0; i < len(tables); i++ {
			ridContext.selectedTables[tables[i]] = false
		}

		return nil
	}

	return errors.New(fmt.Sprint("无此数据库", dbName))
}

//DownloadAll for download table data in cache
func (ridContext *RidContext) DownloadAll() error {
	if ridContext.CurrentDB == nil {
		return errors.New("尚未选择数据库")
	}
	if len(ridContext.selectedTables) == 0 {
		return errors.New("尚未选择任何表")
	}
	if ridContext.Output == "" {
		return errors.New("尚未设置输出目录")
	}

	total := 0
	s := 0
	f := 0

	for _, v := range ridContext.selectedTables {
		if v {
			total++
		}
	}

	runtime.GOMAXPROCS(runtime.NumCPU() * 5)
	ch := make(chan *TableScriptInfo, 10)
	go func() {
		for {
			select {
			case table := <-ch:
				//判断保存是否成功
				table.Save()
			default:
			}
		}
	}()

	wg := sync.WaitGroup{}
	wg.Add(total)
	for k, v := range ridContext.selectedTables {
		if v {
			go func(tab string) {
				tabInfo := ridContext.createTableScriptInfo(tab)
				if !tabInfo.Exists() {
					t, err := ridContext.HttpContext.Download(ridContext.CurrentDB.ID, ridContext.CurrentDB.Name, tab)
					if err != nil {
						tabInfo.Error = &err
						f++
					} else {
						s++
						tabInfo.HTTP = t.HTTP
						tabInfo.Script = t.Script
						tabInfo.Tick = t.Tick
					}
				}
				ch <- tabInfo
				wg.Done()
			}(k)
		}
	}
	wg.Wait()

	log.Succeed("Total:", total, "	Success:", s, "	Failed:", f)

	return nil
}

//SetOutput set the output fold
func (ridContext *RidContext) SetOutput(fold string) error {
	if ridContext.CurrentDB == nil {
		return errors.New("no database selected")
	}

	foldPath := fold + "\\" + ridContext.CurrentDB.Name
	_, err := os.Stat(foldPath)
	if err == nil {
		ridContext.Output = foldPath
		return nil
	}

	if os.IsNotExist(err) {
		err := os.MkdirAll(foldPath, 0777)
		if err != nil {
			return err
		}
		ridContext.Output = foldPath
		return nil
	}

	return err
}

func (ridContext *RidContext) LoadDataBase() ([]string, error) {
	if ridContext.HttpContext == nil || !ridContext.HttpContext.HasLogin {
		return nil, errors.New("login must be first")
	}
	dbArr := []string{}
	if len(ridContext.allDBs) > 0 {
		for _, db := range ridContext.allDBs {
			dbArr = append(dbArr, db.Name)
		}
		return dbArr, nil
	}

	dbs, err := ridContext.HttpContext.LoadDataBase()
	if err != nil {
		return nil, err
	}
	for i := 0; i < len(dbs); i++ {
		if _, ok := ridContext.allDBs[dbs[i].Name]; !ok {
			ridContext.allDBs[dbs[i].Name] = dbs[i]
			dbArr = append(dbArr, dbs[i].Name)
		}
	}

	return dbArr, nil
}

func (ridContext *RidContext) LoadTables(dbName string) ([]string, error) {
	if ridContext.HttpContext == nil || !ridContext.HttpContext.HasLogin {
		return nil, errors.New("login must be first")
	}
	if tables, ok := ridContext.allTables[dbName]; ok {
		return tables, nil
	}

	if dbInfo, hasDb := ridContext.allDBs[dbName]; hasDb {
		tables, err := ridContext.HttpContext.LoadTables(dbInfo.ID, dbInfo.Name)
		if err != nil {
			return nil, err
		}

		ridContext.allTables[dbName] = tables
		return tables, nil
	}

	return nil, errors.New(fmt.Sprint("no such database ", dbName))
}

func (ridContext *RidContext) Login(uid, pwd string) error {
	if len(uid) == 0 {
		return errors.New("uid can not be nil")
	}
	if len(pwd) == 0 {
		return errors.New("pwd can not be nil")
	}

	ridContext.HttpContext = new(HttpContext)
	return ridContext.HttpContext.Login(uid, pwd)
}

func (ridContext *RidContext) AddToCache(name ...string) {
	for _, s := range name {
		if _, ok := ridContext.selectedTables[s]; ok {
			ridContext.selectedTables[s] = true
		}
	}
}

func (ridContext *RidContext) RemoveFromCache(name ...string) {
	for _, s := range name {
		if _, ok := ridContext.selectedTables[s]; ok {
			ridContext.selectedTables[s] = false
		}
	}
}

func (ridContext *RidContext) ClearCache() {
	for k := range ridContext.selectedTables {
		ridContext.selectedTables[k] = false
	}
}
