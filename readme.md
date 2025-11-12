# Simple App Engine app (no Docker)

Structure:
- src/frontend: index.html, styles, scripts
- src/backend: main.go

Local run:
1. cd src/backend
2. go run main.go
3. Open http://localhost:8080

Deploy to App Engine:
1. gcloud init
2. gcloud app create --region=YOUR_REGION   # first time only
3. From project root: gcloud app deploy app.yaml