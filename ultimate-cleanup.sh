#!/bin/bash
# This script performs cleanup and optimization tasks for the ruayChatBot project.

echo "🚀 Ultimate VS Code Performance Cleanup"
echo "======================================"

# Function to show progress
show_progress() {
    echo "🔄 $1..."
}

show_success() {
    echo "✅ $1"
}

# Clean Git repository
show_progress "Cleaning Git repository"
git gc --aggressive --prune=now

# Remove Next.js cache and build artifacts
show_progress "Removing Next.js cache and build artifacts"
rm -rf /Users/kritsadarattanapath/Projects/ruayChatBot/frontend/.next
rm -rf /Users/kritsadarattanapath/Projects/ruayChatBot/frontend/out

# Remove Node.js dependencies
# This is a more drastic step and will require reinstalling dependencies with 'npm install' or 'yarn install'.
# Only uncomment and use if you suspect issues with node_modules or want a completely clean state.
# show_progress "Removing Node.js dependencies"
# rm -rf /Users/kritsadarattanapath/Projects/ruayChatBot/frontend/node_modules

# Remove Go build artifacts
show_progress "Removing Go build artifacts"
rm -rf /Users/kritsadarattanapath/Projects/ruayChatBot/backend/bin
rm -f /Users/kritsadarattanapath/Projects/ruayChatBot/backend/server
rm -rf /Users/kritsadarattanapath/Projects/ruayChatBot/backend/tmp

show_success "Cleanup complete."

# Update file counts
current_files=$(find . -type f | wc -l)
current_size=$(du -sh . | cut -f1)

echo ""
echo "📊 Cleanup Results:"
echo "=================="
echo "📄 Current file count: $current_files"
echo "📁 Current size: $current_size"
echo ""

# Check if node_modules is the problem
if [ -d "frontend-next/node_modules" ]; then
    nm_files=$(find frontend-next/node_modules -type f | wc -l)
    nm_size=$(du -sh frontend-next/node_modules | cut -f1)
    echo "⚠️  node_modules impact:"
    echo "   Files: $nm_files ($(echo "scale=1; $nm_files * 100 / $current_files" | bc)% of total)"
    echo "   Size: $nm_size"
    echo ""
fi

echo "🎯 Performance Impact:"
echo "====================="
if [ $current_files -gt 10000 ]; then
    echo "⚠️  File count still high ($current_files) - consider:"
    echo "   • Move frontend-next to separate workspace"
    echo "   • Use .vscodeignore to exclude more files"
    echo "   • Consider using remote development"
else
    echo "✅ File count is reasonable ($current_files)"
fi

echo ""
echo "🚀 Final Recommendations:"
echo "========================"
echo "1. 🔄 Restart VS Code completely (Cmd+Q)"
echo "2. 📋 File indexing should be much faster"
echo "3. 🔍 Search operations should be quicker"
echo "4. 🤖 Copilot context loading reduced"
echo "5. 💾 Memory usage should be lower"
echo ""
echo "🔧 Monitor performance with:"
echo "   find . -type f | wc -l    # File count"
echo "   du -sh .                  # Total size"
echo "   ./analyze.sh              # Detailed analysis"
