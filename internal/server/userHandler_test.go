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

func TestGetUsersHandler(t *testing.T) {
	server := setupTestServer()
	req := httptest.NewRequest(http.MethodGet, "/users", nil)
	rr := httptest.NewRecorder()

	server.GetUsersHandler(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var users []models.User
	err := json.Unmarshal(rr.Body.Bytes(), &users)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if len(users) != len(server.userDatabase) {
		t.Errorf("handler returned wrong number of users: got %v want %v",
			len(users), len(server.userDatabase))
	}
}

func TestCreateNewUserHandler(t *testing.T) {
	server := setupTestServer()
	newUser := models.User{
		Id:    2,
		Name:  "New User",
		Email: "new@example.com",
	}

	body, _ := json.Marshal(newUser)
	req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewBuffer(body))
	rr := httptest.NewRecorder()

	server.CreateNewUserHandler(rr, req)

	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusCreated)
	}

	var createdUser models.User
	err := json.Unmarshal(rr.Body.Bytes(), &createdUser)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if createdUser.Id != 2 {
		t.Errorf("handler returned wrong user ID: got %v want %v", createdUser.Id, 2)
	}
}

func TestGetSpecificUserHandler(t *testing.T) {
	server := setupTestServer()
	req := httptest.NewRequest(http.MethodGet, "/users/1", nil)
	rr := httptest.NewRecorder()

	chiCtx := chi.NewRouteContext()
	chiCtx.URLParams.Add("userID", "1")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, chiCtx))

	server.GetSpecificUserHandler(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var user models.User
	err := json.Unmarshal(rr.Body.Bytes(), &user)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if user.Id != 1 {
		t.Errorf("handler returned wrong user: got id %v want id %v", user.Id, 1)
	}
}

func TestEditUserHandler(t *testing.T) {
	server := setupTestServer()
	updatedUser := models.User{
		Id:    1,
		Name:  "Updated User",
		Email: "updated@example.com",
	}

	body, _ := json.Marshal(updatedUser)
	req := httptest.NewRequest(http.MethodPut, "/users/1", bytes.NewBuffer(body))
	rr := httptest.NewRecorder()

	chiCtx := chi.NewRouteContext()
	chiCtx.URLParams.Add("userID", "1")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, chiCtx))

	server.EditUserHandler(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var user models.User
	err := json.Unmarshal(rr.Body.Bytes(), &user)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if user.Name != updatedUser.Name || user.Email != updatedUser.Email {
		t.Errorf("handler returned wrong user data: got %v want %v", user, updatedUser)
	}
}

func TestDeleteUserHandler(t *testing.T) {
	server := setupTestServer()
	req := httptest.NewRequest(http.MethodDelete, "/users/1", nil)
	rr := httptest.NewRecorder()

	chiCtx := chi.NewRouteContext()
	chiCtx.URLParams.Add("userID", "1")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, chiCtx))

	server.DeleteUserHandler(rr, req)

	if len(server.userDatabase) != 0 {
		t.Errorf("user was not deleted: database contains %v users",
			len(server.userDatabase))
	}
}

func TestInvalidUserID(t *testing.T) {
	server := setupTestServer()
	tests := []struct {
		name       string
		userID     string
		handler    http.HandlerFunc
		method     string
		wantStatus int
	}{
		{
			name:       "Get Invalid User ID",
			userID:     "invalid",
			handler:    server.GetSpecificUserHandler,
			method:     http.MethodGet,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "Edit Invalid User ID",
			userID:     "invalid",
			handler:    server.EditUserHandler,
			method:     http.MethodPut,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "Delete Invalid User ID",
			userID:     "invalid",
			handler:    server.DeleteUserHandler,
			method:     http.MethodDelete,
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, "/users/"+tt.userID, nil)
			rr := httptest.NewRecorder()

			chiCtx := chi.NewRouteContext()
			chiCtx.URLParams.Add("userID", tt.userID)
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, chiCtx))

			tt.handler(rr, req)

			if status := rr.Code; status != tt.wantStatus {
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, tt.wantStatus)
			}
		})
	}
}
