package handles

import (
	"bdy"
	"encoding/json"
	"io"
	"io/ioutil"
	"models"
	"net/http"
)

func BdyResolve(w http.ResponseWriter, r *http.Request) {
	models.SetCORS(w)
	payloadBts, _ := ioutil.ReadAll(r.Body)
	if len(payloadBts) == 0 {
		w.WriteHeader(400)
		io.WriteString(w, "params error")
		return
	}
	var params map[string]interface{}
	err := json.Unmarshal(payloadBts, &params)
	if err != nil {
		w.WriteHeader(400)
		io.WriteString(w, "params error")
		return
	}
	if err != nil {
		w.WriteHeader(400)
		io.WriteString(w, "params error")
		return
	}
	result, err := bdy.Resolve(params["url"].(string), params["rand"].(string), params["sign"].(string))
	if err != nil {
		w.WriteHeader(500)
		io.WriteString(w, err.Error())
		return
	}
	bts, err := json.Marshal(result["urls"])
	if err != nil {
		w.WriteHeader(500)
		io.WriteString(w, err.Error())
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Write(bts)
}
