## DO NOT:
- ❌ Run services directly with `go run` or `npm run dev`
- ❌ Install dependencies locally
- ❌ Use localhost URLs in code (use service names)
- ❌ Create new Docker files - we already have them

### ALWAYS USE:
- ✅ `docker-compose up` to start services
- ✅ `docker-compose logs -f [service]` to see logs
- ✅ `docker exec` to run commands inside containers
- ✅ Service names for internal communication (e.g., `http://sale-service:8084`)

### SERVICES:
- Frontend: http://localhost:3008 (container: frontend)
- Chatbot Service: http://localhost:8090 (container: chatbot-service)
- Webhook Service: http://localhost:8091 (container: webhook-service)
- Database: localhost:5532 (container: postgres)  