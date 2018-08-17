package main

import (
	"handle"
	"net/http"
	"os"
)

func main() {
	http.HandleFunc("/extension/webhook", handle.WebHook)
	http.ListenAndServe(":"+os.Args[1], nil)
}
