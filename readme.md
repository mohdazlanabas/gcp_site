# GCP Site - Modern Go Web Application on App Engine

A beautiful, production-ready web application built with Go 1.22 and deployed to Google Cloud Platform App Engine. Features a modern blue-themed frontend with real-time visitor tracking and elegant animations.

**Live Demo**: https://net1io-web-475001.et.r.appspot.com

![Status](https://img.shields.io/badge/status-live-brightgreen)
![Go Version](https://img.shields.io/badge/go-1.22-blue)
![Platform](https://img.shields.io/badge/platform-GCP%20App%20Engine-blue)
![Region](https://img.shields.io/badge/region-Singapore-orange)

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

## ✨ Features

### Backend
- **Go 1.22 HTTP Server**: Fast, efficient backend with RESTful API
- **Static File Serving**: Serves frontend assets from `src/frontend/`
- **JSON API Endpoint**: `/api/hello` for backend communication
- **Production Ready**: Deployed on GCP App Engine Standard Environment

### Frontend Design
- **Modern Blue Theme**: Professional gradient design with custom color palette
- **Responsive Layout**: Works perfectly on mobile, tablet, and desktop
- **Smooth Animations**: Floating logo, hover effects, and smooth transitions
- **Card-Based UI**: Clean, modern card layout with shadows and effects
- **No Framework**: Pure HTML/CSS/JavaScript for minimal dependencies

### Interactive Features
- **Real-Time Clock**: Live time display in footer, updates every second
- **Visitor Location Tracking**: Automatic geolocation with flag emoji
- **API Testing**: Interactive button to test backend API connection
- **Animated Status**: Pulsing green indicator for live status
- **Server Information**: Displays Singapore hosting location
- **GitHub Integration**: Direct link to source code repository

### Technical Details
- **Instance Class**: F1 (cost-optimized for low traffic)
- **Region**: asia-southeast2 (Singapore)
- **HTTPS**: Automatic SSL/TLS encryption
- **Auto-Scaling**: Scales down when idle to minimize costs
- **One-Command Deployment**: Automated script for git + deployment

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

## 🎨 Frontend Features

The frontend showcases modern web design principles:

### Visual Design
- **Blue Gradient Theme**: Professional color scheme with primary blue (#2563eb)
- **Header**: Animated logo with floating effect and gradient background
- **Cards**: Elevated card design with hover animations
- **Info Dashboard**: 4-card grid showing Region, Runtime, Instance, and Status
- **Footer**: Glassmorphism effect with gradient background

### User Experience
- **Live Clock**: Real-time display showing: `Fri, Nov 12, 2025, 08:37:18 PM`
- **Location Detection**: Uses IP geolocation API (ipapi.co) to show visitor's location
- **Flag Emojis**: Automatically displays country flag based on detected location
- **Smooth Transitions**: All interactive elements have smooth hover effects
- **Responsive Grid**: Adapts to any screen size (mobile, tablet, desktop)

### Technical Implementation
- **CSS Variables**: Modern CSS custom properties for theming
- **Flexbox & Grid**: Modern layout techniques
- **Animations**: Keyframe animations for logo, status indicator, and slide-ins
- **Fetch API**: Modern JavaScript for API calls
- **Geolocation API**: Third-party service integration for location tracking

## 🔌 API Endpoints

### GET /
**Description**: Returns the main HTML page
**Content-Type**: `text/html`

### GET /api/hello
**Description**: Returns a JSON greeting message from the backend

**Response**:
```json
{
  "message": "\n\nHello from the Go backend on App Engine"
}
```

**Example**:
```bash
curl https://net1io-web-475001.et.r.appspot.com/api/hello
```

### GET /static/*
**Description**: Serves static files (CSS, JS) from `src/frontend/`

**Examples**:
- `/static/styles.css` - Main stylesheet
- `/static/scripts.js` - JavaScript functionality

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

## 🎨 Design System

### Color Palette
The application uses a professional blue theme with carefully selected colors:

```css
--primary-blue: #2563eb        /* Main blue */
--primary-blue-dark: #1e40af   /* Darker shade */
--primary-blue-light: #3b82f6  /* Lighter shade */
--bg-primary: #f8fafc          /* Light background */
--bg-secondary: #ffffff        /* White cards */
--text-primary: #1e293b        /* Main text */
--text-secondary: #64748b      /* Secondary text */
```

### Gradients
- **Header/Footer**: `linear-gradient(135deg, #1e3a8a 0%, #1e40af 50%, #2563eb 100%)`
- **Buttons**: `linear-gradient(135deg, #2563eb 0%, #1e40af 100%)`
- **Card Headers**: `linear-gradient(135deg, #dbeafe 0%, #bfdbfe 100%)`

### Typography
- **Font Family**: System fonts (system-ui, -apple-system, Segoe UI, Roboto)
- **Header**: 2rem, bold (700)
- **Body Text**: 1rem, regular (400)
- **Line Height**: 1.6 for readability

### Spacing & Shadows
- **Spacing Scale**: 0.5rem, 1rem, 1.5rem, 2rem, 3rem
- **Border Radius**: 0.375rem to 1rem for soft corners
- **Shadows**: Layered shadows for depth perception

## 💰 Cost Optimization

- **F1 Instance Class**: Minimal compute resources (~$0.05/hour when active)
- **No Always-On Instances**: Scales down when idle to $0
- **Standard Environment**: Faster cold starts than Flexible
- **Minimal Dependencies**: Reduces deployment and runtime overhead
- **CDN-Free**: No external CDN costs, served directly from App Engine
- **Zero Database**: No database costs for this simple application

## Security Considerations

- HTTPS enforced by App Engine automatically
- No authentication implemented (add as needed)
- CORS not configured (add if needed for API)
- Environment variables for sensitive configuration
- IAM roles for service account permissions

## 🛠️ Technologies Used

### Backend
- **Go 1.22**: Primary programming language
- **net/http**: Standard library HTTP server
- **encoding/json**: JSON response handling
- **path/filepath**: File path operations

### Frontend
- **HTML5**: Semantic markup
- **CSS3**: Custom properties (variables), Flexbox, Grid, Animations
- **JavaScript (ES6+)**: Fetch API, async/await, Date API, Geolocation
- **ipapi.co**: Third-party IP geolocation service

### Infrastructure
- **Google Cloud Platform**: Cloud provider
- **App Engine Standard**: Hosting environment
- **Cloud Storage**: File upload staging
- **GitHub**: Version control and code repository

### Development Tools
- **Git**: Version control
- **gcloud CLI**: Deployment and management
- **deploy.sh**: Custom automation script

## 📚 Further Reading

- [App Engine Go Standard Environment](https://cloud.google.com/appengine/docs/standard/go)
- [Go Runtime Configuration](https://cloud.google.com/appengine/docs/standard/go/runtime)
- [App Engine Pricing](https://cloud.google.com/appengine/pricing)
- [Deployment Lessons](reference/lessons.md)
- [GitHub Repository](https://github.com/mohdazlanabas/gcp_site)

## 📝 Project Information

**Repository**: https://github.com/mohdazlanabas/gcp_site
**Live URL**: https://net1io-web-475001.et.r.appspot.com
**Hosting**: Google Cloud App Engine (Singapore)
**Status**: ✅ Production & Active

## 📄 License

© 2025 Net1io - Built with ☕ and deliberately few dependencies.

## 🤝 Support

For deployment issues and lessons learned, see: `reference/lessons.md`

For GCP-related questions, visit: [Google Cloud Community](https://cloud.google.com/community)

## 🌟 Highlights

- ✨ Beautiful blue gradient design
- ⚡ Lightning-fast Go backend
- 🌍 Real-time visitor location tracking
- 🕐 Live clock in footer
- 📱 Fully responsive design
- 🚀 One-command deployment
- 💰 Cost-optimized for minimal expenses
- 🔒 HTTPS by default

---

**Made with 💙 by Roger using Claude Code**
