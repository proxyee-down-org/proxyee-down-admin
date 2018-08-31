package models

import (
	"common"
	_ "github.com/go-sql-driver/mysql"
	"time"
)

type ExtensionDown struct {
	Id         int64
	ExtId      int64
	Version    float64
	Ip         string
	Ua         string
	CreateTime time.Time
}

func (extensionDown *ExtensionDown) Insert() {
	db, err := common.GetDb()
	defer db.Close()
	if err != nil {
		return
	}
	stmt, err := db.Prepare("insert into extension_down (ext_id,version,ip,ua,create_time) values (?,?,?,?,?)")
	defer stmt.Close()
	if err != nil {
		return
	}
	stmt.Exec(extensionDown.ExtId, extensionDown.Version, extensionDown.Ip, extensionDown.Ua, extensionDown.CreateTime)
}
