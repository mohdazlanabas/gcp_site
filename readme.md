# GCP Site - Simple Go Web Application on App Engine

A minimal, production-ready web application built with Go and deployed to Google Cloud Platform App Engine.

## Project Structure

```
gcp_site/
├── app.yaml              # App Engine configuration
├── go.mod                # Go module definition
├── main.go               # Go HTTP server and routing
├── deploy.sh             # Deployment automation script
├── cloudbuild.yaml       # Cloud Build config (archived)
├── src/
│   └── frontend/         # Static web assets
│       ├── index.html    # Main HTML page
│       ├── scripts.js    # Client-side JavaScript
│       └── styles.css    # Styling
├── reference/            # Documentation
│   ├── lessons.md        # Deployment lessons learned
│   └── readme.pdf        # PDF documentation
└── readme.md             # This file
```

## Features

- **Backend**: Go 1.22 HTTP server with RESTful API
- **Frontend**: Pure HTML/CSS/JavaScript (no frameworks)
- **Deployment**: Google Cloud Platform App Engine Standard Environment
- **Instance Class**: F1 (cost-optimized)
- **Region**: asia-southeast2

## Prerequisites

- [Go 1.22+](https://golang.org/dl/)
- [Google Cloud SDK](https://cloud.google.com/sdk/docs/install)
- GCP account with billing enabled
- Git

## Local Development

### 1. Run the Backend Server

From the project root:

```bash
go run main.go
```

The server will start on `http://localhost:8080`

### 2. Access the Application

Open your browser and navigate to:
- **Homepage**: http://localhost:8080
- **API Endpoint**: http://localhost:8080/api/hello
- **Static Assets**: http://localhost:8080/static/styles.css

### 3. Make Changes

- Backend: Edit `main.go`
- Frontend: Edit files in `src/frontend/`
- Restart the server to see backend changes
- Refresh the browser to see frontend changes

## Deployment to Google Cloud App Engine

### First-Time Setup

1. **Initialize gcloud CLI**:
```bash
gcloud init
```

2. **Set your project**:
```bash
gcloud config set project YOUR_PROJECT_ID
```

3. **Create App Engine application** (first time only):
```bash
gcloud app create --region=asia-southeast2
```

### Deploy Application

**Option 1: Automated Script (Recommended)**

Use the included `deploy.sh` script for one-command deployment:

```bash
./deploy.sh "Your commit message here"
```

The script automatically:
1. Shows current git status
2. Adds all changes to git staging
3. Creates a commit with your message
4. Pushes to GitHub (main branch)
5. Deploys to App Engine
6. Shows deployment success with app URL

**Example:**
```bash
./deploy.sh "Add new feature to user dashboard"
```

**Option 2: Manual Step-by-Step**

If you prefer manual control:

```bash
# 1. Check what files changed
git status

# 2. Stage your changes
git add -A

# 3. Commit with a message
git commit -m "Your descriptive message"

# 4. Push to GitHub
git push origin main

# 5. Deploy to App Engine
gcloud app deploy app.yaml --quiet
```

### Deployment Workflow Notes

- **No automatic CI/CD**: The project uses manual deployment workflow
- **GitHub**: Used for version control only
- **Cloud Build**: Currently disabled (see `cloudbuild.yaml` for details)
- **Deployment Time**: Typically 2-3 minutes
- **Current URL**: https://net1io-web-475001.et.r.appspot.com

### View Deployed Application

```bash
gcloud app browse
```

Or visit: `https://YOUR_PROJECT_ID.et.r.appspot.com`

### View Logs

```bash
gcloud app logs tail -s default
```

## API Endpoints

### GET /
Returns the main HTML page

### GET /api/hello
Returns a JSON greeting message

**Response**:
```json
{
  "message": "Hello from the Go backend on App Engine"
}
```

### GET /static/*
Serves static files (CSS, JS) from `src/frontend/`

## Configuration

### app.yaml

The App Engine configuration file (in project root):

```yaml
runtime: go122              # Go 1.22 runtime
service: default            # Default service
instance_class: F1          # Smallest instance (cost-effective)

env_variables:
  EXAMPLE_ENV: "production" # Environment variables
```

The `main.go` file is located in the project root and serves static files from `src/frontend/`.

### go.mod

Go module configuration:

```go
module example.com/my-app
go 1.22
```

## Development Tips

### Running Tests

```bash
go test ./...
```

### Build Binary Locally

```bash
go build -o app
./app
```

### Check for Updates

```bash
go get -u ./...
go mod tidy
```

## Deployment Troubleshooting

If you encounter deployment issues, refer to `reference/lessons.md` which documents common problems and solutions including:

- Deprecated runtime versions
- Artifact Registry permissions
- Service naming requirements
- File path configurations
- Static file serving issues

## Cost Optimization

- **F1 Instance Class**: Minimal compute resources
- **No Always-On Instances**: Scales down when idle
- **Standard Environment**: Faster cold starts than Flexible
- **Minimal Dependencies**: Reduces deployment and runtime overhead

## Security Considerations

- HTTPS enforced by App Engine automatically
- No authentication implemented (add as needed)
- CORS not configured (add if needed for API)
- Environment variables for sensitive configuration
- IAM roles for service account permissions

## Further Reading

- [App Engine Go Standard Environment](https://cloud.google.com/appengine/docs/standard/go)
- [Go Runtime Configuration](https://cloud.google.com/appengine/docs/standard/go/runtime)
- [App Engine Pricing](https://cloud.google.com/appengine/pricing)
- [Deployment Lessons](reference/lessons.md)

## License

© 2025 Roger - Built with ☕ and deliberately few dependencies.

## Support

For deployment issues and lessons learned, see: `reference/lessons.md`

For GCP-related questions, visit: [Google Cloud Community](https://cloud.google.com/community)
