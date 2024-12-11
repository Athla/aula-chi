package server

import (
	"aula-chi/internal/models"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandler(t *testing.T) {
	s := &Server{}
	server := httptest.NewServer(http.HandlerFunc(s.HelloWorldHandler))
	defer server.Close()
	resp, err := http.Get(server.URL)
	if err != nil {
		t.Fatalf("error making request to server. Err: %v", err)
	}
	defer resp.Body.Close()
	// Assertions
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status OK; got %v", resp.Status)
	}
	expected := "{\"message\":\"Hello World\"}"
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("error reading response body. Err: %v", err)
	}
	if expected != string(body) {
		t.Errorf("expected response body to be %v; got %v", expected, string(body))
	}
}

func setupTestServer() *Server {
	return &Server{
		productsDatabase: []models.Product{
			{
				Id:       1,
				Name:     "Test Product",
				Price:    "99.99",
				Category: "Geral",
				Owner:    1,
			},
		},
		userDatabase: []models.User{
			{
				Id:    1,
				Name:  "Test User",
				Email: "test@example.com",
				Products: []models.Product{
					{
						Id:       1,
						Name:     "Test Product",
						Price:    "99.99",
						Category: "Geral",
						Owner:    1,
					},
				},
			},
		},
	}
}
