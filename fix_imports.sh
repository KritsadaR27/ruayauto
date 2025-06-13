#!/bin/bash

# Fix all import paths from ruaymanagement to correct ruaychatbot paths

echo "🔧 Fixing import paths throughout the project..."

# Find all .go files and replace ruaymanagement imports
find backend -name "*.go" -type f -exec sed -i '' 's|"ruaymanagement/internal/chatbot-core|"ruaychatbot/backend/internal/chatbot-core|g' {} \;
find backend -name "*.go" -type f -exec sed -i '' 's|"ruaymanagement/external/facebook-connect|"ruaychatbot/backend/external/facebook-connect|g' {} \;
find backend -name "*.go" -type f -exec sed -i '' 's|"ruaymanagement/libs/|"ruaychatbot/backend/libs/|g' {} \;

echo "✅ Fixed all import paths!"

# Sync workspace
echo "🔄 Syncing Go workspace..."
go work sync

echo "🎉 All done! Import paths have been corrected."
