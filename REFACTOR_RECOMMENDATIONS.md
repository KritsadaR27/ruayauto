# 🔄 Refactoring Recommendations for Facebook Webhook Integration

## 📋 Overview
Recommendations for integrating Facebook Webhook service into the ruaymanagement ecosystem.

---

## 🎯 **1. Directory Structure Refactoring**

### Current Structure:
```
backend/
├── cmd/server/
├── internal/
│   ├── config/
│   ├── handler/
│   ├── models/
│   ├── service/
│   └── database/
```

### Target Structure (in ruaymanagement):

#### 🔌 **External Service - FacebookConnect** (Port 8085)
```
backend/external/FacebookConnect/           # Facebook API Integration
├── go.mod                                 # Go 1.22+
├── Dockerfile
├── cmd/
│   └── server/
│       └── main.go
├── handlers/                              # HTTP handlers
│   ├── webhook.go                         # Webhook receiver
│   └── verification.go                    # Webhook verification
├── services/                              # Facebook API services
│   └── facebook_api.go                    # Facebook Graph API client
├── models/                                # Facebook data models
│   └── webhook.go                         # Webhook payload models
├── config/                                # Configuration
│   └── config.go
└── utils/                                 # Facebook utilities
    └── facebook_client.go
```

#### 🏢 **Internal Service - AutoReplyManagement** (Port 8086)
```
backend/internal/AutoReplyManagement/       # Auto Reply Business Logic
├── go.mod                                 # Go 1.22+
├── Dockerfile
├── cmd/
│   └── server/
│       └── main.go
├── handlers/                              # HTTP handlers
│   ├── keywords.go                        # Keyword management
│   ├── analytics.go                       # Analytics endpoints
│   └── health.go                          # Health check
├── services/                              # Business logic
│   ├── autoreply.go                       # Auto reply logic
│   ├── keyword.go                         # Keyword matching
│   └── analytics.go                       # Analytics service
├── models/                                # Data models
│   ├── keyword.go                         # Keyword models
│   └── analytics.go                       # Analytics models
├── repository/                            # Data access
│   ├── keyword_repo.go
│   └── analytics_repo.go
├── config/                                # Configuration
│   └── config.go
└── utils/                                 # Business utilities
    └── text_matcher.go
```

---

## 🔧 **2. Code Refactoring Points**

### A. Module Naming
```go
// Current
module ruayautomsg

// Target
module ruaymanagement/backend/external/FacebookConnect
```

### B. Package Import Paths
```go
// Current
import "ruayautomsg/internal/config"

// Target
import "ruaymanagement/backend/external/FacebookConnect/config"
```

### C. Port Configuration
```go
// Current: Port 8100
// Target: Port 8085 (following ruaymanagement pattern)

// In docker-compose.yml
services:
  facebook-connect:
    ports:
      - "8085:8085"
```

---

## 🏗️ **3. Architecture Improvements**

### A. Service Layer Pattern
```go
type FacebookService interface {
    HandleWebhook(ctx context.Context, payload WebhookPayload) error
    VerifyWebhook(mode, token, challenge string) (string, error)
    ProcessComment(ctx context.Context, comment CommentData) error
}

type KeywordService interface {
    MatchKeyword(ctx context.Context, text string) (*MatchResult, error)
    GetAllKeywords(ctx context.Context) ([]models.Keyword, error)
    SaveKeywords(ctx context.Context, keywords []KeywordPair) error
}
```

### B. Repository Pattern (Already Good!)
```go
// Keep existing repository pattern
type KeywordRepository interface {
    GetAll(ctx context.Context) ([]models.Keyword, error)
    GetByKeyword(ctx context.Context, keyword string) (*models.Keyword, error)
    Create(ctx context.Context, keyword models.Keyword) error
    // ...
}
```

### C. Error Handling Enhancement
```go
// Add structured error types
type FacebookError struct {
    Code    string `json:"code"`
    Message string `json:"message"`
    Type    string `json:"type"`
}

func (e FacebookError) Error() string {
    return fmt.Sprintf("[%s] %s: %s", e.Type, e.Code, e.Message)
}
```

---

## 📊 **4. Database Integration**

### A. Connection Management
```go
// Use ruaymanagement's PostgreSQL instance (Port 5433)
// Database: loyverse_cache (shared database)

// Create dedicated schema for Facebook service
CREATE SCHEMA facebook_webhook;

// Tables with proper prefixes
facebook_webhook.keywords
facebook_webhook.message_analytics
facebook_webhook.webhook_logs
```

### B. Migration Strategy
```sql
-- Create dedicated schema
CREATE SCHEMA IF NOT EXISTS facebook_webhook;

-- Move existing tables
ALTER TABLE keywords SET SCHEMA facebook_webhook;
ALTER TABLE message_analytics SET SCHEMA facebook_webhook;
ALTER TABLE webhook_logs SET SCHEMA facebook_webhook;
```

---

## 🔌 **5. API Standardization**

### A. REST API Endpoints
```go
// Current endpoints (keep compatible)
GET  /webhook/facebook     → Webhook verification
POST /webhook/facebook     → Webhook handler
GET  /api/keywords         → Get all keywords
POST /api/keywords         → Save keywords

// Additional endpoints for integration
GET  /health               → Health check
GET  /metrics              → Service metrics
GET  /api/analytics        → Analytics data
```

