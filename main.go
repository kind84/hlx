package main

import (
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/rs/cors"

	"github.com/kind84/hlx/handlers"
	"github.com/kind84/hlx/repo"
)

func main() {
	r := &repo.Repo{
		ConnStr: "localhost:9080",
	}

	mux := httprouter.New()

	mux.GET("/api", handlers.GetInfo)
	mux.POST("/api/categories/load", handlers.LoadCategories(r))

	c := cors.New(cors.Options{
		AllowCredentials: true,
		AllowedHeaders:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST"},
		AllowedOrigins:   []string{"*"},
	})

	handler := c.Handler(mux)

	fmt.Println("Listening on port :8080")
	http.ListenAndServe(":8080", handler)
}