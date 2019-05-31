package handlers

import (
	"io"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func GetInfo(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	io.WriteString(w, "Hello from HLX")
}
