package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/dgraph-io/dgo/protos/api"
	"github.com/julienschmidt/httprouter"

	"github.com/kind84/hlx/models"
	"github.com/kind84/hlx/repo"
)

func TestGetFolder(t *testing.T) {
	r := &repo.Repo{
		ConnStr: "localhost:9080",
	}
	defer func() {
		c := r.NewClient()
		err := c.Alter(context.Background(), &api.Operation{DropAll: true})
		if err != nil {
			t.Errorf("Error deleteing data")
		}
	}()

	router := httprouter.New()
	router.POST("/api/categories/load", LoadCategories(r))
	router.POST("/api/categories/leaves", GetCategoryLeaves(r))

	rr := httptest.NewRecorder()

	jf, err := os.Open("../categories.json")
	if err != nil {
		t.Errorf("Cannot open json file.")
	}
	defer jf.Close()

	req, _ := http.NewRequest("POST", "/api/categories/load", jf)

	router.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Wrong status, got %v, want %v\n", status, http.StatusOK)
	}

	decoder := json.NewDecoder(rr.Body)
	var res []models.Category
	decoder.Decode(&res)

	if len(res) != 90 {
		t.Errorf("Expected 90 results, got %d", len(res))
	}

	body := struct {
		Name         string `json:"name"`
		SubLayerName string `json:"subLayerName,omitempty"`
	}{
		Name:         "Comics",
		SubLayerName: "Books & Literature",
	}
	bs, err := json.Marshal(body)

	req, _ = http.NewRequest("POST", "/api/categories/leaves", bytes.NewBuffer(bs))

	router.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Wrong status, got %v, want %v\n", status, http.StatusOK)
	}

	decoder = json.NewDecoder(rr.Body)
	decoder.Decode(&res)

	if len(res) != 1 {
		t.Errorf("Expected 1 result, got %d", len(res))
	}
	if res[0].ID != 60 {
		t.Errorf("Wrong ID: expected 60, got %v", res[0].ID)
	}
}
