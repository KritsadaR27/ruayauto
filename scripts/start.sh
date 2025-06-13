#!/bin/bash

echo "🚀 Starting ruayChatBot Services..."

# Build and start services
docker-compose up --build -d

echo ""
echo "✅ Services started successfully!"
echo ""
echo "📋 Service URLs:"
echo "   🤖 ChatBotCore:     http://localhost:8090"
echo "   📘 FacebookConnect: http://localhost:8091"
echo ""
echo "🔍 Health Checks:"
echo "   🤖 ChatBotCore:     http://localhost:8090/health"
echo "   📘 FacebookConnect: http://localhost:8091/health"
echo ""
echo "📊 View logs:   docker-compose logs -f"
echo "🛑 Stop:        docker-compose down"
echo ""
echo "🔄 Waiting for services to be ready..."
sleep 5

# Check service health
echo "🔍 Checking service health..."
curl -s http://localhost:8090/health | jq '.' || echo "❌ ChatBotCore not ready yet"
curl -s http://localhost:8091/health | jq '.' || echo "❌ FacebookConnect not ready yet"
