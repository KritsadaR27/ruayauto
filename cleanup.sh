#!/bin/zsh

echo "🧹 Starting VS Code Performance Cleanup..."

# Clean Next.js build cache
echo "📦 Cleaning Next.js cache..."
if [ -d "frontend-next/.next" ]; then
    rm -rf frontend-next/.next
    echo "✅ Removed .next directory"
fi

if [ -d "frontend-next/.swc" ]; then
    rm -rf frontend-next/.swc
    echo "✅ Removed .swc directory"
fi

if [ -d "frontend-next/out" ]; then
    rm -rf frontend-next/out
    echo "✅ Removed out directory"
fi

# Clean Go binaries
echo "🔧 Cleaning Go binaries..."
if [ -d "backend/bin" ]; then
    rm -rf backend/bin
    echo "✅ Removed backend/bin directory"
fi

if [ -f "backend/server" ]; then
    rm -f backend/server
    echo "✅ Removed backend/server binary"
fi

if [ -f "backend/migrate" ]; then
    rm -f backend/migrate
    echo "✅ Removed backend/migrate binary"
fi

# Clean temporary files
echo "🗑️  Cleaning temporary files..."
find . -name "*.log" -type f -delete 2>/dev/null || true
find . -name "*.tmp" -type f -delete 2>/dev/null || true
find . -name ".DS_Store" -type f -delete 2>/dev/null || true
find . -name "Thumbs.db" -type f -delete 2>/dev/null || true

# Git cleanup
echo "🔄 Git cleanup..."
git gc --aggressive --prune=now

# VS Code workspace cleanup
echo "💻 VS Code cleanup..."
if [ -d ".vscode" ]; then
    # Keep only settings.json, remove other cache
    find .vscode -name "*.log" -delete 2>/dev/null || true
fi

# Clean npm cache if exists
if command -v npm &> /dev/null; then
    echo "📚 Cleaning npm cache..."
    npm cache clean --force 2>/dev/null || true
fi

echo ""
echo "✨ Cleanup completed!"
echo ""
echo "📊 Current project size:"
du -sh . 2>/dev/null
echo ""
echo "📄 File count: $(find . -type f | wc -l)"
echo ""
echo "🚀 Recommendations:"
echo "1. Restart VS Code completely"
echo "2. If still slow, run: cd frontend-next && npm ci"
echo "3. Monitor file count with: find . -type f | wc -l"
echo "4. Check largest folders: du -sh */ | sort -hr"
