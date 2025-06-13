# Database Package

This package provides PostgreSQL database connectivity using pgx/v5 with connection pooling, health monitoring, and automatic migrations.

## Features

- **PostgreSQL Connection**: Using pgx/v5 for high-performance PostgreSQL connectivity
- **Connection Pooling**: Configurable connection pool with health checks
- **Auto Migrations**: Automatic database migrations on startup
- **Health Monitoring**: Continuous health monitoring with metrics
- **Graceful Shutdown**: Proper connection cleanup on application shutdown

## Components

### Database Connection (`database.go`)
- Main database connection manager
- Connection pooling configuration
- Transaction support
- Query execution helpers

### Health Monitoring (`health.go`)
- Continuous health monitoring
- Pool statistics tracking
- Health status reporting
- Configurable check intervals

### Migration System (`migrations.go`)
- Automatic migration execution
- Rollback support
- Migration status tracking
- File-based migration source

### Manager (`manager.go`)
- High-level database manager
- Integration of all components
- Graceful shutdown handling
- Configuration management

## Usage

### Basic Setup

```go
import "ruayautomsg/internal/database"

// Create manager with auto-migrations
mgr, err := database.NewManagerFromURL(databaseURL, "./migrations", true)
if err != nil {
    log.Fatal(err)
}
defer mgr.Shutdown(context.Background())

// Get database connection
db := mgr.GetDB()
```

### Configuration

```go
config := &database.Config{
    Host:            "localhost",
    Port:            "5432", 
    User:            "username",
    Password:        "password",
    Database:        "dbname",
    SSLMode:         "disable",
    MaxConnections:  25,
    MinConnections:  5,
    MaxLifetime:     time.Hour,
    MaxIdleTime:     time.Minute * 30,
    ConnectTimeout:  time.Second * 10,
    HealthCheckTime: time.Minute,
}

mgr, err := database.NewManager(&database.ManagerConfig{
    Database: config,
    Migration: &database.MigrationConfig{
        MigrationsPath: "./migrations",
        DatabaseName:   "mydb",
    },
    AutoMigrate: true,
})
```

### Health Monitoring

```go
// Check if database is healthy
if mgr.IsHealthy() {
    log.Println("Database is healthy")
}

// Get detailed health status
status := mgr.GetHealthStatus()
fmt.Printf("Pool stats: %+v\n", status.PoolStats)
```

### Migrations

```go
// Manual migration commands
mgr.RunMigrations()           // Run pending migrations
mgr.RollbackMigration()       // Rollback last migration

// Using the migrate CLI tool
go run ./cmd/migrate -up      // Run migrations up
go run ./cmd/migrate -down    // Rollback migration  
go run ./cmd/migrate -status  // Check migration status
```

## Migration Files

Migration files should be placed in `internal/migrations/` directory:

```
internal/migrations/
├── 000001_create_keywords_table.up.sql
├── 000001_create_keywords_table.down.sql
├── 000002_create_analytics_tables.up.sql
└── 000002_create_analytics_tables.down.sql
```

### Example Migration (up)

```sql
-- 000001_create_keywords_table.up.sql
CREATE TABLE IF NOT EXISTS keywords (
    id SERIAL PRIMARY KEY,
    keyword VARCHAR(255) NOT NULL UNIQUE,
    response TEXT NOT NULL,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_keywords_keyword ON keywords(keyword);
```

### Example Migration (down)

```sql
-- 000001_create_keywords_table.down.sql
DROP INDEX IF EXISTS idx_keywords_keyword;
DROP TABLE IF EXISTS keywords;
```

## Environment Variables

```bash
# Database Configuration
DATABASE_URL="postgres://user:password@localhost:5432/dbname?sslmode=disable"
DB_HOST="localhost"
DB_PORT="5432"
DB_USER="username"
DB_PASSWORD="password"
DB_NAME="dbname"
DB_SSLMODE="disable"
```

## CLI Tools

### Migrate Command

```bash
# Build the migrate tool
go build ./cmd/migrate

# Run migrations up
./migrate -up

# Rollback last migration
./migrate -down

# Check migration status
./migrate -status

# Create new migration
./migrate -create create_users_table
```

## Testing

```go
// Test database connection
err := database.TestConnection(databaseURL)
if err != nil {
    log.Fatal(err)
}

// Setup database with migrations
err = database.SetupDatabase(databaseURL)
if err != nil {
    log.Fatal(err)
}
```

## Production Recommendations

1. **Connection Pool Sizing**: Adjust `MaxConnections` based on your application load
2. **Health Check Interval**: Set appropriate health check intervals for your environment
3. **Migration Strategy**: Always test migrations in staging before production
4. **Monitoring**: Monitor pool statistics and health metrics
5. **Backup**: Ensure proper database backups before running migrations

## Error Handling

The package provides detailed error messages and proper error wrapping:

```go
mgr, err := database.NewManagerFromURL(databaseURL, "./migrations", true)
if err != nil {
    // Handle connection errors
    log.Printf("Database connection failed: %v", err)
    return
}

if !mgr.IsHealthy() {
    status := mgr.GetHealthStatus()
    log.Printf("Database unhealthy: %s", status.LastError)
}
```
