package models

import (
	"net/http"
	"strings"
)

type Page struct {
	PageNum    int           `json:"pageNum"`    //当前页数
	PageSize   int           `json:"pageSize"`   //每页条数
	TotalPage  int           `json:"totalPage"`  //总页数
	TotalCount int           `json:"totalCount"` //总数
	Data       []interface{} `json:"data"`       //数据
}

func GetIp(r *http.Request) string {
	ip := r.Header.Get("X-Forwarded-For")
	if len(ip) > 0 && "unknown" != ip && "127.0.0.1" != ip {
		// 多次反向代理后会有多个IP值，第一个为真实IP。
		index := strings.Index(ip, ",")
		if index != -1 {
			return ip[0:index]
		} else {
			return ip
		}
	} else {
		return strings.Split(r.RemoteAddr, ":")[0]
	}
}

func SetCORS(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "*")
	w.Header().Set("Access-Control-Allow-Headers", "*")
}
