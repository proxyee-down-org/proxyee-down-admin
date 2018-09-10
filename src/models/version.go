package models

import (
	"common"
	_ "github.com/go-sql-driver/mysql"
	"strconv"
	"strings"
	"time"
)

type Version struct {
	id          int64     `json:"id"`
	Version     float64   `json:"version"`
	Path        string    `json:"path"`
	BakPaths    []string  `json:"bakPaths"`
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
	defer rows.Close()
	if err != nil {
		return nil, err
	}
	if rows.Next() {
		var version Version
		err = rows.Scan(&version.Version, &version.Path, &version.Description)
		if err != nil {
			return nil, err
		} else {
			version.BakPaths = []string{common.Config.App.Server + version.Path}
			v := strconv.FormatFloat(version.Version, 'f', -1, 64)
			if strings.Index(v, ".") == -1 {
				v += ".0"
			}
			version.Path = "https://github.com/proxyee-down-org/proxyee-down/releases/download/" + v + "/proxyee-down-main.jar"
			return &version, nil
		}
	}
	return nil, nil
}
