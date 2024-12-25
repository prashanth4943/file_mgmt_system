package handlers

import (
	"file_mgmt_system/internal/storage"
	"io"
	"net/http"
)

type DownloadHandler struct {
	Storage storage.Storage
}

func NewDownloadHandler(s storage.Storage) *DownloadHandler {
	return &DownloadHandler{Storage: s}
}

func (h *DownloadHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	objectName := r.URL.Query().Get("name")
	if objectName == "" {
		http.Error(w, "File name is required", http.StatusBadRequest)
		return
	}

	reader, err := h.Storage.DownloadFile(objectName)
	if err != nil {
		http.Error(w, "Failed to download file: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Disposition", "attachment; filename="+objectName)
	w.WriteHeader(http.StatusOK)
	io.Copy(w, reader)
}