### B. Response Format Standardization
```go
type APIResponse struct {
    Success bool        `json:"success"`
    Data    interface{} `json:"data,omitempty"`
    Error   *APIError   `json:"error,omitempty"`
    Meta    *Meta       `json:"meta,omitempty"`
}

type APIError struct {
    Code    string `json:"code"`
    Message string `json:"message"`
}
```

---

## 🐳 **6. Docker Configuration**

### A. Dockerfile Optimization
```dockerfile
# Multi-stage build
FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o server ./cmd/server

FROM alpine:latest
RUN apk --no-cache add ca-certificates tzdata
WORKDIR /root/
COPY --from=builder /app/server ./
EXPOSE 8085
CMD ["./server"]
```

### B. Environment Variables
```env
# Facebook Configuration
FACEBOOK_PAGE_TOKEN=your_page_token
FACEBOOK_VERIFY_TOKEN=your_verify_token
FACEBOOK_API_VERSION=v18.0

# Database Configuration (shared with ruaymanagement)
DB_HOST=postgres
DB_PORT=5432
DB_NAME=loyverse_cache
DB_USER=postgres
DB_PASSWORD=password
DB_SCHEMA=facebook_webhook

# Service Configuration
PORT=8085
GIN_MODE=release
LOG_LEVEL=info

# Integration Configuration
ENABLE_ANALYTICS=true
ENABLE_RATE_LIMITING=true
RATE_LIMIT_REQUESTS=100
RATE_LIMIT_WINDOW=60
```

---

## 🔄 **7. Migration Steps**

### Phase 1: Code Refactoring
1. Update module name and import paths
2. Restructure directories to match target pattern
3. Update configuration management
4. Add proper error handling

### Phase 2: Database Integration  
1. Create facebook_webhook schema in shared database
2. Migrate existing tables to new schema
3. Update repository layer for schema-aware queries
4. Test database connectivity

### Phase 3: Service Integration
1. Update docker-compose.yml in ruaymanagement
2. Add FacebookConnect service configuration
3. Update port mappings (8085)
4. Test service communication

### Phase 4: API Integration
1. Standardize API response formats
2. Add health check endpoints
3. Implement metrics collection
4. Update documentation

---

## 🔄 **Service Communication Architecture**

### A. Data Flow
```
Facebook → FacebookConnect → AutoReplyManagement → FacebookConnect → Facebook
         (Webhook)        (Process Logic)      (Send Reply)
```

### B. Service Responsibilities

#### 🔌 **FacebookConnect** (Port 8085)
- ✅ รับ webhook จาก Facebook
- ✅ Verify webhook signature
- ✅ ส่งข้อความตอบกลับผ่าน Facebook Graph API
- ✅ จัดการ Facebook API credentials
- ❌ **ไม่มี** business logic การประมวลผลคำตอบ

#### 🏢 **AutoReplyManagement** (Port 8086)  
- ✅ จัดการ keywords และ responses
- ✅ ประมวลผลข้อความและหาคำตอบที่เหมาะสม
- ✅ Analytics และ logging
- ✅ Machine learning / AI processing (future)
- ❌ **ไม่รู้จัก** Facebook API directly

### C. API Communication
```go
// FacebookConnect calls AutoReplyManagement
POST http://autoreply-management:8086/api/process-message
{
    "message": "สวัสดีครับ",
    "sender_id": "123456789",
    "context": {...}
}

// AutoReplyManagement responds
{
    "should_reply": true,
    "response": "สวัสดีค่ะ ยินดีต้อนรับ",
    "matched_keyword": "สวัสดี",
    "confidence": 0.95
}
```

---

## 📋 **8. Benefits After Refactoring**

### A. Consistency
- ✅ Follows ruaymanagement patterns
- ✅ Consistent naming conventions
- ✅ Standardized error handling
- ✅ Unified configuration management

### B. Maintainability
- ✅ Clear separation of concerns
- ✅ Proper dependency injection
- ✅ Repository pattern for data access
- ✅ Service layer for business logic

### C. Scalability
- ✅ Microservice architecture
- ✅ Database schema isolation
- ✅ Independent deployment
- ✅ Health monitoring

### D. Integration Ready
- ✅ Compatible with existing services
- ✅ Shared database connection
- ✅ Standardized API responses
- ✅ Docker compose integration

---

## ⚠️ **9. Migration Risks & Mitigation**

### Risks:
1. **Breaking Changes**: Module name and import path changes
2. **Database Migration**: Schema changes and data migration
3. **Service Dependencies**: Port conflicts and configuration changes

### Mitigation:
1. **Backward Compatibility**: Keep existing API endpoints
2. **Gradual Migration**: Phase-by-phase implementation
3. **Testing**: Comprehensive testing at each phase
4. **Rollback Plan**: Ability to revert changes if needed

---

## 🎯 **10. Next Steps**

1. **Review** this refactoring plan
2. **Create** backup of current code
3. **Start** with Phase 1 (Code Refactoring)
4. **Test** each phase thoroughly
5. **Document** any issues or changes
6. **Deploy** to development environment first

---

*This refactoring will make the Facebook Webhook service ready for integration into the ruaymanagement ecosystem while maintaining backward compatibility and improving code quality.*
