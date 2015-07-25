package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type YbParams struct {
	DbHost   string   `json:"dbhost"`
	DbPort   string   `json:"dbport"`
	User     string   `json:"user"`
	Pass     string   `json:"pass"`
	Database string   `json:"database"`
	Schema   string   `json:"schema"`
	Tables   []string `json:"table"`
	OutDir   string   `json:"outdir"`
}

func InitCfg(file string) (para YbParams, err error) {

	f, err := os.Open(file)
	raw, err := ioutil.ReadAll(f)

	err = json.Unmarshal(raw, &para)

	return
}

func (para YbParams) ConnStr() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", para.User, para.Pass, para.DbHost, para.DbPort, para.Database)
}

type TableStru struct {
	ColumName string `db:"COLUMN_NAME"`
	ColumType string `db:"DATA_TYPE"`
}

const MODEL = `// GENERATED CODE - DO NOT EDIT
package Model


{{if .Time}}
import (
	"time"
)
{{end}}

type {{.Table}} struct { {{range $i, $colum := .Colums}}
	{{$colum.ColumName}}  {{$colum.ColumType}}` + "    `" + `db:"{{$colum.ColumName}}"    json:"{{$colum.ColumName}}"` + "`" +
	`{{end}}
}
`
