package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"

	_ "github.com/lib/pq"
)

type FileInfo struct {
	Filename  string `json:"filename"`
	Timestamp string `json:"timestamp"`
	Uploader  string `json:"uploader"`
	Size      int64  `json:"size"`
}

func ListHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	db, err := sql.Open("postgres", "your_connection_string_here")
	if err != nil {
		http.Error(w, "Failed to connect to the database", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	rows, err := db.Query("SELECT filename, timestamp, uploader, size FROM files")
	if err != nil {
		http.Error(w, "Failed to retrieve files", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var files []FileInfo
	for rows.Next() {
		var file FileInfo
		err := rows.Scan(&file.Filename, &file.Timestamp, &file.Uploader, &file.Size)
		if err != nil {
			http.Error(w, "Failed to parse files", http.StatusInternalServerError)
			return
		}
		files = append(files, file)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(files)
}
