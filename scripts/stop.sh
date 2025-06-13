#!/bin/bash

echo "🛑 Stopping ruayChatBot Services..."

# Stop and remove containers
docker-compose down

echo ""
echo "✅ All services stopped!"
echo ""
echo "🧹 To clean up everything:"
echo "   docker-compose down --volumes --remove-orphans"
echo "   docker system prune -f"
