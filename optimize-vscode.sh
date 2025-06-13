#!/bin/zsh

echo "⚡ VS Code Performance Optimizer"
echo "================================"

# Create .vscode directory if not exists
mkdir -p .vscode

echo "📝 Creating optimized VS Code settings..."

# Backup existing settings
if [ -f ".vscode/settings.json" ]; then
    cp .vscode/settings.json .vscode/settings.json.backup
    echo "✅ Backed up existing settings to settings.json.backup"
fi

# Write optimized settings
cat > .vscode/settings.json << 'EOF'
{
  "files.exclude": {
    "**/node_modules": true,
    "**/.next": true,
    "**/.swc": true,
    "**/bin": true,
    "**/server": true,
    "**/migrate": true,
    "**/tmp": true,
    "**/*.log": true,
    "**/.DS_Store": true,
    "**/Thumbs.db": true,
    "**/package-lock.json": true,
    "**/go.sum": true,
    "backend/internal/*/server": true,
    "backend/external/*/server": true
  },
  "files.watcherExclude": {
    "**/node_modules/**": true,
    "**/.next/**": true,
    "**/.swc/**": true,
    "**/bin/**": true,
    "**/server/**": true,
    "**/migrate/**": true,
    "**/tmp/**": true,
    "**/*.log": true,
    "backend/internal/**/server": true,
    "backend/external/**/server": true
  },
  "search.exclude": {
    "**/node_modules": true,
    "**/.next": true,
    "**/bin": true,
    "**/server": true,
    "**/migrate": true,
    "**/tmp": true,
    "**/*.log": true,
    "**/package-lock.json": true,
    "**/go.sum": true
  },
  "typescript.preferences.excludePackageJsonAutoImports": "on",
  "typescript.suggest.autoImports": false,
  "go.toolsManagement.autoUpdate": false,
  "files.associations": {
    "*.go": "go"
  },
  "editor.formatOnSave": true,
  "go.formatTool": "goimports",
  "files.maxMemoryForLargeFilesMB": 4096,
  "search.smartCase": true,
  "search.useGlobalIgnoreFiles": true,
  "typescript.disableAutomaticTypeAcquisition": true,
  "extensions.autoUpdate": false,
  "git.autorefresh": false,
  "git.scanRepositories": [],
  "postman.settings.dotenv-detection-notification-visibility": false
}
EOF

echo "✅ Created optimized VS Code settings"

# Create extensions.json for recommended extensions
cat > .vscode/extensions.json << 'EOF'
{
  "recommendations": [
    "golang.go",
    "bradlc.vscode-tailwindcss",
    "ms-vscode.vscode-typescript-next"
  ],
  "unwantedRecommendations": [
    "ms-vscode.node-debug2",
    "ms-vscode.references-view"
  ]
}
EOF

echo "✅ Created extensions recommendations"

# Update .copilotignore to be more aggressive
cat > .copilotignore << 'EOF'
# Dependencies - Completely ignore
node_modules/
frontend-next/node_modules/
**/node_modules/

# Build outputs - Completely ignore
frontend-next/.next/
frontend-next/out/
frontend-next/.swc/
backend/bin/
backend/server
backend/migrate
backend/tmp/
backend/internal/*/server
backend/external/*/server

# Documentation - Reduce context (keep only essential)
*.md
!README.md
!MVP_SPECIFICATION.md

# Configuration - Ignore most configs
docker-compose*.yml
Dockerfile*
*.sh
*.sql
nginx.conf
*.yaml
*.yml

# Tests - Reduce testing context
*_test.go
*.test.js
test/
tests/

# Logs, databases, temp files
*.log
*.db
*.sqlite*
*.tmp

# Environment and secrets
.env*

# Large generated files
package-lock.json
go.sum
go.mod.new

# Git and IDE
.git/
.vscode/
.idea/

# OS files
.DS_Store
Thumbs.db
EOF

echo "✅ Updated .copilotignore for better performance"

echo ""
echo "🎯 Performance optimizations applied:"
echo "   ✓ Excluded node_modules from VS Code indexing"
echo "   ✓ Disabled file watching for large directories"
echo "   ✓ Limited TypeScript auto-imports"
echo "   ✓ Reduced Copilot context loading"
echo "   ✓ Disabled unnecessary auto-updates"
echo ""
echo "🚀 Next steps:"
echo "   1. Restart VS Code completely (Cmd+Q then reopen)"
echo "   2. File count should feel much faster now"
echo "   3. Monitor with: find . -type f | wc -l"
echo ""
echo "💡 If still slow, consider:"
echo "   - Move frontend-next to separate workspace"
echo "   - Use VS Code workspace file for multi-root setup"
