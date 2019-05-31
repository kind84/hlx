package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/kind84/hlx/models"
	"github.com/kind84/hlx/repo"
)

func LoadPsychos(r *repo.Repo) func(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	return func(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
		ctx := context.TODO()

		var ps []models.Psycho

		defer req.Body.Close()

		decoder := json.NewDecoder(req.Body)
		err := decoder.Decode(&ps)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(`{"code": 400, "msg": "Bad Request."}`))
			return
		}

		cs, err := r.SavePsychos(ctx, ps)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf(`{"code": 500, "msg": "Internal Server Error: %s"}`, err.Error())))
			return
		}

		encoder := json.NewEncoder(w)
		encoder.Encode(cs)
	}
}
