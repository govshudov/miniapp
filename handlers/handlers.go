package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"io"
	"miniapp/models"
	"net/http"
	"strconv"
)

type HTTPHandler struct {
	repository Repository
}

func NewHTTPHandler(repository Repository) *HTTPHandler {
	return &HTTPHandler{
		repository: repository,
	}
}
func (handler *HTTPHandler) RegisterRoutes(router chi.Router) {
	router.Method("POST", "/upsert-client", http.HandlerFunc(handler.UpsertClient))
	router.Method("GET", "/get-own-clients/{id}", http.HandlerFunc(handler.GetOwnClients))
	router.Method("POST", "/search-client", http.HandlerFunc(handler.SearchClients))
	router.Method("POST", "/login", http.HandlerFunc(handler.Login))
}
func (handler *HTTPHandler) UpsertClient(w http.ResponseWriter, r *http.Request) {
	body, errBody := io.ReadAll(r.Body)
	if errBody != nil {
		http.Error(w, errBody.Error(), http.StatusBadRequest)
		return
	}

	dto := models.Client{}
	errData := json.Unmarshal(body, &dto)
	if errData != nil || (dto.Currency != "usd" && dto.Currency != "eur") || (dto.USD == 0 && dto.EUR == 0) {
		http.Error(w, errData.Error(), http.StatusBadRequest)
		return
	}
	err := handler.repository.UpsertClient(dto)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func (handler *HTTPHandler) GetOwnClients(w http.ResponseWriter, r *http.Request) {
	idString := chi.URLParam(r, "id")
	id, _ := strconv.Atoi(idString)
	clients, err := handler.repository.GetOwnClients(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	handler.writeJSON(w, clients, http.StatusOK)
}

func (handler *HTTPHandler) SearchClients(w http.ResponseWriter, r *http.Request) {
	body, errBody := io.ReadAll(r.Body)
	if errBody != nil {
		http.Error(w, errBody.Error(), http.StatusBadRequest)
		return
	}

	dto := models.SearchClient{}
	errData := json.Unmarshal(body, &dto)
	if errData != nil {
		http.Error(w, errData.Error(), http.StatusBadRequest)
		return
	}
	fmt.Println(dto)
	clients, err := handler.repository.SearchClient(dto.Passport)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	handler.writeJSON(w, clients, http.StatusOK)
}

func (handler *HTTPHandler) Login(w http.ResponseWriter, r *http.Request) {
	body, errBody := io.ReadAll(r.Body)
	if errBody != nil {
		http.Error(w, errBody.Error(), http.StatusBadRequest)
		return
	}

	dto := models.Login{}
	errData := json.Unmarshal(body, &dto)
	if errData != nil {
		http.Error(w, errData.Error(), http.StatusBadRequest)
		return
	}
	password, err := handler.repository.GetUserData(dto.UserName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if password != dto.Password {
		http.Error(w, errData.Error(), http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	handler.writeJSON(w, "login successfully", http.StatusOK)
}

func (handler *HTTPHandler) writeJSON(w http.ResponseWriter, data interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
	}
}
