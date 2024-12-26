package handlers

import (
	"encoding/json"
	"file_mgmt_system/helper"
	"file_mgmt_system/internal/models"
	"file_mgmt_system/internal/service"
	"net/http"
)

type UploadHandler struct {
	service *service.UploadService
}

func NewUploadHandler(service *service.UploadService) *UploadHandler {
	return &UploadHandler{
		service: service,
	}
}

func (handler *UploadHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Failed to read file: "+err.Error(), http.StatusBadRequest)
		return
	}
	defer file.Close()

	email, ok := helper.GetEmailFromContext(r.Context())

	if !ok {
		email := r.FormValue("email")
		if email == "" {
			writeErrorResponse(w, "Email is required", "ERR_MISSING_EMAIL", http.StatusBadRequest)
			return
		}
	}
	req := models.UploadRequest{
		FileName: header.Filename,
		FileType: header.Header.Get("Content-Type"),
		FileSize: header.Size,
		Email:    email,
	}

	// Call the service layer
	fileMetadata, err := handler.service.Upload(file, header, req)
	if err != nil {
		writeErrorResponse(w, "Failed to upload file: "+err.Error(), "ERR_UPLOAD_FAILED", http.StatusInternalServerError)
		return
	}
	response := map[string]interface{}{
		"status": "success",
		"file":   fileMetadata,
	}
	writeJSONResponse(w, response, http.StatusOK)
}

func writeErrorResponse(w http.ResponseWriter, message, code string, status int) {
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "error",
		"error": map[string]interface{}{
			"message": message,
			"code":    code,
		},
	})
}

func writeJSONResponse(w http.ResponseWriter, data map[string]interface{}, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}
