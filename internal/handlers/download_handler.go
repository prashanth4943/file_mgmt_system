package handlers

import (
	"bytes"
	"file_mgmt_system/internal/service"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

type DownloadHandler struct {
	Service *service.DownloadService
}

func NewDownloadHandler(service *service.DownloadService) *DownloadHandler {
	return &DownloadHandler{
		Service: service,
	}
}

func (h *DownloadHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	vars := mux.Vars(r)
	fileID, ok := vars["fileID"]
	if !ok || fileID == "" {
		http.Error(w, "fileID is required in the path", http.StatusBadRequest)
		return
	}

	fileName, fileStream, err := h.Service.Download(fileID)
	if err != nil {
		http.Error(w, "Failed to download file: "+err.Error(), http.StatusInternalServerError)
		return
	}
	// defer fileStream.Close()

	w.Header().Set("Content-Disposition", "attachment; filename="+fileName)
	// w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", fileName))
	// w.Header().Set("Content-Type", "application/octet-stream")
	w.WriteHeader(http.StatusOK)
	_, err = io.Copy(w, fileStream)
	if err != nil {
		http.Error(w, "Failed to send file: "+err.Error(), http.StatusInternalServerError)
	}
}

func (h *DownloadHandler) GetThumbnail(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	vars := mux.Vars(r)
	fileID, ok := vars["fileName"]
	if !ok || fileID == "" {
		http.Error(w, "fileName is required in the path", http.StatusBadRequest)
		return
	}

	_, fileData, err := h.Service.Download(fileID) // Assuming this function already exists
	if err != nil {
		http.Error(w, "Failed to fetch file from storage: "+err.Error(), http.StatusInternalServerError)
		return
	}

	fileBytes, err := io.ReadAll(fileData) // Convert io.Reader to []byte
	if err != nil {
		http.Error(w, "Failed to read file data: "+err.Error(), http.StatusInternalServerError)
		return
	}
	// log.Println("prashanth")
	// log.Println(fileBytes)
	thumbnailData, err := h.Service.GenerateThumbnail(fileBytes)
	if err != nil {
		http.Error(w, "Failed to generate thumbnail: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "image/jpeg") // Or the appropriate type
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(thumbnailData)
	if err != nil {
		http.Error(w, "Failed to write thumbnail to response: "+err.Error(), http.StatusInternalServerError)
	}

}

func (h *DownloadHandler) GetPreview(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	vars := mux.Vars(r)
	fileID, ok := vars["fileID"]
	if !ok || fileID == "" {
		http.Error(w, "fileID is required in the path", http.StatusBadRequest)
		return
	}

	fileName, fileData, err := h.Service.Download(fileID) // Assuming this function already exists
	if err != nil {
		http.Error(w, "Failed to fetch file from storage: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// detect MIME type based on file name extension
	var mimeType string
	switch {
	case strings.HasSuffix(strings.ToLower(fileName), ".jpg"),
		strings.HasSuffix(strings.ToLower(fileName), ".jpeg"):
		mimeType = "image/jpeg"
	case strings.HasSuffix(strings.ToLower(fileName), ".png"):
		mimeType = "image/png"
	case strings.HasSuffix(strings.ToLower(fileName), ".gif"):
		mimeType = "image/gif"
	default:
		mimeType = "application/octet-stream" // Fallback for unknown types
	}

	if mimeType == "image/jpeg" || mimeType == "image/png" {
		log.Println("opiopipo")
		compressedData, err := h.Service.CompressImage(fileData, mimeType, 800)
		log.Println(compressedData)
		log.Println(err)
		if err == nil {
			log.Println("here")
			log.Println(compressedData)
			fileData = bytes.NewReader(compressedData)
		}
	}

	// w.Header().Set("Content-Type", "application/octet-stream")               // Use appropriate MIME type if known
	w.WriteHeader(http.StatusOK)
	// w.Header().Set("Content-Type", "image/jpeg")
	w.Header().Set("Content-Type", mimeType)
	w.Header().Set("Cache-Control", "public, max-age=21600")
	w.Header().Set("Content-Disposition", `inline; filename="`+fileName+`"`) // "inline" for preview, "attachment" for download
	// w.Header().Set("Content-Length", fmt.Sprintf("%d", len(fileData)))

	// Write the file data to the response body
	// _, writeErr := w.Write(fileData)
	_, err = io.Copy(w, fileData)
	if err != nil {
		http.Error(w, "Failed to write file to response: "+err.Error(), http.StatusInternalServerError)
		return
	}

}
