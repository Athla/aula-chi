package server

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func (s *Server) RegisterRoutes() http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Get("/", s.HelloWorldHandler)
	r.Route("/users", func(r chi.Router) {
		r.Get("/", s.GetUsersHandler)
		r.Post("/", s.CreateNewUserHandler)
		r.Route("/{userID}", func(r chi.Router) {
			r.Get("/", s.GetSpecificUserHandler)
			r.Put("/", s.EditUserHandler)
			r.Delete("/", s.DeleteUserHandler)
		})
	})
	r.Route("/products", func(r chi.Router) {
		r.Get("/", s.GetProductsHandler)
		r.Post("/", s.CreateNewProductHandler)
		r.Route("/{ProductID}", func(r chi.Router) {
			r.Get("/", s.GetSpecificProductHandler)
			r.Put("/", s.EditProductHandler)
			r.Delete("/", s.DeleteProductHandler)
		})

	})

	log.Println("Server alive and running!")
	return r
}

func (s *Server) HelloWorldHandler(w http.ResponseWriter, r *http.Request) {
	resp := make(map[string]string)
	resp["message"] = "Hello World"

	jsonResp, err := json.Marshal(resp)
	if err != nil {
		log.Fatalf("error handling JSON marshal. Err: %v", err)
	}

	_, _ = w.Write(jsonResp)
}
