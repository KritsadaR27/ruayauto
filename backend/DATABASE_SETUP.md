# Database Setup Guide

This guide explains how to set up and use the PostgreSQL database for the ruayAutoMsg project.

## Quick Start

### 1. Prerequisites

- PostgreSQL 16+ installed locally OR
- Docker and Docker Compose (recommended)

### 2. Using Docker (Recommended)

The easiest way to run the application with PostgreSQL:

```bash
# Copy environment variables
cp .env.example .env

# Edit .env file with your Facebook credentials
nano .env

# Start everything with Docker Compose
docker-compose up
```

This will:
- Start PostgreSQL on port 55432
- Start the backend on port 3006
- Start the frontend on port 5173
- Automatically run database migrations

### 3. Manual PostgreSQL Setup

If you prefer to run PostgreSQL manually:

```bash
# Install PostgreSQL (macOS)
brew install postgresql@16

# Start PostgreSQL service
brew services start postgresql@16

# Create database and user
psql postgres
CREATE DATABASE "ruayAutoMsg";
CREATE USER ruay WITH PASSWORD 'ruay1234';
GRANT ALL PRIVILEGES ON DATABASE "ruayAutoMsg" TO ruay;
\q

# Update .env file
cp .env.example .env
# Edit DATABASE_URL in .env to match your setup
```

## Environment Configuration

### Required Configuration

1. **Copy the example environment file:**
   ```bash
   cp .env.example .env
   ```

2. **Configure database connection in `.env`:**
   ```bash
   # For Docker setup (default)
   DATABASE_URL=postgres://ruay:ruay1234@localhost:55432/ruayAutoMsg?sslmode=disable

   # For local PostgreSQL setup
   DATABASE_URL=postgres://ruay:ruay1234@localhost:5432/ruayAutoMsg?sslmode=disable
   ```

3. **Configure Facebook credentials (REQUIRED):**
   ```bash
   FACEBOOK_VERIFY_TOKEN=your_custom_verify_token_here
   FACEBOOK_PAGE_ACCESS_TOKEN=your_facebook_page_access_token_here
   ```

### Facebook Setup

To get your Facebook credentials:

1. Go to [Facebook Developers](https://developers.facebook.com/)
2. Create a new app or use an existing one
3. Add the "Messenger" product
4. Generate a Page Access Token for your Facebook page
5. Set up a custom verify token (any string you choose)
6. Configure your webhook URL: `https://your-domain.com/webhook/facebook`

## Database Operations

### Running Migrations

Migrations run automatically when you start the server. You can also run them manually:

```bash
# Build the migration tool
go build -o bin/migrate ./cmd/migrate

# Run all pending migrations
./bin/migrate -up

# Check migration status
./bin/migrate -status

# Rollback last migration
./bin/migrate -down

# Create a new migration
./bin/migrate -create add_new_feature
```

### Database Schema

The database includes these tables:

1. **keywords** - Stores keyword-response mappings for auto-reply
2. **message_analytics** - Tracks message statistics and performance
3. **webhook_logs** - Logs all webhook requests for debugging

### Testing Database Connection

```bash
# Build and run the server
go build -o bin/server ./cmd/server
./bin/server

# Check health endpoint
curl http://localhost:3006/health
```

## Development Workflow

### 1. Local Development

```bash
# Start PostgreSQL (if using Docker)
docker-compose up db

# Run the backend
go run ./cmd/server/main.go

# The server will:
# - Connect to PostgreSQL
# - Run pending migrations automatically
# - Start health monitoring
# - Listen on port 3006
```

### 2. Adding New Migrations

```bash
# Create a new migration file
go run ./cmd/migrate/main.go -create add_user_preferences

# This creates:
# internal/migrations/000003_add_user_preferences.up.sql
# internal/migrations/000003_add_user_preferences.down.sql

# Edit the .up.sql file with your schema changes
# Edit the .down.sql file with rollback instructions

# Test the migration
go run ./cmd/migrate/main.go -up
```

### 3. Database Connection Pooling

The application uses connection pooling with these settings:

- **Max Connections**: 25
- **Min Connections**: 5
- **Health Check Interval**: 1 minute
- **Connection Timeout**: 30 seconds

You can adjust these in `internal/database/database.go`.

## Production Deployment

### Environment Variables

Ensure these are set in production:

```bash
# Database
DATABASE_URL=postgres://user:password@host:port/dbname?sslmode=require

# Facebook (REQUIRED)
FACEBOOK_VERIFY_TOKEN=your_production_verify_token
FACEBOOK_PAGE_ACCESS_TOKEN=your_production_page_token

# Security
GIN_MODE=release
API_KEY=your_secure_api_key
JWT_SECRET=your_secure_jwt_secret

# Performance
MAX_CONNECTIONS=50
CACHE_TTL=600
```

### SSL Configuration

For production, enable SSL:

```bash
DATABASE_URL=postgres://user:password@host:port/dbname?sslmode=require
```

### Monitoring

The application provides health monitoring:

- Health endpoint: `GET /health`
- Database statistics in health response
- Automatic reconnection on connection loss
- Graceful shutdown on application termination

## Troubleshooting

### Common Issues

1. **Connection refused**
   ```
   Error: failed to connect to database
   ```
   - Check if PostgreSQL is running
   - Verify DATABASE_URL is correct
   - Check firewall/network settings

2. **Migration errors**
   ```
   Error: migration failed
   ```
   - Check migration file syntax
   - Verify database permissions
   - Check for conflicting changes

3. **Facebook webhook errors**
   ```
   Error: Facebook verify token mismatch
   ```
   - Verify FACEBOOK_VERIFY_TOKEN is set correctly
   - Check Facebook app configuration
   - Ensure webhook URL is accessible

### Debug Mode

Enable debug logging:

```bash
LOG_LEVEL=debug
GIN_MODE=debug
```

### Database Logs

Check PostgreSQL logs:

```bash
# Docker
docker-compose logs db

# Local PostgreSQL
tail -f /usr/local/var/log/postgresql@16.log
```

## API Endpoints

The application provides these database-related endpoints:

- `GET /health` - Health check including database status
- `GET /api/keywords` - List all keyword-response mappings
- `POST /api/keywords` - Create/update keyword mappings
- `POST /webhook/facebook` - Facebook webhook (uses database)

## Next Steps

After setting up the database:

1. Configure your Facebook app webhook
2. Test with a Facebook page
3. Add custom keyword-response mappings
4. Monitor analytics and logs
5. Scale connection pooling as needed

For more details, see:
- `internal/database/README.md` - Technical database documentation
- `internal/models/` - Data model documentation
- `cmd/migrate/` - Migration tool documentation
