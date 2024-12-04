package handlers

import (
	"curse_serv/logger"
	"fmt"
	"log"
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
			http.Error(w, "File not found", 404)
			log.Printf("Failed to open file '%s': %v", filePath, err)
			return
		}
		defer file.Close()

		if err := logger.Log("Download", fmt.Sprintf("File '%s' downloaded by client", filename)); err != nil {
			log.Printf("Failed to log download action: %v", err)
		}

		http.ServeFile(w, r, filePath)
	}
}
