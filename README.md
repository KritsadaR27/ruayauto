# 🤖 Ruay ChatBot - Facebook Auto Reply System

## 🚀 Quick Start

### 1. Environment Setup
```bash
# Copy environment template
cp .env.example .env

# Edit .env file with your actual values
nano .env
```

### 2. Required Environment Variables

**⚠️ IMPORTANT: Never commit real API keys to git!**

```bash
# Facebook Configuration (Required)
FACEBOOK_APP_ID=your_facebook_app_id_here
NEXT_PUBLIC_FACEBOOK_APP_ID=your_facebook_app_id_here
FACEBOOK_APP_SECRET=your_facebook_app_secret_here
FACEBOOK_VERIFY_TOKEN=your_custom_verify_token
FACEBOOK_PAGE_ACCESS_TOKEN=your_page_access_token

# Database Configuration
DB_NAME=ruay_chatbot
DB_USER=ruay_user
DB_PASSWORD=your_secure_password
DB_HOST=postgres
DB_PORT=5432

# Security
NEXTAUTH_SECRET=your_nextauth_secret_here
JWT_SECRET=your_jwt_secret_here
```

### 3. Start Services

```bash
# Start all services
docker-compose up -d

# Check logs
docker-compose logs -f web
docker-compose logs -f chatbot
docker-compose logs -f webhook
```

### 4. Access Points

- **Web Interface**: http://localhost:3008
- **Chatbot API**: http://localhost:8090
- **Webhook API**: http://localhost:8091
- **Database**: localhost:5532

## 🔒 Security Notes

1. **Never commit `.env` file to git**
2. **Use strong passwords and secrets**
3. **Rotate API keys regularly**
4. **Use HTTPS in production**

## 🛠️ Development

```bash
# View logs
docker-compose logs -f [service_name]

# Restart specific service
docker-compose restart [service_name]

# Execute commands inside container
docker exec -it ruaychatbot_web_1 sh
docker exec -it ruaychatbot_chatbot_1 sh

# Stop all services
docker-compose down
```

## 📋 Service Configuration

### Chatbot Service (Port 8090)
- Handles business logic
- Database connections
- Facebook API integration

### Webhook Service (Port 8091)  
- Receives Facebook webhooks
- Processes incoming messages
- Forwards to chatbot service

### Web Service (Port 3008)
- Next.js frontend
- User interface
- Facebook OAuth integration

## 🐛 Troubleshooting

### Common Issues

1. **"App ID is incorrect" error**
   - Check `NEXT_PUBLIC_FACEBOOK_APP_ID` is set correctly
   - Verify Facebook App configuration

2. **Database connection issues**
   - Check database credentials in `.env`
   - Ensure postgres service is healthy

3. **Service not starting**
   - Check Docker logs: `docker-compose logs [service]`
   - Verify all required environment variables are set

### Debug Commands

```bash
# Check running containers
docker-compose ps

# Check service health
docker-compose exec chatbot wget -qO- http://localhost:8090/health
docker-compose exec webhook wget -qO- http://localhost:8091/health

# Reset everything
docker-compose down -v
docker-compose up --build
```
