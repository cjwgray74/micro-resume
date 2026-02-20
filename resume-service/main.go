package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/jackc/pgx/v5"
)

type ExperienceEntry struct {
	Title       string `json:"title"`
	Company     string `json:"company"`
	Dates       string `json:"dates"`
	Description string `json:"description"`
}

type Resume struct {
	Name       string            `json:"name"`
	Email      string            `json:"email"`
	Phone      string            `json:"phone"`
	Summary    string            `json:"summary"`
	Education  string            `json:"education"`
	Skills     []string          `json:"skills"`
	Experience []ExperienceEntry `json:"experience"`
}

var db *pgx.Conn

func main() {
	var err error

	db, err = pgx.Connect(context.Background(),
		"postgres://postgres:pacman@localhost:5433/resume_db")
	if err != nil {
		log.Fatal("DB connection failed:", err)
	}
	defer db.Close(context.Background())

	http.HandleFunc("/resume", handleResume)

	log.Println("Resume service running on port 8081...")
	http.ListenAndServe(":8081", nil)
}

func handleResume(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		resume, err := loadResume()
		if err != nil {
			http.Error(w, "Failed to load resume", 500)
			return
		}
		json.NewEncoder(w).Encode(resume)

	case "POST":
		var resume Resume
		json.NewDecoder(r.Body).Decode(&resume)

		err := saveResume(resume)
		if err != nil {
			http.Error(w, "Failed to save resume", 500)
			return
		}

		w.Write([]byte(`{"status":"saved"}`))
	}
}

func loadResume() (Resume, error) {
	var data []byte
	err := db.QueryRow(context.Background(),
		"SELECT data FROM resumes WHERE id = 1").Scan(&data)

	if err != nil {
		return Resume{}, nil
	}

	var resume Resume
	json.Unmarshal(data, &resume)
	return resume, nil
}

func saveResume(resume Resume) error {
	data, _ := json.Marshal(resume)

	_, err := db.Exec(context.Background(),
		`INSERT INTO resumes (id, data)
         VALUES (1, $1)
         ON CONFLICT (id)
         DO UPDATE SET data = EXCLUDED.data`,
		data)

	return err
}
