package common

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
)

type configuration struct {
	Port int
	App  struct {
		Server string
	}
	Extension struct {
		Secret string
		Repo   string
	}
	Db struct {
		Host     string
		Port     int
		Uri      string
		Username string
		Password string
	}
}

var Config configuration

//加载配置文件
func init() {
	file, _ := os.Open("config.json")
	defer file.Close()

	decoder := json.NewDecoder(file)
	err := decoder.Decode(&Config)
	if err != nil {
		fmt.Println("Error:", err)
	}
}

func GetDb() (*sql.DB, error) {
	return sql.Open("mysql", Config.Db.Username+":"+Config.Db.Password+"@tcp("+Config.Db.Host+":"+strconv.Itoa(Config.Db.Port)+")"+Config.Db.Uri+"?parseTime=true&loc=Local")
}
