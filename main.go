package main

import (
	// "bytes"
	"container/list"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"gopkg.in/gorp.v1"
	"log"
	"os"
	"os/exec"
	"sync"
	"text/template"
)

var (
	Cfg   YbParams
	Dba   *gorp.DbMap
	GLock sync.Mutex
)

func init() {
	//read config
	cfg, err := InitCfg("config.json")
	if err != nil {
		fmt.Println("init config failed! ", err)
		return
	}
	Cfg = cfg

	//get db instance
	db, err := InitDb()
	if err != nil {
		fmt.Println("init db failed! ", err)
		return
	}
	Dba = db
}

func InitDb() (dbMap *gorp.DbMap, err error) {
	db, err := sql.Open("mysql", Cfg.ConnStr())
	if err != nil {
		fmt.Println("open sql failed")
		return nil, err
	}
	dbMap = new(gorp.DbMap)
	dbMap.Db = db
	dbMap.Dialect = gorp.MySQLDialect{}
	err = dbMap.Db.Ping()
	return dbMap, err
}

func genSource(table, templateSource string, args map[string]interface{}) (err error) {
	GLock.Lock()
	defer GLock.Unlock()

	//new template parser
	tmpl, err := template.New(table).Parse(templateSource)
	if err != nil {
		return err
	}

	filePath := Cfg.OutDir + "/" + table + ".go"
	f, err := os.OpenFile(filePath, os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0766)
	defer f.Close()
	if err != nil {
		return err
	}

	//wirte to file
	err = tmpl.Execute(f, args)
	tmpl = nil

	//format code
	cmd := exec.Command("gofmt", "-s", "-w", filePath)
	err = cmd.Run()
	return
}

func ProduceModel() {

	task := make(chan int, 4)
	done := make(chan bool)
	failed := list.New()

	count := len(Cfg.Tables)
	if count == 0 {
		fmt.Println("not found table to procee")
		return
	}

	if Cfg.OutDir == "." || Cfg.OutDir == "./" {
		err := os.Mkdir("Models", os.ModePerm)
		if err != nil && !os.IsExist(err) {
			fmt.Println("create dir failed ", err)
			return
		}
	} else {
		err := os.MkdirAll(Cfg.OutDir, os.ModePerm)
		if err != nil && !os.IsExist(err) {
			fmt.Println("create dir failed ", err)
			return
		}
	}
	Cfg.OutDir += "/Models"

	ToFile := func(table string) {
		fmt.Printf("to process %s table\n", table)

		sql := `select COLUMN_NAME,DATA_TYPE from COLUMNS where table_name = '%s' and table_schema ='%s';`
		sql = fmt.Sprintf(sql, table, Cfg.Schema)
		var columns []TableStru
		_, err := Dba.Select(&columns, sql)
		if err != nil {
			fmt.Printf("scan table column from database failed for : %s \n", table)
			failed.PushBack(table)
		}

		if len(columns) == 0 {
			failed.PushBack(table)
		}
		fmt.Println("cloums is ", columns)
		templateArgs := map[string]interface{}{
			"Table":  table,
			"Colums": columns,
			"Time":   false,
		}
		err = genSource(table, MODEL, templateArgs)
		if err != nil {
			panic(err)
		}

		//a task successful
		// GLock.Lock()
		currTask := <-task
		// GLock.Unlock()
		if currTask == count-1 {
			done <- true
		}
	}

	fmt.Println("count is ", count)

	for index, table := range Cfg.Tables {
		fmt.Println("begin to do jobs")
		//assign task
		task <- index
		go ToFile(table)
	}

	//wait end
	<-done

	fmt.Printf("%d table failed\n", failed.Len())
	if failed.Len() > 0 {
		fmt.Println("the failed table list follown:")
		table := 1
		for ele := failed.Front(); ele != nil; ele = ele.Next() {
			if table%5 == 0 {
				fmt.Println()
			} else {
				fmt.Printf("%s\t", ele.Value.(string))
			}
			table++
		}
	}
	fmt.Printf("\ncongratulataions to see file in %s dir\n", Cfg.OutDir)
}

func main() {
	// buf := new(bytes.Buffer)
	Dba.TraceOn("", log.New(os.Stdout, "model", log.Lshortfile|log.Ltime))
	ProduceModel()
}
