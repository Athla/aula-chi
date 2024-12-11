package server

import (
	"aula-chi/internal/models"
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
)

func TestGetProductsHandler(t *testing.T) {
	server := setupTestServer()
	req := httptest.NewRequest(http.MethodGet, "/products", nil)
	rr := httptest.NewRecorder()

	server.GetProductsHandler(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var products []models.Product
	err := json.Unmarshal(rr.Body.Bytes(), &products)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if len(products) != len(server.productsDatabase) {
		t.Errorf("handler returned wrong number of products: got %v want %v",
			len(products), len(server.productsDatabase))
	}
}

func TestCreateNewProductHandler(t *testing.T) {
	server := setupTestServer()
	newProduct := models.Product{
		Name:  "New Product",
		Price: "29.99",
	}

	body, _ := json.Marshal(newProduct)
	req := httptest.NewRequest(http.MethodPost, "/products", bytes.NewBuffer(body))
	rr := httptest.NewRecorder()

	server.CreateNewProductHandler(rr, req)

	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusCreated)
	}

	var createdProduct models.Product
	err := json.Unmarshal(rr.Body.Bytes(), &createdProduct)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if createdProduct.Id != 2 {
		t.Errorf("handler returned wrong product ID: got %v want %v", createdProduct.Id, 2)
	}
}

func TestGetSpecificProductHandler(t *testing.T) {
	server := setupTestServer()
	req := httptest.NewRequest(http.MethodGet, "/products/1", nil)
	rr := httptest.NewRecorder()

	chiCtx := chi.NewRouteContext()
	chiCtx.URLParams.Add("productID", "1")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, chiCtx))

	server.GetSpecificProductHandler(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var product models.Product
	err := json.Unmarshal(rr.Body.Bytes(), &product)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if product.Id != 1 {
		t.Errorf("handler returned wrong product: got id %v want id %v", product.Id, 1)
	}
}

func TestEditProductHandler(t *testing.T) {
	server := setupTestServer()
	updatedProduct := models.Product{
		Name:  "Updated Product",
		Price: "199.99",
	}

	body, _ := json.Marshal(updatedProduct)
	req := httptest.NewRequest(http.MethodPut, "/products/1", bytes.NewBuffer(body))
	rr := httptest.NewRecorder()

	chiCtx := chi.NewRouteContext()
	chiCtx.URLParams.Add("productID", "1")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, chiCtx))

	server.EditProductHandler(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var product models.Product
	err := json.Unmarshal(rr.Body.Bytes(), &product)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if product.Name != updatedProduct.Name {
		t.Errorf("handler returned wrong product name: got %v want %v",
			product.Name, updatedProduct.Name)
	}
}

func TestDeleteProductHandler(t *testing.T) {
	server := setupTestServer()
	req := httptest.NewRequest(http.MethodDelete, "/products/1", nil)
	rr := httptest.NewRecorder()

	chiCtx := chi.NewRouteContext()
	chiCtx.URLParams.Add("productID", "1")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, chiCtx))

	server.DeleteProductHandler(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	if len(server.productsDatabase) != 0 {
		t.Errorf("product was not deleted: database contains %v products",
			len(server.productsDatabase))
	}
}
