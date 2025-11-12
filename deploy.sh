#!/bin/bash

# Simple deployment script for GCP App Engine
# Usage: ./deploy.sh "Your commit message"

set -e  # Exit on error

# Check if commit message is provided
if [ -z "$1" ]; then
    echo "❌ Error: Please provide a commit message"
    echo "Usage: ./deploy.sh \"Your commit message\""
    exit 1
fi

COMMIT_MSG="$1"

echo "🔍 Checking git status..."
git status

echo ""
echo "📝 Adding changes to git..."
git add -A

echo ""
echo "💾 Committing changes..."
git commit -m "$COMMIT_MSG

🤖 Generated with [Claude Code](https://claude.com/claude-code)

Co-Authored-By: Claude <noreply@anthropic.com>"

echo ""
echo "⬆️  Pushing to GitHub..."
git push origin main

echo ""
echo "🚀 Deploying to App Engine..."
gcloud app deploy app.yaml --quiet

echo ""
echo "✅ Deployment complete!"
echo ""
echo "🌐 View your app: https://net1io-web-475001.et.r.appspot.com"
echo "📊 View logs: gcloud app logs tail -s default"
echo "🔍 View in browser: gcloud app browse"
