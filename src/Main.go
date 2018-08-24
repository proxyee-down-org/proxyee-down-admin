package main

import (
	"common"
	"handles"
	"net/http"
	"strconv"
)

func main() {
	http.HandleFunc("/extension/webhook", handles.WebHook)
	http.HandleFunc("/extension/search", handles.Search)
	http.HandleFunc("/extension/down", handles.Down)
	http.ListenAndServe(":"+strconv.Itoa(common.Config.Port), nil)
}
