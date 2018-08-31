package models

import (
	"common"
	_ "github.com/go-sql-driver/mysql"
	"time"
)

type Version struct {
	id          int64     `json:"id"`
	Version     float64   `json:"version"`
	Path        string    `json:"path"`
	Description string    `json:"description"`
	createTime  time.Time `json:"createTime"`
}

func GetNewVersion() (*Version, error) {
	db, err := common.GetDb()
	defer db.Close()
	if err != nil {
		return nil, err
	}
	rows, err := db.Query("select version,path,description from version order by version desc limit 1")
	if err != nil {
		return nil, err
	}
	if rows.Next() {
		var version Version
		err = rows.Scan(&version.Version, &version.Path, &version.Description)
		if err != nil {
			return nil, err
		} else {
			version.Path = common.Config.App.Server + version.Path
			return &version, nil
		}
	}
	return nil, nil
}
