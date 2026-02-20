package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
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

func main() {
	http.HandleFunc("/pdf", handlePDF)

	fmt.Println("PDF service running on port 8082...")
	http.ListenAndServe(":8082", nil)
}

func handlePDF(w http.ResponseWriter, r *http.Request) {
	var resume Resume
	json.NewDecoder(r.Body).Decode(&resume)

	html := buildHTML(resume)

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

	if err != nil {
		http.Error(w, "PDF generation failed", 500)
		return
	}

	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition", "attachment; filename=resume.pdf")
	w.Write(pdfBytes)
}

func buildHTML(resume Resume) string {
	skills := strings.Join(resume.Skills, ", ")

	expHTML := ""
	for _, e := range resume.Experience {
		expHTML += fmt.Sprintf(`
            <h3>%s</h3>
            <p><strong>%s</strong></p>
            <p><em>%s</em></p>
            <p>%s</p>
            <hr/>
        `, e.Title, e.Company, e.Dates, e.Description)
	}

	return fmt.Sprintf(`
        <html>
        <body style="font-family: Arial; padding: 40px;">
            <h1>%s</h1>
            <p>%s • %s</p>

            <h2>Summary</h2>
            <p>%s</p>

            <h2>Education</h2>
            <p>%s</p>

            <h2>Skills</h2>
            <p>%s</p>

            <h2>Experience</h2>
            %s
        </body>
        </html>
    `, resume.Name, resume.Email, resume.Phone, resume.Summary, resume.Education, skills, expHTML)
}
