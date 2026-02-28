package main

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"text/template"

	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
)

// Incoming JSON shape from gateway:
//
//	{
//	  "title": "...",
//	  "data": { ...resume fields... },
//	  "theme": "modern"
//	}
type ResumeEnvelope struct {
	Title string          `json:"title"`
	Data  json.RawMessage `json:"data"`
	Theme string          `json:"theme"`
}

// This matches your HTML templates
type Resume struct {
	Name       string `json:"name"`
	Email      string `json:"email"`
	Phone      string `json:"phone"`
	Summary    string `json:"summary"`
	Education  string `json:"education"`
	Experience []struct {
		Title       string `json:"title"`
		Company     string `json:"company"`
		Dates       string `json:"dates"`
		Description string `json:"description"`
	} `json:"experience"`
	Skills []string `json:"skills"`
}

func generatePDF(html string) ([]byte, error) {

	allocOpts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.ExecPath(`C:\Program Files\Google\Chrome\Application\chrome.exe`),
		chromedp.NoSandbox,
	)

	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), allocOpts...)
	defer cancel()

	ctx, cancel := chromedp.NewContext(allocCtx)
	defer cancel()

	var pdfBytes []byte

	escaped := url.PathEscape(html)

	err := chromedp.Run(ctx,
		chromedp.Navigate("data:text/html,"+escaped),
		chromedp.ActionFunc(func(ctx context.Context) error {
			var err error
			pdfBytes, _, err = page.PrintToPDF().
				WithMarginTop(0.4).
				WithMarginBottom(0.4).
				WithMarginLeft(0.4).
				WithMarginRight(0.4).
				WithPreferCSSPageSize(true).
				Do(ctx)
			return err
		}),
	)

	return pdfBytes, err
}

func pdfHandler(w http.ResponseWriter, r *http.Request) {

	// 1. Decode outer envelope
	var env ResumeEnvelope
	if err := json.NewDecoder(r.Body).Decode(&env); err != nil {
		log.Println("JSON decode error:", err)
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// 2. Decode inner resume.data into Resume struct
	var resume Resume
	if err := json.Unmarshal(env.Data, &resume); err != nil {
		log.Println("Inner resume decode error:", err)
		http.Error(w, "Invalid resume data", http.StatusBadRequest)
		return
	}

	// 3. Select template based on theme
	templateName := "resume_default.html"
	if env.Theme != "" {
		templateName = "resume_" + env.Theme + ".html"
	}

	tmpl, err := template.ParseFiles("templates/" + templateName)
	if err != nil {
		log.Println("Template error:", err)
		http.Error(w, "Template not found", http.StatusInternalServerError)
		return
	}

	// 4. Render HTML
	var htmlOutput bytes.Buffer
	if err := tmpl.Execute(&htmlOutput, resume); err != nil {
		log.Println("Template execution error:", err)
		http.Error(w, "Template execution failed", http.StatusInternalServerError)
		return
	}

	// 5. Convert HTML → PDF
	pdf, err := generatePDF(htmlOutput.String())
	if err != nil {
		log.Println("PDF generation error:", err)
		http.Error(w, "PDF generation failed", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition", "attachment; filename=resume.pdf")
	w.Write(pdf)
}

func main() {
	http.HandleFunc("/api/pdf", pdfHandler)

	log.Println("PDF service running on port 8082...")
	http.ListenAndServe(":8082", nil)
}
