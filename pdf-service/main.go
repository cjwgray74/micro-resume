package main

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"text/template"

	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
)

type Resume struct {
	Name       string `json:"name"`
	Email      string `json:"email"`
	Phone      string `json:"phone"`
	Summary    string `json:"summary"`
	Experience []struct {
		Title       string `json:"title"`
		Company     string `json:"company"`
		Dates       string `json:"dates"`
		Description string `json:"description"`
	} `json:"experience"`
	Education string   `json:"education"`
	Skills    []string `json:"skills"`
}

func generatePDF(html string) ([]byte, error) {
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	var pdfBytes []byte
	err := chromedp.Run(ctx,
		chromedp.Navigate("data:text/html,"+html),
		chromedp.ActionFunc(func(ctx context.Context) error {
			var err error
			pdfBytes, _, err = page.PrintToPDF().Do(ctx)
			return err
		}),
	)

	return pdfBytes, err
}

func pdfHandler(w http.ResponseWriter, r *http.Request) {
	var resume Resume

	if err := json.NewDecoder(r.Body).Decode(&resume); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	tmpl, err := template.ParseFiles("templates/resume.html")
	if err != nil {
		http.Error(w, "Template not found", http.StatusInternalServerError)
		return
	}

	var htmlOutput bytes.Buffer
	if err := tmpl.Execute(&htmlOutput, resume); err != nil {
		http.Error(w, "Template execution failed", http.StatusInternalServerError)
		return
	}

	pdf, err := generatePDF(htmlOutput.String())
	if err != nil {
		http.Error(w, "PDF generation failed", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition", "attachment; filename=resume.pdf")
	w.Write(pdf)
}

func main() {
	http.HandleFunc("/api/pdf", pdfHandler)

	log.Println("PDF service running on port 8080...")
	http.ListenAndServe(":8080", nil)
}
