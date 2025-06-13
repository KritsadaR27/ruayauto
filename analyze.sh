#!/bin/zsh

echo "📊 Project Size Analysis"
echo "======================="

echo "📁 Largest directories:"
du -sh */ 2>/dev/null | sort -hr | head -10

echo ""
echo "📄 File count by type:"
echo "Total files: $(find . -type f | wc -l)"
echo "JavaScript/TypeScript: $(find . -name '*.js' -o -name '*.ts' -o -name '*.tsx' -o -name '*.jsx' | wc -l)"
echo "Go files: $(find . -name '*.go' | wc -l)"
echo "Markdown files: $(find . -name '*.md' | wc -l)"
echo "JSON files: $(find . -name '*.json' | wc -l)"

echo ""
echo "🔍 Large files (>1MB):"
find . -type f -size +1M -exec ls -lh {} \; 2>/dev/null | awk '{print $5 " " $9}' | head -10

echo ""
echo "📦 Dependencies size:"
if [ -d "frontend-next/node_modules" ]; then
    echo "node_modules: $(du -sh frontend-next/node_modules 2>/dev/null | cut -f1)"
    echo "node_modules file count: $(find frontend-next/node_modules -type f | wc -l)"
fi

if [ -d "frontend-next/.next" ]; then
    echo ".next cache: $(du -sh frontend-next/.next 2>/dev/null | cut -f1)"
fi

if [ -d "backend/bin" ]; then
    echo "backend/bin: $(du -sh backend/bin 2>/dev/null | cut -f1)"
fi

echo ""
echo "🚨 VS Code Performance Impact:"
echo "Problem directories found:"
find . -name node_modules -o -name .next -o -name .swc -o -name bin | head -10

echo ""
echo "💡 Recommendations:"
if [ $(find . -type f | wc -l) -gt 10000 ]; then
    echo "⚠️  File count is high ($(find . -type f | wc -l)) - VS Code may be slow"
    echo "   Run: ./cleanup.sh to reduce file count"
fi

if [ -d "frontend-next/node_modules" ]; then
    echo "📦 node_modules found - ensure it's in .gitignore and VS Code exclude"
fi

if [ -d "backend/bin" ]; then
    echo "🔧 Go binaries found - ensure backend/bin is in .gitignore"
fi
