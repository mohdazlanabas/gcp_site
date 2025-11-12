# Deployment Lessons: GCP App Engine Go Application

## Issues Encountered and Solutions

### Issue 1: Deprecated Go Runtime (go120)
**Error:** `Runtime go120 is deprecated and no longer allowed`

**Cause:** The original `app.yaml` specified `runtime: go120`, which has been deprecated by Google Cloud.

**Solution:**
- Updated `app.yaml` to use `runtime: go122`
- Updated `go.mod` to specify `go 1.22`

```yaml
# Before
runtime: go120

# After
runtime: go122
```

---

### Issue 2: First Service Must Be Named "default"
**Error:** `The first service (module) you upload to a new application must be the 'default' service`

**Cause:** App Engine requires the first deployed service to be named "default". The original `app.yaml` specified `service: gcp-serve`.

**Solution:** Changed service name from `gcp-serve` to `default` in `app.yaml`

```yaml
# Before
service: gcp-serve

# After
service: default
```

---

### Issue 3: Missing Staging Bucket
**Error:** `Staging bucket staging.net1io-web-475001.appspot.com is not available`

**Cause:** The App Engine staging bucket didn't exist or was corrupted.

**Solution:** Ran the App Engine repair command to recreate infrastructure:

```bash
gcloud beta app repair
```

---

### Issue 4: Artifact Registry Permissions
**Error:** `Permission "artifactregistry.repositories.uploadArtifacts" denied`

**Cause:** Multiple permission issues with Artifact Registry:
1. Artifact Registry API needed to be enabled
2. Cloud Build service accounts lacked necessary permissions
3. gcr.io redirect repositories didn't exist

**Solution:**

1. **Enabled Artifact Registry API:**
```bash
gcloud services enable artifactregistry.googleapis.com
```

2. **Created gcr.io-compatible repositories:**
```bash
# GCR repository must be in 'us' region
gcloud artifacts repositories create gcr.io \
  --repository-format=docker \
  --location=us \
  --description="GCR redirect repository"

# Asia-specific repository for regional builds
gcloud artifacts repositories create asia.gcr.io \
  --repository-format=docker \
  --location=asia \
  --description="GCR Asia redirect repository"
```

3. **Granted permissions to Cloud Build service accounts:**
```bash
# Grant to Cloud Build default service account
gcloud artifacts repositories add-iam-policy-binding asia.gcr.io \
  --location=asia \
  --member="serviceAccount:591411435095@cloudbuild.gserviceaccount.com" \
  --role="roles/artifactregistry.writer"

# Grant to Cloud Build service agent
gcloud artifacts repositories add-iam-policy-binding asia.gcr.io \
  --location=asia \
  --member="serviceAccount:service-591411435095@gcp-sa-cloudbuild.iam.gserviceaccount.com" \
  --role="roles/artifactregistry.writer"

# Grant to App Engine service account
gcloud artifacts repositories add-iam-policy-binding asia.gcr.io \
  --location=asia \
  --member="serviceAccount:net1io-web-475001@appspot.gserviceaccount.com" \
  --role="roles/artifactregistry.writer"
```

4. **Granted Storage Admin to Cloud Build:**
```bash
gcloud projects add-iam-policy-binding net1io-web-475001 \
  --member="serviceAccount:591411435095@cloudbuild.gserviceaccount.com" \
  --role="roles/storage.admin"
```

---

### Issue 5: Frontend Files Not Deployed
**Error:** `404 page not found` when accessing the root URL

**Cause:** The deployment was initiated from `src/backend/` directory, but `main.go` referenced frontend files at `src/frontend/` using relative paths. Since App Engine only deploys files in the directory containing `app.yaml`, the frontend files were not included in the deployment.

**Solution:**

1. **Copied frontend files into backend directory:**
```bash
cp -r src/frontend src/backend/frontend
```

