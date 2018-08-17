package handle

import (
	"fmt"
	"io"
	"net/http"
)

func WebHook(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	fmt.Printf("%v\n", r.Form)
	io.WriteString(w, "ok")
}
