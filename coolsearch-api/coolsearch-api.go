package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/render"
	"gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"
)

// SearchResult stores the jon message from 'searchengine-result' topic
type search struct {
	UUID       string `json:"uuid"`
	Searchword string `json:"searchword"`
}

func routes() *chi.Mux {
	router := chi.NewRouter()

	// Basic CORS
	cors := cors.New(cors.Options{
		// AllowedOrigins: []string{"https://foo.com"}, // Use this to allow specific origin hosts
		AllowedOrigins: []string{"*"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	})
	router.Use(cors.Handler)
	router.Use(
		render.SetContentType(render.ContentTypeJSON),
		middleware.Logger,
		middleware.DefaultCompress,
		middleware.RedirectSlashes,
		middleware.Recoverer,
	)
	// add sub-routers here
	// sub-router url started with prefix v1
	/*router.Route("/v1", func(r chi.Router) {
		r.Mount("/api/postquery", handleSearch())
	})*/
	router.Post("/api/postquery", postNewQuery)
	return router
}

func main() {
	router := routes()
	walkRoutes := func(method string, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		log.Printf("%s %s\n", method, route)
		return nil
	}
	if err := chi.Walk(router, walkRoutes); err != nil {
		log.Panicf("Logging err: %s\n", err.Error())
	}
	log.Fatal(http.ListenAndServe(":8082", router))
}

func postNewQuery(w http.ResponseWriter, r *http.Request) {
	var post search
	json.NewDecoder(r.Body).Decode(&post)
	query, err := json.Marshal(post)
	if err == nil {
		postSearchResult("search", query)
	} else {
		log.Println("Error decoding search query, ", err)
	}
}

func postSearchResult(topic string, data []byte) {
	p, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": "localhost:9092"})
	if err != nil {
		panic(err)
	}
	p.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          data,
	}, nil)
}
