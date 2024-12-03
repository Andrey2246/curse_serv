package handlers

import (
	"net/http"
	"os"
	"path/filepath"
)

func DownloadHandler(storageDir string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
			return
		}

		filename := r.URL.Query().Get("filename")
		if filename == "" {
			http.Error(w, "Filename is required", http.StatusBadRequest)
			return
		}

		filePath := filepath.Join(storageDir, filename)
		file, err := os.Open(filePath)
		if err != nil {
			http.Error(w, "File not found", http.StatusNotFound)
			return
		}
		defer file.Close()

		http.ServeFile(w, r, filePath)
	}
}
