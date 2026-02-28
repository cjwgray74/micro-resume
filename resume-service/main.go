package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

var db *sql.DB

// Resume model
type Resume struct {
	ID        string          `json:"id"`
	Title     string          `json:"title"`
	Data      json.RawMessage `json:"data"`
	Theme     string          `json:"theme"` // theme now fully supported
	CreatedAt time.Time       `json:"created_at"`
	UpdatedAt time.Time       `json:"updated_at"`
}

func main() {
	var err error

	db, err = sql.Open("postgres",
		"postgres://postgres:pacman@localhost:5433/resume_db?sslmode=disable")
	if err != nil {
		log.Fatal("Failed to connect to DB:", err)
	}

	if err = db.Ping(); err != nil {
		log.Fatal("Database unreachable:", err)
	}

	log.Println("Connected to PostgreSQL")

	r := mux.NewRouter()

	// CRUD routes
	r.HandleFunc("/resumes", listResumes).Methods("GET")
	r.HandleFunc("/resumes", createResume).Methods("POST")
	r.HandleFunc("/resumes/{id}", getResume).Methods("GET")
	r.HandleFunc("/resumes/{id}", updateResume).Methods("PUT")
	r.HandleFunc("/resumes/{id}", deleteResume).Methods("DELETE")
	r.HandleFunc("/resumes/{id}/duplicate", duplicateResume).Methods("POST")

	log.Println("Resume service running on :8081")
	http.ListenAndServe(":8081", r)
}

// List all resumes
func listResumes(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	rows, err := db.Query(`
		SELECT id, title, data, theme, created_at, updated_at
		FROM resumes ORDER BY created_at DESC`)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	defer rows.Close()

	var resumes []Resume

	for rows.Next() {
		var res Resume
		if err := rows.Scan(&res.ID, &res.Title, &res.Data, &res.Theme, &res.CreatedAt, &res.UpdatedAt); err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		resumes = append(resumes, res)
	}

	json.NewEncoder(w).Encode(resumes)
}

// Create a new resume
func createResume(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var input struct {
		Title string          `json:"title"`
		Data  json.RawMessage `json:"data"`
		Theme string          `json:"theme"`
	}

	json.NewDecoder(r.Body).Decode(&input)

	// Default theme if none provided
	if input.Theme == "" {
		input.Theme = "default"
	}

	var id string
	err := db.QueryRow(
		`INSERT INTO resumes (title, data, theme)
		 VALUES ($1, $2, $3) RETURNING id`,
		input.Title, input.Data, input.Theme,
	).Scan(&id)

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"id": id})
}

// Get a resume by ID
func getResume(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id := mux.Vars(r)["id"]

	var res Resume
	err := db.QueryRow(
		`SELECT id, title, data, theme, created_at, updated_at
		 FROM resumes WHERE id = $1`,
		id,
	).Scan(&res.ID, &res.Title, &res.Data, &res.Theme, &res.CreatedAt, &res.UpdatedAt)

	if err != nil {
		http.Error(w, "resume not found", 404)
		return
	}

	json.NewEncoder(w).Encode(res)
}

// Update a resume
func updateResume(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id := mux.Vars(r)["id"]

	var input struct {
		Title string          `json:"title"`
		Data  json.RawMessage `json:"data"`
		Theme string          `json:"theme"`
	}

	json.NewDecoder(r.Body).Decode(&input)

	_, err := db.Exec(
		`UPDATE resumes
		 SET title=$1, data=$2, theme=$3, updated_at=NOW()
		 WHERE id=$4`,
		input.Title, input.Data, input.Theme, id,
	)

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.WriteHeader(200)
}

// Delete a resume
func deleteResume(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id := mux.Vars(r)["id"]

	_, err := db.Exec(`DELETE FROM resumes WHERE id=$1`, id)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.WriteHeader(200)
}

// Duplicate a resume
func duplicateResume(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id := mux.Vars(r)["id"]

	var title string
	var data json.RawMessage
	var theme string

	err := db.QueryRow(
		`SELECT title, data, theme FROM resumes WHERE id=$1`,
		id,
	).Scan(&title, &data, &theme)

	if err != nil {
		http.Error(w, "resume not found", 404)
		return
	}

	var newID string
	err = db.QueryRow(
		`INSERT INTO resumes (title, data, theme)
		 VALUES ($1, $2, $3) RETURNING id`,
		title+" (Copy)", data, theme,
	).Scan(&newID)

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"id": newID})
}
