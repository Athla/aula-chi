package server

import (
	"aula-chi/internal/models"
	"aula-chi/internal/utils"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

func (s *Server) GetProductsHandler(w http.ResponseWriter, r *http.Request) {
	utils.RespondWithJSON(w, http.StatusOK, s.productsDatabase)
}

func (s *Server) CreateNewProductHandler(w http.ResponseWriter, r *http.Request) {
	var product models.Product
	if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	product.Id = uint(len(s.productsDatabase) + 1)
	s.productsDatabase = append(s.productsDatabase, product)
	utils.RespondWithJSON(w, http.StatusCreated, product)
}

func (s *Server) GetSpecificProductHandler(w http.ResponseWriter, r *http.Request) {
	productId, err := strconv.ParseUint(chi.URLParam(r, "productID"), 10, 32)
	if err != nil {
		log.Println(err)
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid product ID.")
		return
	}
	for _, product := range s.productsDatabase {
		if product.Id == uint(productId) {
			utils.RespondWithJSON(w, http.StatusOK, product)
			return
		}
	}
	utils.RespondWithError(w, http.StatusNotFound, "product not found in database.")
}

func (s *Server) EditProductHandler(w http.ResponseWriter, r *http.Request) {
	productId, err := strconv.ParseUint(chi.URLParam(r, "productID"), 10, 32)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid product ID.")
		return
	}

	var updatedProducts models.Product
	if err := json.NewDecoder(r.Body).Decode(&updatedProducts); err != nil {
		log.Println(err)
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request paylod.")
		return
	}

	for i, product := range s.productsDatabase {
		if product.Id == uint(productId) {
			updatedProducts.Id = product.Id
			s.productsDatabase[i] = updatedProducts
			utils.RespondWithJSON(w, http.StatusOK, updatedProducts)
			return
		}
	}
	utils.RespondWithError(w, http.StatusNotFound, "product not found in database.")
}

func (s *Server) DeleteProductHandler(w http.ResponseWriter, r *http.Request) {
	productId, err := strconv.ParseUint(chi.URLParam(r, "productID"), 10, 32)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid Product ID.")
		return
	}

	for i, prod := range s.productsDatabase {
		if prod.Id == uint(productId) {
			s.productsDatabase = append(s.productsDatabase[:i], s.productsDatabase[i+1:]...)
		}
		return
	}
	utils.RespondWithError(w, http.StatusNotFound, "Product not found in database.")

}
