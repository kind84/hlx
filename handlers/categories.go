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

func LoadCategories(r *repo.Repo) func(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	return func(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
		ctx := context.TODO()

		var cats []models.Category

		defer req.Body.Close()

		decoder := json.NewDecoder(req.Body)
		err := decoder.Decode(&cats)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(`{"code": 400, "msg": "Bad Request."}`))
			return
		}

		cs, err := r.SaveCategories(ctx, cats)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf(`{"code": 500, "msg": "Internal Server Error: %s"}`, err.Error())))
			return
		}

		encoder := json.NewEncoder(w)
		encoder.Encode(cs)
	}
}

func GetCategoryLeaves(r *repo.Repo) func(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	return func(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
		ctx := context.TODO()

		var payload struct {
			Name         string `json:"name"`
			SubLayerName string `json:"subLayerName,omitempty"`
		}

		defer req.Body.Close()

		decoder := json.NewDecoder(req.Body)
		err := decoder.Decode(&payload)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(`{"code": 400, "msg": "Bad Request."}`))
			return
		}

		cs, err := r.GetCategoryLeaf(ctx, payload.Name, payload.SubLayerName)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf(`{"code": 500, "msg": "Internal Server Error: %s"}`, err.Error())))
			return
		}

		encoder := json.NewEncoder(w)
		encoder.Encode(cs)
	}
}
