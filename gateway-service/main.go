package main

import (
	"io"
	"net/http"
	"os"
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

	// Read service URLs from environment variables
	resumeServiceURL := os.Getenv("RESUME_SERVICE_URL") // http://resume-service:8081
	pdfServiceURL := os.Getenv("PDF_SERVICE_URL")       // http://pdf-service:8082

	// -------------------------------
	// Multi‑resume API
	// -------------------------------

	// GET /api/resumes
	// POST /api/resumes
	http.HandleFunc("/api/resumes", func(w http.ResponseWriter, r *http.Request) {
		forward(w, r, resumeServiceURL+"/resumes")
	})

	// EVERYTHING under /api/resumes/*
	http.HandleFunc("/api/resumes/", func(w http.ResponseWriter, r *http.Request) {
		path := strings.TrimPrefix(r.URL.Path, "/api")
		target := resumeServiceURL + path
		forward(w, r, target)
	})

	// -------------------------------
	// PDF route
	// -------------------------------
	http.HandleFunc("/api/pdf", func(w http.ResponseWriter, r *http.Request) {
		forward(w, r, pdfServiceURL+"/api/pdf")
	})

	// -------------------------------
	// START GATEWAY
	// -------------------------------
	http.ListenAndServe(":8080", nil)
}
