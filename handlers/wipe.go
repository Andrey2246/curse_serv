package handlers

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	_ "github.com/lib/pq"
)

func WipeHandler(storageDir string, passphrase string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
			return
		}

		// Optional passphrase for additional protection
		if passphrase != "" {
			providedPassphrase := r.FormValue("passphrase")
			if providedPassphrase != passphrase {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}
		}

		// Delete all files in the storage directory
		files, err := ioutil.ReadDir(storageDir)
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
		db, err := sql.Open("postgres", "your_connection_string_here")
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

		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "All data wiped successfully")
	}
}
