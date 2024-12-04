package handlers

import (
	"curse_serv/logger"
	"database/sql"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"

	_ "github.com/lib/pq"
)

func UploadHandler(storageDir string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
			return
		}

		// Parse the multipart form
		err := r.ParseMultipartForm(20 << 20) // 20 MB max
		if err != nil {
			http.Error(w, "Failed to parse form", http.StatusBadRequest)
			return
		}

		// Retrieve file and metadata
		file, header, err := r.FormFile("file")
		if err != nil {
			http.Error(w, "Failed to retrieve file", http.StatusBadRequest)
			return
		}
		defer file.Close()

		uploader := r.FormValue("uploader")
		if uploader == "" {
			http.Error(w, "Uploader name is required", http.StatusBadRequest)
			return
		}

		// Save the file
		destPath := filepath.Join(storageDir, header.Filename)
		destFile, err := os.Create(destPath)
		if err != nil {
			http.Error(w, "Failed to save file", http.StatusInternalServerError)
			return
		}
		defer destFile.Close()

		size, err := io.Copy(destFile, file)
		if err != nil {
			http.Error(w, "Failed to save file", http.StatusInternalServerError)
			return
		}

		// Log metadata in the database
		db, err := sql.Open("postgres", "postgresql://cursework:security@localhost:5432/")
		if err != nil {
			http.Error(w, "Failed to connect to the database", http.StatusInternalServerError)
			return
		}
		defer db.Close()

		_, err = db.Exec(
			"INSERT INTO files (filename, timestamp, uploader, size) VALUES ($1, NOW(), $2, $3)",
			header.Filename, uploader, size,
		)
		if err != nil {
			http.Error(w, "Failed to save metadata", http.StatusInternalServerError)
			return
		}

		if err := logger.Log("Upload", fmt.Sprintf("File '%s' uploaded by '%s'", header.Filename, uploader)); err != nil {
			log.Printf("Failed to log upload action: %v", err)
		}

		fmt.Fprintf(w, "File uploaded successfully")
	}
}
