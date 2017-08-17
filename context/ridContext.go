package context

import (
	"fmt"
	"os"
	"errors"
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
	DBName string
	Name   string
	Script string
	Tick   float64
	HTTP   bool
	Path   string
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

	//创建临时文件夹
	var outDir = ridContext.Output + ridContext.CurrentDB.Name
	err := os.MkdirAll(outDir, 0777)
	if err != nil {
		return err
	}

	var existsFile = func(path string) bool {
		_, err := os.Stat(path)
		if err == nil {
			return true
		}
		if os.IsNotExist(err) {
			return false
		}
		return false
	}

	runtime.GOMAXPROCS(runtime.NumCPU())
	ch := make(chan *TableScriptInfo)
	go func() {
		for {
			select {
			case table := <-ch:
				if table.HTTP {
					writeStringToFile(table.Path, table.Script)
				}
				log.Succeed(table.DBName+"	"+table.Name+"	%vs\r\n", table.Tick)
			default:
			}
		}
	}()

	wg := sync.WaitGroup{}
	wg.Add(len(ridContext.selectedTables))

	for t := range ridContext.selectedTables {
		go func(tab string) {
			path := outDir + "/" + tab + ".sql"
			tabInfo := TableScriptInfo{
				DBName: ridContext.CurrentDB.Name,
				Name:   tab,
				Path:   path,
			}
			if !existsFile(path) {
				tinfo, err := ridContext.HttpContext.Download(ridContext.CurrentDB.ID, ridContext.CurrentDB.Name, tab)
				if err != nil {
					fmt.Println(err)
				}

				tabInfo.HTTP = tinfo.HTTP
				tabInfo.Script = tinfo.Script
				tabInfo.Tick = tinfo.Tick
			}

			ch <- &tabInfo
			wg.Done()
		}(t)
	}

	wg.Wait()

	fmt.Println("finish")
	return nil
}

//SetOutput set the output fold
func (ridContext *RidContext) SetOutput(fold string) {
	_, err := os.Stat(fold)
	if err == nil {
		ridContext.Output = fold
		return
	}

	if os.IsNotExist(err) {
		err := os.MkdirAll(fold, 0777)
		if err != nil {
			log.Error(err)
		}
		return
	}
	log.Error(err)
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

func writeStringToFile(fileName, fileContent string) {
	f, err := os.Create(fileName)
	if err != nil {
		log.Error(err)
	}
	defer f.Close()
	f.WriteString(fileContent)
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
	for k, _ := range ridContext.selectedTables {
		ridContext.selectedTables[k] = false
	}
}
