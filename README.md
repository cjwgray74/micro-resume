📄 Resume Builder — Microservices Architecture
https://img.shields.io/badge/Backend-Go-00ADD8?logo=go&logoColor=white https://img.shields.io/badge/Frontend-React-61DAFB?logo=react&logoColor=black https://img.shields.io/badge/Architecture-Microservices-blue https://img.shields.io/badge/Status-In_Development-yellow https://img.shields.io/badge/License-MIT-green
A microservices‑based Resume Builder application demonstrating real‑world engineering skills across frontend development, backend services, API gateway routing, and PDF generation. Users can create, edit, duplicate, rename, view, and export resumes as downloadable PDFs.

🏗️ System Architecture
                   ┌──────────────────────────┐
                   │        Frontend          │
                   │          React           │
                   │   localhost:3000         │
                   └─────────────┬────────────┘
                                 │
                                 ▼
                   ┌──────────────────────────┐
                   │       API Gateway        │
                   │            Go            │
                   │   localhost:8080         │
                   └──────┬──────────┬────────┘
                          │          │
        ┌─────────────────┘          └──────────────────┐
        ▼                                                ▼
┌──────────────────────┐                      ┌────────────────────────┐
│   Resume Service      │                      │     PDF Service        │
│          Go           │                      │   Go + Chromedp        │
│   localhost:8081      │                      │   localhost:8082       │
└──────────────────────┘                      └────────────────────────┘



🚀 Features
Frontend (React)
- Resume list with grid layout
- Create resume
- Edit resume
- Delete resume
- Duplicate resume
- Rename resume
- View resume (full preview)
- Download PDF button
Resume Service (Go)
- Full CRUD operations
- JSON-based resume model
- Supports themes
- Stable API responses
API Gateway (Go)
- Central entry point for all frontend requests
- Forwards:
- /api/resumes → Resume Service
- /api/pdf → PDF Service
- Handles routing and service communication
PDF Service (Go + Chromedp)
- Converts resume JSON → HTML → PDF
- Uses templates/resume.html
- Headless Chrome rendering
- PDF generation pipeline implemented

🛠️ Work in Progress
- Correcting gateway routing for PDF (/api/pdf should not become /api/pdf/:id)
- Aligning ports between gateway and PDF service
- Adding Chrome executable path for Chromedp on Windows
- Improving PDF layout and adding themes

📁 Project Structure
project-root/
│
├── frontend/               # React UI
│     ├── src/
│     └── package.json
│
├── gateway-service/        # API Gateway
│     └── main.go
│
├── resume-service/         # Resume CRUD microservice
│     └── main.go
│
├── pdf-service/            # PDF generator microservice
│     ├── main.go
│     └── templates/
│           └── resume.html
│
└── README.md



🧰 Tech Stack
Frontend
- React
- React Router
- Fetch API
Backend
- Go
- Chromedp (PDF generation)
- net/http
- JSON APIs
Architecture
- Microservices
- API Gateway pattern
- Service isolation
- JSON-based communication

🧪 Running the Project
1. Start the Resume Service
cd resume-service
go run main.go


2. Start the PDF Service
cd pdf-service
go run main.go


3. Start the API Gateway
cd gateway-service
go run main.go


4. Start the Frontend
cd frontend
npm install
npm start



⚙️ Requirements
- Go 1.20+
- Node.js 18+
- Chrome or Edge installed (required for Chromedp)
- Windows, macOS, or Linux

📌 Next Steps
- Fix PDF routing through the gateway
- Add multiple resume themes
- Add Docker support for all services
- Add PostgreSQL persistence
- Deploy to AWS or Render

📜 License
MIT License — free to use, modify, and distribute.


