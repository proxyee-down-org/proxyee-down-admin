package models

import (
	"common"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"math"
	"strings"
	"time"
)

type Extension struct {
	Id          int64     `json:"id"`
	Title       string    `json:"title"`
	Version     float64   `json:"version"`
	Description string    `json:"description"`
	Path        string    `json:"path"`
	Files       string    `json:"files"`
	CreateTime  time.Time `json:"createTime"`
	UpdateTime  time.Time `json:"updateTime"`
}

func SelectExtensionByPath(path string) (*Extension, error) {
	db, err := common.GetDb()
	defer db.Close()
	if err != nil {
		return nil, err
	}
	rows, err := db.Query("select id,version from extension where path = ?", path)
	if err != nil {
		return nil, err
	}
	if rows.Next() {
		var extension Extension
		err = rows.Scan(&extension.Id, &extension.Version)
		if err != nil {
			return nil, err
		} else {
			return &extension, nil
		}
	}
	return nil, nil
}

func SelectExtensionByKeyword(keyword string, pageNum int, pageSize int) (*Page, error) {
	page := Page{PageNum: pageNum, PageSize: pageSize}
	db, err := common.GetDb()
	defer db.Close()
	if err != nil {
		return nil, err
	}
	var where = ""
	var params []interface{}
	if len(strings.TrimSpace(keyword)) > 0 {
		where = " where title like CONCAT('%',?,'%') or description like CONCAT('%',?,'%')"
		params = append(params, []interface{}{keyword, keyword}...)
	}

	stmt, err := db.Prepare("select count(*) from extension" + where)
	if err != nil {
		return nil, err
	}
	var rows *sql.Rows
	var queryErr error
	if len(params) != 0 {
		rows, queryErr = stmt.Query(params...)
	} else {
		rows, queryErr = stmt.Query()
	}
	if queryErr != nil {
		return nil, queryErr
	}
	if rows.Next() {
		var count int
		err = rows.Scan(&count)
		if err != nil {
			return nil, err
		} else {
			page.TotalCount = count
			page.TotalPage = int(math.Ceil(float64(count) / float64(pageSize)))
		}
	}

	stmt, err = db.Prepare("select id,title,version,description,path,files,create_time,update_time from extension" + where + " limit ?,?")
	if err != nil {
		return nil, err
	}
	params = append(params, []interface{}{(pageNum - 1) * pageSize, pageSize}...)
	if len(params) != 0 {
		rows, queryErr = stmt.Query(params...)
	} else {
		rows, queryErr = stmt.Query()
	}
	if queryErr != nil {
		return nil, queryErr
	}
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var extension Extension
		err = rows.Scan(&extension.Id, &extension.Title, &extension.Version, &extension.Description, &extension.Path, &extension.Files, &extension.CreateTime, &extension.UpdateTime)
		if err != nil {
			return nil, err
		} else {
			page.Data = append(page.Data, extension)
		}
	}
	return &page, nil
}

func (extension *Extension) Update() {
	db, err := common.GetDb()
	defer db.Close()
	if err != nil {
		return
	}
	stmt, err := db.Prepare("update extension set title=?,version=?,description=?,files=?,update_time=? where id=?")
	if err != nil {
		return
	}
	stmt.Exec(extension.Title, extension.Version, extension.Description, extension.Files, extension.UpdateTime, extension.Id)
}

func (extension *Extension) Insert() {
	db, err := common.GetDb()
	defer db.Close()
	if err != nil {
		return
	}
	stmt, err := db.Prepare("insert into extension (title,version,description,path,files,create_time,update_time) values (?,?,?,?,?,?,?)")
	if err != nil {
		return
	}
	stmt.Exec(extension.Title, extension.Version, extension.Description, extension.Path, extension.Files, extension.CreateTime, extension.CreateTime)
}
