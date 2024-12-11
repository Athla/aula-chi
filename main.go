package main

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"net/http"
	"strconv"
)

type Product struct {
	Id       uint   `json:"id"`
	Name     string `json:"name"`
	Price    string `json:"price"`
	Category string `json:"category"`
	Owner    string `json:"owner"`
}

type User struct {
	Id       uint      `json:"id"`
	Name     string    `json:"name"`
	Email    string    `json:"email"`
	Products []Product `json:"products"`
}

// Wrapper/Abstracao para a resposta de erro
func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

// Wrapper/Abstracao para a resposta de JSON
func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func main() {
	var usersDatabase []User
	var productDatabase []Product

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.AllowContentType("application/json"))

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello World!"))
	})

	// Top level route definida para manter a logica num so lugar. Em um projeto real, provavel que estaria num arquivo separado (pt 2 se der tempo)
	r.Route("/users", func(r chi.Router) {
		// Retorna todos os users
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			respondWithJSON(w, http.StatusOK, usersDatabase)
		})

		// Cria novo user, fazendo a json validation atraves do decoder, usando um pointer ao user
		r.Post("/", func(w http.ResponseWriter, r *http.Request) {
			var user User
			if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
				respondWithError(w, http.StatusBadRequest, "Invalid request payload")
				return
			}
			user.Id = uint(len(usersDatabase) + 1)
			usersDatabase = append(usersDatabase, user)
			respondWithJSON(w, http.StatusCreated, user)
		})

		// Nesting de rotas para as funcionalidades referente a user spec
		// Util quando sao mts funcionalidade
		r.Route("/{userID}", func(r chi.Router) {
			r.Get("/", func(w http.ResponseWriter, r *http.Request) {
				userID, err := strconv.ParseUint(chi.URLParam(r, "userID"), 10, 32)
				if err != nil {
					respondWithError(w, http.StatusBadRequest, "Invalid user ID")
					return
				}

				for _, user := range usersDatabase {
					if user.Id == uint(userID) {
						respondWithJSON(w, http.StatusOK, user)
						return
					}
				}
				respondWithError(w, http.StatusNotFound, "User not found")
			})

			r.Put("/", func(w http.ResponseWriter, r *http.Request) {
				userID, err := strconv.ParseUint(chi.URLParam(r, "userID"), 10, 32)
				if err != nil {
					respondWithError(w, http.StatusBadRequest, "Invalid user ID")
					return
				}

				var updatedUser User
				if err := json.NewDecoder(r.Body).Decode(&updatedUser); err != nil {
					respondWithError(w, http.StatusBadRequest, "Invalid request payload")
					return
				}

				for i, user := range usersDatabase {
					if user.Id == uint(userID) {
						updatedUser.Id = user.Id
						usersDatabase[i] = updatedUser
						respondWithJSON(w, http.StatusOK, updatedUser)
						return
					}
				}
				respondWithError(w, http.StatusNotFound, "User not found")
			})

			r.Delete("/", func(w http.ResponseWriter, r *http.Request) {
				userID, err := strconv.ParseUint(chi.URLParam(r, "userID"), 10, 32)
				if err != nil {
					respondWithError(w, http.StatusBadRequest, "Invalid user ID")
					return
				}

				for i, user := range usersDatabase {
					if user.Id == uint(userID) {
						usersDatabase = append(usersDatabase[:i], usersDatabase[i+1:]...)
						respondWithJSON(w, http.StatusOK, map[string]string{"message": "User deleted"})
						return
					}
				}
				respondWithError(w, http.StatusNotFound, "User not found")
			})
		})
	})

	// Outra top level route, definida para desacoplar as coisas
	r.Route("/products", func(r chi.Router) {
		// Pega todos os produtos
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			respondWithJSON(w, http.StatusOK, productDatabase)
		})

		// Cria um novo produto
		r.Post("/", func(w http.ResponseWriter, r *http.Request) {
			var product Product
			if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
				respondWithError(w, http.StatusBadRequest, "Invalid request payload")
				return
			}
			product.Id = uint(len(productDatabase) + 1)
			productDatabase = append(productDatabase, product)

			// Update user's products
			for i, user := range usersDatabase {
				if user.Name == product.Owner {
					usersDatabase[i].Products = append(usersDatabase[i].Products, product)
					break
				}
			}

			respondWithJSON(w, http.StatusCreated, product)
		})

		// Funcionalidades referente a item especifico
		r.Route("/{productID}", func(r chi.Router) {
			r.Get("/", func(w http.ResponseWriter, r *http.Request) {
				productID, err := strconv.ParseUint(chi.URLParam(r, "productID"), 10, 32)
				if err != nil {
					respondWithError(w, http.StatusBadRequest, "Invalid product ID")
					return
				}

				for _, product := range productDatabase {
					if product.Id == uint(productID) {
						respondWithJSON(w, http.StatusOK, product)
						return
					}
				}
				respondWithError(w, http.StatusNotFound, "Product not found")
			})

			r.Put("/", func(w http.ResponseWriter, r *http.Request) {
				productID, err := strconv.ParseUint(chi.URLParam(r, "productID"), 10, 32)
				if err != nil {
					respondWithError(w, http.StatusBadRequest, "Invalid product ID")
					return
				}

				var updatedProduct Product
				if err := json.NewDecoder(r.Body).Decode(&updatedProduct); err != nil {
					respondWithError(w, http.StatusBadRequest, "Invalid request payload")
					return
				}

				for i, product := range productDatabase {
					if product.Id == uint(productID) {
						updatedProduct.Id = product.Id
						productDatabase[i] = updatedProduct
						respondWithJSON(w, http.StatusOK, updatedProduct)
						return
					}
				}
				respondWithError(w, http.StatusNotFound, "Product not found")
			})

			r.Delete("/", func(w http.ResponseWriter, r *http.Request) {
				productID, err := strconv.ParseUint(chi.URLParam(r, "productID"), 10, 32)
				if err != nil {
					respondWithError(w, http.StatusBadRequest, "Invalid product ID")
					return
				}

				for i, product := range productDatabase {
					if product.Id == uint(productID) {
						productDatabase = append(productDatabase[:i], productDatabase[i+1:]...)
						respondWithJSON(w, http.StatusOK, map[string]string{"message": "Product deleted"})
						return
					}
				}
				respondWithError(w, http.StatusNotFound, "Product not found")
			})
		})
	})

	fmt.Println("Server running on port 3000")
	http.ListenAndServe(":3000", r)
}
