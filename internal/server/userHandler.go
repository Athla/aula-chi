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

func (s *Server) GetUsersHandler(w http.ResponseWriter, r *http.Request) {
	utils.RespondWithJSON(w, http.StatusOK, s.userDatabase)
}

func (s *Server) CreateNewUserHandler(w http.ResponseWriter, r *http.Request) {
	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	user.Id = uint(len(s.userDatabase) + 1)
	s.userDatabase = append(s.userDatabase, user)
	utils.RespondWithJSON(w, http.StatusCreated, user)
}

func (s *Server) GetSpecificUserHandler(w http.ResponseWriter, r *http.Request) {
	userId, err := strconv.ParseUint(chi.URLParam(r, "userID"), 10, 32)
	if err != nil {
		log.Println(err)
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid User ID.")
		return
	}
	for _, user := range s.userDatabase {
		if user.Id == uint(userId) {
			utils.RespondWithJSON(w, http.StatusOK, user)
			return
		}
	}
	utils.RespondWithError(w, http.StatusNotFound, "User not found in database.")
}

func (s *Server) EditUserHandler(w http.ResponseWriter, r *http.Request) {
	userId, err := strconv.ParseUint(chi.URLParam(r, "userID"), 10, 32)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid User ID.")
		return
	}

	var updatedUser models.User
	if err := json.NewDecoder(r.Body).Decode(&updatedUser); err != nil {
		log.Println(err)
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request paylod.")
		return
	}

	for i, user := range s.userDatabase {
		if user.Id == uint(userId) {
			updatedUser.Id = user.Id
			s.userDatabase[i] = updatedUser
			utils.RespondWithJSON(w, http.StatusOK, updatedUser)
			return
		}
	}
	utils.RespondWithError(w, http.StatusNotFound, "User not found in database.")
}

func (s *Server) DeleteUserHandler(w http.ResponseWriter, r *http.Request) {
	userId, err := strconv.ParseUint(chi.URLParam(r, "userID"), 10, 32)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid User ID.")
		return
	}

	for i, user := range s.userDatabase {
		if user.Id == uint(userId) {
			s.userDatabase = append(s.userDatabase[:i], s.userDatabase[i+1:]...)
		}
		return
	}
	utils.RespondWithError(w, http.StatusNotFound, "User not found in database.")

}
