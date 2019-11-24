package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/cors"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	"github.com/go-redis/redis"
)

var (
	serviceHost string = os.Getenv("SERVICE_HOST")
	redisHost   string = os.Getenv("REDIS_HOST")
)

var (
	rClient *redis.Client
)

func routes() *chi.Mux {

	router := chi.NewRouter()

	cors := cors.New(cors.Options{
		// AllowedOrigins: []string{"https://foo.com"}, // Use this to allow specific origin hosts
		AllowedOrigins: []string{"*"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token", "Access-Control-Allow-Origin"},
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
	router.Route("/v1", func(r chi.Router) {
		r.Mount("/api/searchresult", handleSearchResult())
		r.Mount("/api/last10", handlelast10())
	})
	return router
}

func main() {
	if "" == serviceHost {
		serviceHost = "localhost:8082"
	}

	if "" == redisHost {
		redisHost = "localhost:6379"
	}
	rClient = redisClient()
	router := routes()
	walkRoutes := func(method string, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		log.Printf("%s %s\n", method, route)
		return nil
	}
	if err := chi.Walk(router, walkRoutes); err != nil {
		log.Panicf("Logging err: %s\n", err.Error())
	}
	log.Fatal(http.ListenAndServe(serviceHost, router))
}

func handleSearchResult() *chi.Mux {
	router := chi.NewRouter()
	router.Get("/{uuid}", getSearchResult)
	return router
}

func handlelast10() *chi.Mux {
	router := chi.NewRouter()
	router.Get("/{uuid}", getLast10Result)
	return router
}

func getSearchResult(w http.ResponseWriter, r *http.Request) {

	uuid := chi.URLParam(r, "uuid")
	cachedResult, err := rClient.Get(uuid).Result()
	if err != nil {
		log.Print("Error fetching cached result for ", uuid, "\n", err)
		return
	}
	var result map[string]interface{}
	json.Unmarshal([]byte(cachedResult), &result)
	fmt.Print("Rendering results \n")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	render.JSON(w, r, result)
}

func getLast10Result(w http.ResponseWriter, r *http.Request) {

	uuid := chi.URLParam(r, "uuid")
	cachedResult, err := rClient.LRange("uuid", -10, -1).Result()

	if err != nil {
		log.Print("Error fetching cached result for ", uuid)
		return
	}
	var result map[string]interface{}
	var last10 [10]map[string]interface{}
	for index, element := range cachedResult {
		json.Unmarshal([]byte(element), &result)
		r, ok := result["uuid"].(map[string]interface{})
		if ok {
			last10[index] = r
		}
	}
	render.JSON(w, r, last10)
}

func redisClient() *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     redisHost,
		Password: "",
		DB:       0,
	})
	pong, err := client.Ping().Result()
	if err != nil {
		print("Failed to connect to redis store, error - ", err, "\n")
	} else {
		print("Connected to redis store, ping result - ", pong, "\n")
	}
	return client
}
