package handles

import (
	"common"
	"encoding/json"
	"io"
	"models"
	"net/http"
	"strconv"
	"strings"
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

func Download(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	version, err := models.GetNewVersion()
	if err != nil {
		w.WriteHeader(500)
		io.WriteString(w, "error")
	} else {
		os := r.Form.Get("os")
		if len(os) == 0 {
			os = "windows"
		}
		v := strconv.FormatFloat(version.Version, 'f', -1, 64)
		if strings.Index(v, ".") == -1 {
			v += ".0"
		}
		t := "zip"
		if os == "windows" {
			t = "7z"
		}
		w.Header().Set("Location", common.Config.App.Server+"/dist/release/"+v+"/"+"Proxyee Down."+v+"."+os+"."+t)
		w.WriteHeader(302)
	}
}
