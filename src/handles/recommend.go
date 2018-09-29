package handles

import (
	"encoding/json"
	"io"
	"models"
	"net/http"
)

func RecommendSoft(w http.ResponseWriter, r *http.Request) {
	models.SetCORS(w)
	list, err := models.SelectRecommendSoft()
	if err != nil {
		w.WriteHeader(500)
		io.WriteString(w, "error")
	} else {
		bts, _ := json.Marshal(list)
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.Write(bts)
	}
}
