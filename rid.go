package main

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"runtime"
	"sync"

	"github.com/yidane/rid/log"
)

//RidClient for http request with rid
type RidClient struct {
	Cookies        []*http.Cookie
	allDBs         map[string]*DBInfo
	CurrentDB      *DBInfo
	allTables      map[string][]string
	selectedTables map[string]int
	Output         string
}

//DBInfo 数据库对象信息
type DBInfo struct {
	ID   int    `json:"id"`
	IP   string `json:"ip"`
	Name string `json:"name"`
	Port string `json:"port"`
	Type int    `json:"type"`
}

type TableScriptInfo struct {
	DBName string
	Name   string
	Script string
	Tick   float64
	HTTP   bool
	Path   string
}

//AllDataBases 返回数据库列表
func (ridClient *RidClient) AllDataBases() []string {
	result := []string{}
	for _, db := range ridClient.allDBs {
		result = append(result, db.Name)
	}

	return result
}

//SelectedTables 返回数据库表集合
func (ridClient RidClient) SelectedTables() []string {
	result := []string{}
	for k := range ridClient.selectedTables {
		result = append(result, k)
	}
	return result
}

//SetCurrentDatabase choose database
func (ridClient *RidClient) SetCurrentDatabase(dbName string) bool {
	if ridClient.allDBs == nil {
		return false
	}

	if dbInfo, hasDb := ridClient.allDBs[dbName]; hasDb {
		ridClient.CurrentDB = dbInfo
		return true
	}

	return false
}

//DownloadAll for download table data in cache
func (ridClient *RidClient) DownloadAll() error {
	if ridClient.CurrentDB == nil {
		return errors.New("尚未选择数据库")
	}

	if len(ridClient.selectedTables) == 0 {
		return errors.New("尚未选择任何表")
	}

	if ridClient.Output == "" {
		return errors.New("尚未设置输出目录")
	}

	//创建临时文件夹
	var outDir = ridClient.Output + ridClient.CurrentDB.Name
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
	wg.Add(len(ridClient.selectedTables))

	for t := range ridClient.selectedTables {
		go func(tab string) {
			path := outDir + "/" + tab + ".sql"
			tabInfo := TableScriptInfo{
				DBName: ridClient.CurrentDB.Name,
				Name:   tab,
				Path:   path,
			}
			if !existsFile(path) {
				tinfo, err := ridClient.Download(ridClient.CurrentDB.Name, tab)
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
func (ridClient *RidClient) SetOutput(fold string) {
	_, err := os.Stat(fold)
	if err == nil {
		ridClient.Output = fold
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

func writeStringToFile(fileName, fileContent string) {
	f, err := os.Create(fileName)
	if err != nil {
		log.Error(err)
	}
	defer f.Close()
	f.WriteString(fileContent)
}
