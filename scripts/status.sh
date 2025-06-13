#!/bin/bash

echo "📊 ruayChatBot Services Status"
echo "==============================="

# Check if containers are running
docker-compose ps

echo ""
echo "🔍 Service Health Checks:"
echo "========================="

# ChatBotCore health
echo -n "🤖 ChatBotCore:     "
if curl -s http://localhost:8090/health >/dev/null 2>&1; then
    echo "✅ Healthy"
    curl -s http://localhost:8090/health | jq '.status'
else
    echo "❌ Not responding"
fi

# FacebookConnect health  
echo -n "📘 FacebookConnect: "
if curl -s http://localhost:8091/health >/dev/null 2>&1; then
    echo "✅ Healthy"
    curl -s http://localhost:8091/health | jq '.status'
else
    echo "❌ Not responding"
fi

echo ""
echo "📋 Useful Commands:"
echo "=================="
echo "📊 View logs:        docker-compose logs -f"
echo "📊 View specific:    docker-compose logs -f chatbot-core"
echo "🔄 Restart:          docker-compose restart"
echo "🛑 Stop:             docker-compose down"
