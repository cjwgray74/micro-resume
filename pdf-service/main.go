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
	json.NewDecoder(r.Body).Decode(&resume)

	tmpl, _ := template.ParseFiles("templates/resume.html")

	var htmlOutput bytes.Buffer
	tmpl.Execute(&htmlOutput, resume)

	pdf, err := generatePDF(htmlOutput.String())
	if err != nil {
		http.Error(w, "PDF generation failed", 500)
		return
	}

	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition", "attachment; filename=resume.pdf")
	w.Write(pdf)
}

func main() {
	http.HandleFunc("/pdf", pdfHandler)
	log.Println("PDF service running on port 8082...")
	http.ListenAndServe(":8082", nil)
}
