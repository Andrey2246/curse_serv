package handlers

import (
	"curse_serv/logger"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/lib/pq"
)

func WipeHandler(storageDir string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
			return
		}

		// Delete all files in the storage directory
		files, err := os.ReadDir(storageDir)
		if err != nil {
			http.Error(w, "Failed to read storage directory", http.StatusInternalServerError)
			return
		}

		for _, file := range files {
			err := os.Remove(fmt.Sprintf("%s/%s", storageDir, file.Name()))
			if err != nil {
				http.Error(w, "Failed to delete files", http.StatusInternalServerError)
				return
			}
		}

		// Clear the database
		db, err := sql.Open("postgres", "postgresql://cursework:security@localhost:5432/")
		if err != nil {
			http.Error(w, "Failed to connect to the database", http.StatusInternalServerError)
			return
		}
		defer db.Close()

		_, err = db.Exec("DELETE FROM files")
		if err != nil {
			http.Error(w, "Failed to clear database", http.StatusInternalServerError)
			return
		}
		log.Println("Database cleared")

		if err := logger.Log("Wipe", "All files wiped from server"); err != nil {
			log.Printf("Failed to log wipe action: %v", err)
		}

		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "All data wiped successfully")
	}
}
