package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/nmowens95/Chirpy-Web-Server/internal/database"
)

type apiConfig struct {
	fileserverHits int
	DB             *database.DB
}

// using our own files for a DB currently

func main() {
	const filePathRoot = "."
	const port = "8080"

	db, err := database.NewDB("database.json")
	if err != nil {
		log.Fatal(err)
	}

	apiCfg := apiConfig{
		fileserverHits: 0,
		DB:             db,
	}

	router := chi.NewRouter()
	fsHandler := apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(filePathRoot))))
	router.Handle("/app", fsHandler)
	router.Handle("/app/*", fsHandler)
	// Front end focus

	apiRouter := chi.NewRouter()
	apiRouter.Get("/healthz", handlerReadiness)
	apiRouter.Get("/reset", apiCfg.handlerReset)

	apiRouter.Post("/login", apiCfg.handlerLogin)
	apiRouter.Post("/users", apiCfg.handlerUsersCreate)

	apiRouter.Post("/chirps", apiCfg.handlerChirpsCreate)
	apiRouter.Get("/chirps", apiCfg.handlerChirpsRetrieve)
	apiRouter.Get("/chirps/{chirpID}", apiCfg.handlerChirpsGet)
	router.Mount("/api", apiRouter)
	// Backend focus

	adminRouter := chi.NewRouter()
	adminRouter.Get("/metrics", apiCfg.handlerMetrics)
	router.Mount("/admin", adminRouter)
	// Backend focus

	corsMux := middlewareCors(router) // Cross-Origin-Resource-Sharing

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: corsMux,
	}

	log.Printf("Serving files from %s on port: %s\n", filePathRoot, port)
	log.Fatal(srv.ListenAndServe())
}
