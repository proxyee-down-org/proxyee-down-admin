package handles

import (
	"encoding/json"
	"io"
	"models"
	"net/http"
)

func CheckUpdate(w http.ResponseWriter, r *http.Request) {
	models.SetCORS(w)
	version, err := models.GetNewVersion()
	if err != nil {
		w.WriteHeader(500)
		io.WriteString(w, "error")
	} else {
		bts, _ := json.Marshal(version)
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.Write(bts)
	}
}
