package main

import (
	"io"
	"net/http"
)

func main() {

	// -------------------------------
	// ROUTE: /api/resume → resume-service
	// -------------------------------
	http.HandleFunc("/api/resume", func(w http.ResponseWriter, r *http.Request) {
		// CORS for React
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")

		// Handle preflight request
		if r.Method == http.MethodOptions {
			return
		}

		// Forward request to resume-service
		req, err := http.NewRequest(r.Method, "http://localhost:8081/resume", r.Body)
		if err != nil {
			http.Error(w, "bad gateway", http.StatusBadGateway)
			return
		}
		req.Header = r.Header

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			http.Error(w, "resume service unavailable", http.StatusBadGateway)
			return
		}
		defer resp.Body.Close()

		// Copy headers and status code
		for k, v := range resp.Header {
			for _, vv := range v {
				w.Header().Add(k, vv)
			}
		}
		w.WriteHeader(resp.StatusCode)

		// Copy body
		io.Copy(w, resp.Body)
	})

	// -------------------------------
	// ROUTE: /api/pdf → pdf-service
	// -------------------------------
	http.HandleFunc("/api/pdf", func(w http.ResponseWriter, r *http.Request) {
		// CORS for React
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")

		// Handle preflight request
		if r.Method == http.MethodOptions {
			return
		}

		// Forward request to pdf-service
		req, err := http.NewRequest(r.Method, "http://localhost:8082/pdf", r.Body)
		if err != nil {
			http.Error(w, "bad gateway", http.StatusBadGateway)
			return
		}
		req.Header = r.Header

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			http.Error(w, "pdf service unavailable", http.StatusBadGateway)
			return
		}
		defer resp.Body.Close()

		// Copy headers and status code
		for k, v := range resp.Header {
			for _, vv := range v {
				w.Header().Add(k, vv)
			}
		}
		w.WriteHeader(resp.StatusCode)

		// Copy body (PDF bytes)
		io.Copy(w, resp.Body)
	})

	// -------------------------------
	// START GATEWAY
	// -------------------------------
	http.ListenAndServe(":8080", nil)
}
