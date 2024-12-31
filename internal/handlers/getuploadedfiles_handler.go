package handlers

import (
	"encoding/json"
	"file_mgmt_system/internal/service"
	"net/http"

	"github.com/gorilla/mux"
)

type GetFilesHandler struct {
	service *service.GetFilesService
}

func NewGetFilesHandler(service *service.GetFilesService) *GetFilesHandler {
	return &GetFilesHandler{service: service}
}

func (handler *GetFilesHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	vars := mux.Vars(r)
	email, ok := vars["email"]
	if !ok || email == "" {
		http.Error(w, "Email is required in the path", http.StatusBadRequest)
		return
	}

	enrichedFiles, err := handler.service.GetUploadedFilesByEmail(email)
	if err != nil {
		http.Error(w, "Failed to fetch uploaded files: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	// json.NewEncoder(w).Encode(files)
	json.NewEncoder(w).Encode(enrichedFiles)
}
