package main

import (
	"io"
	"net/http"
	"strings"
)

func forward(w http.ResponseWriter, r *http.Request, target string) {
	// CORS
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")

	if r.Method == http.MethodOptions {
		return
	}

	req, err := http.NewRequest(r.Method, target, r.Body)
	if err != nil {
		http.Error(w, "bad gateway", http.StatusBadGateway)
		return
	}
	req.Header = r.Header

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		http.Error(w, "service unavailable", http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	for k, v := range resp.Header {
		for _, vv := range v {
			w.Header().Add(k, vv)
		}
	}
	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
}

func main() {

	// -------------------------------
	// NEW ROUTES: Multi‑resume API
	// -------------------------------

	// GET /api/resumes
	// POST /api/resumes
	http.HandleFunc("/api/resumes", func(w http.ResponseWriter, r *http.Request) {
		forward(w, r, "http://localhost:8081/resumes")
	})

	// EVERYTHING under /api/resumes/*
	// GET /api/resumes/{id}
	// PUT /api/resumes/{id}
	// DELETE /api/resumes/{id}
	// POST /api/resumes/{id}/duplicate
	http.HandleFunc("/api/resumes/", func(w http.ResponseWriter, r *http.Request) {
		path := strings.TrimPrefix(r.URL.Path, "/api")
		target := "http://localhost:8081" + path
		forward(w, r, target)
	})

	// -------------------------------
	// ROUTE: /api/pdf → pdf-service
	// -------------------------------
	http.HandleFunc("/api/pdf", func(w http.ResponseWriter, r *http.Request) {
		forward(w, r, "http://localhost:8082/pdf")
	})

	// -------------------------------
	// START GATEWAY
	// -------------------------------
	http.ListenAndServe(":8080", nil)
}