2. **Updated file paths in `main.go`:**
```go
// Before
fs := http.FileServer(http.Dir(filepath.Join(".", "src", "frontend")))
http.ServeFile(w, r, filepath.Join("src", "frontend", "index.html"))

// After
fs := http.FileServer(http.Dir("frontend"))
http.ServeFile(w, r, filepath.Join("frontend", "index.html"))
```

---

### Issue 6: Static Files 404 Error
**Error:** `404 page not found` when accessing `/static/styles.css` and `/static/scripts.js`

**Cause:** The frontend directory structure had a nested `static` folder (`frontend/static/`), but the Go code expected files directly in the `frontend` directory. The `/static/` URL path was mapped to serve files from the `frontend` directory root.

**Solution:**

1. **Flattened the frontend directory structure:**
```bash
cd src/backend/frontend
mv static/* .
rmdir static
```

2. **Updated HTML file paths to use absolute paths:**
```html
<!-- Before -->
<link rel="stylesheet" href="./static/styles.css">
<script src="./static/scripts.js"></script>

<!-- After -->
<link rel="stylesheet" href="/static/styles.css">
<script src="/static/scripts.js"></script>
```

**Final directory structure:**
```
src/backend/frontend/
├── index.html
├── scripts.js
└── styles.css
```

---

## Final Working Configuration

### Project Structure
```
src/backend/
├── app.yaml          # App Engine configuration
├── .gcloudignore     # Deployment exclusions
├── go.mod            # Go module definition (go 1.22)
├── main.go           # Go server code
└── frontend/         # Frontend files (copied and flattened)
    ├── index.html
    ├── scripts.js
    └── styles.css
```

### app.yaml
```yaml
runtime: go122
service: default
instance_class: F1

env_variables:
  EXAMPLE_ENV: "production"
```

### Deployment Command
```bash
cd src/backend
gcloud app deploy app.yaml --quiet
```

---

## Key Learnings

1. **Always use current runtime versions** - Check GCP documentation for supported runtimes before deployment

2. **First service must be named "default"** - This is a hard requirement for new App Engine applications

3. **Artifact Registry requires proper setup** - When gcr.io has been migrated to Artifact Registry:
   - Create proper redirect repositories in correct regions
   - Grant permissions to all relevant service accounts (Cloud Build SA, service agents, App Engine SA)
   - Allow time for IAM permissions to propagate (5-10 seconds)

4. **File paths matter in deployed environments** - Files referenced in code must be included in the deployment:
   - Only files in the directory containing `app.yaml` (and subdirectories) are deployed
   - Use relative paths from the deployment directory
   - Test file accessibility in the deployed environment
   - Directory structure must match URL routing expectations

5. **Use `gcloud beta app repair`** - This command fixes common App Engine infrastructure issues like missing buckets

6. **Permission propagation takes time** - IAM policy changes may take several seconds to take effect

7. **URL paths and directory structure must align** - When serving static files with `http.StripPrefix`, ensure:
   - The URL path (e.g., `/static/`) matches the stripping prefix
   - The FileServer directory contains the actual files (not subdirectories)
   - HTML references use absolute paths (e.g., `/static/file.css` not `./static/file.css`)

---

## Successful Deployment Result

**Application URL:** https://net1io-web-475001.et.r.appspot.com

**Verified Endpoints:**
- ✅ Frontend: Accessible at root URL (`/`)
- ✅ API: Working at `/api/hello` - returns JSON response
- ✅ Static CSS: Accessible at `/static/styles.css`
- ✅ Static JS: Accessible at `/static/scripts.js`

**Configuration:**
- Runtime: Go 1.22
- Instance class: F1 (cost-optimized)
- Region: asia-southeast2
- Service: default

---

## Testing Commands

```bash
# Test homepage
curl https://net1io-web-475001.et.r.appspot.com/

# Test API
curl https://net1io-web-475001.et.r.appspot.com/api/hello

# Test static assets
curl https://net1io-web-475001.et.r.appspot.com/static/styles.css
curl https://net1io-web-475001.et.r.appspot.com/static/scripts.js

# View in browser
gcloud app browse
```
