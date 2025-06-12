# 🌐 Multi-Channel Auto Reply Architecture

## 📋 Overview
Architecture design for supporting multiple social media and e-commerce platforms with centralized auto-reply management.

---

## 🎯 **Target Multi-Channel Structure**

```
ruaymanagement/
├── frontend/                               # Next.js Dashboard
├── backend/
│   ├── 🔌 external/                      # EXTERNAL INTEGRATIONS
│   │   ├── FacebookConnect/              # Port 8085 - Facebook API
│   │   │   ├── go.mod
│   │   │   ├── cmd/server/main.go
│   │   │   ├── handlers/webhook.go
│   │   │   ├── services/facebook_api.go
│   │   │   └── models/webhook.go
│   │   │
│   │   ├── LineConnect/                  # Port 8089 - LINE API (existing)
│   │   │   ├── go.mod
│   │   │   ├── cmd/server/main.go
│   │   │   ├── handlers/webhook.go
│   │   │   ├── services/line_api.go
│   │   │   └── models/webhook.go
│   │   │
│   │   ├── TikTokConnect/                # Port 8087 - TikTok API
│   │   │   ├── go.mod
│   │   │   ├── cmd/server/main.go
│   │   │   ├── handlers/webhook.go
│   │   │   ├── services/tiktok_api.go
│   │   │   └── models/webhook.go
│   │   │
│   │   ├── LazadaConnect/                # Port 8090 - Lazada API
│   │   │   ├── go.mod
│   │   │   ├── cmd/server/main.go
│   │   │   ├── handlers/webhook.go
│   │   │   ├── services/lazada_api.go
│   │   │   └── models/webhook.go
│   │   │
│   │   └── ShopeeConnect/                # Port 8091 - Shopee API
│   │       ├── go.mod
│   │       ├── cmd/server/main.go
│   │       ├── handlers/webhook.go
│   │       ├── services/shopee_api.go
│   │       └── models/webhook.go
│   │
│   └── 🏢 internal/                      # INTERNAL BUSINESS LOGIC
│       ├── AutoReplyManagement/          # Port 8086 - Centralized Auto Reply
│       │   ├── go.mod
│       │   ├── cmd/server/main.go
│       │   ├── handlers/
│       │   │   ├── keywords.go           # Keyword management
│       │   │   ├── analytics.go          # Analytics
│       │   │   ├── channels.go           # Channel management
│       │   │   └── rules.go              # Business rules
│       │   ├── services/
│       │   │   ├── autoreply.go          # Auto reply logic
│       │   │   ├── keyword_matcher.go    # Keyword matching
│       │   │   ├── channel_router.go     # Channel routing
│       │   │   └── analytics.go          # Analytics service
│       │   ├── models/
│       │   │   ├── message.go            # Universal message model
│       │   │   ├── channel.go            # Channel model
│       │   │   ├── keyword.go            # Keyword model
│       │   │   └── analytics.go          # Analytics model
│       │   └── repository/
│       │       ├── keyword_repo.go
│       │       ├── message_repo.go
│       │       ├── channel_repo.go
│       │       └── analytics_repo.go
│       │
│       ├── ChannelManagement/            # Port 8092 - Channel Configuration
│       │   ├── go.mod
│       │   ├── cmd/server/main.go
│       │   ├── handlers/
│       │   │   ├── channels.go           # Channel CRUD
│       │   │   ├── settings.go           # Channel settings
│       │   │   └── credentials.go        # API credentials
│       │   ├── services/
│       │   │   ├── channel.go            # Channel management
│       │   │   ├── validation.go         # Credential validation
│       │   │   └── health_check.go       # Channel health monitoring
│       │   └── models/
│       │       ├── channel.go            # Channel configuration
│       │       └── credential.go         # API credentials
│       │
│       └── MessageRouter/                # Port 8093 - Message Orchestration
│           ├── go.mod
│           ├── cmd/server/main.go
│           ├── handlers/
│           │   ├── webhook.go            # Universal webhook receiver
│           │   └── router.go             # Message routing
│           ├── services/
│           │   ├── message_router.go     # Route messages to channels
│           │   ├── rate_limiter.go       # Rate limiting
│           │   └── queue_manager.go      # Message queue management
│           └── models/
│               ├── routing_rule.go       # Routing rules
│               └── queue_message.go      # Queue message model
```

---

## 🔄 **Data Flow Architecture**

### A. **Incoming Message Flow**
```
Platform → ChannelConnect → MessageRouter → AutoReplyManagement → MessageRouter → ChannelConnect → Platform
   ↓             ↓              ↓                    ↓                   ↓              ↓           ↓
Facebook     FacebookConnect   Route to         Process Auto        Route Reply    Send via     Facebook
LINE         LineConnect       AutoReply        Reply Logic         Back           Facebook     LINE
TikTok       TikTokConnect     Management       + Analytics         to Channel     API          TikTok
Lazada       LazadaConnect                                         Connect                      Lazada
Shopee       ShopeeConnect                                                                      Shopee
```

### B. **Service Communication**
```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│  ChannelConnect │    │  MessageRouter  │    │ AutoReplyMgmt   │
│  (External)     │    │   (Internal)    │    │   (Internal)    │
├─────────────────┤    ├─────────────────┤    ├─────────────────┤
│ • Webhook       │───▶│ • Route Message │───▶│ • Match Keyword │
│ • API Calls     │    │ • Rate Limit    │    │ • Generate Reply│
│ • Auth/Verify   │◀───│ • Queue Mgmt    │◀───│ • Analytics     │
└─────────────────┘    └─────────────────┘    └─────────────────┘
         ▲                       ▲                       ▲
         │                       │                       │
         ▼                       ▼                       ▼
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│ ChannelMgmt     │    │   Database      │    │   Analytics     │
│ (Internal)      │    │  (PostgreSQL)   │    │   Dashboard     │
└─────────────────┘    └─────────────────┘    └─────────────────┘
```

---

## 📊 **Database Schema Design**

### A. **Multi-Channel Schema**
```sql
-- Channel Management
CREATE SCHEMA channel_management;

-- Channels table
CREATE TABLE channel_management.channels (
    id SERIAL PRIMARY KEY,
    name VARCHAR(50) NOT NULL,           -- 'facebook', 'line', 'tiktok', etc.
    display_name VARCHAR(100) NOT NULL,  -- 'Facebook Pages', 'LINE Official', etc.
    is_active BOOLEAN DEFAULT true,
    port INTEGER,                        -- Service port
    api_version VARCHAR(20),
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- Channel Credentials
CREATE TABLE channel_management.channel_credentials (
    id SERIAL PRIMARY KEY,
    channel_id INTEGER REFERENCES channel_management.channels(id),
    credential_type VARCHAR(50),         -- 'page_token', 'verify_token', etc.
    credential_value TEXT,               -- Encrypted credential
    is_active BOOLEAN DEFAULT true,
    expires_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT NOW()
);

-- Auto Reply Schema
CREATE SCHEMA auto_reply;

-- Keywords (Multi-channel support)
CREATE TABLE auto_reply.keywords (
    id SERIAL PRIMARY KEY,
    channel_id INTEGER REFERENCES channel_management.channels(id),
    keyword TEXT NOT NULL,
    response TEXT NOT NULL,
    is_active BOOLEAN DEFAULT true,
    match_type VARCHAR(20) DEFAULT 'exact', -- 'exact', 'partial', 'fuzzy'
    priority INTEGER DEFAULT 0,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- Universal Message Log
CREATE TABLE auto_reply.message_logs (
    id SERIAL PRIMARY KEY,
    channel_id INTEGER REFERENCES channel_management.channels(id),
    external_message_id VARCHAR(255),   -- Platform-specific message ID
    sender_id VARCHAR(255),
    sender_name VARCHAR(255),
    message_text TEXT,
    matched_keyword VARCHAR(500),
    response_sent TEXT,
    processing_time_ms INTEGER,
    created_at TIMESTAMP DEFAULT NOW()
);

-- Channel Analytics
CREATE TABLE auto_reply.channel_analytics (
    id SERIAL PRIMARY KEY,
    channel_id INTEGER REFERENCES channel_management.channels(id),
    date DATE,
    messages_received INTEGER DEFAULT 0,
    messages_replied INTEGER DEFAULT 0,
    unique_senders INTEGER DEFAULT 0,
    avg_response_time_ms DECIMAL(10,2),
    created_at TIMESTAMP DEFAULT NOW()
);
```

---

## 🔌 **Universal API Contracts**

### A. **Channel Connect Interface**
```go
// Universal Channel Interface
type ChannelConnector interface {
    // Webhook handling
    VerifyWebhook(ctx context.Context, req WebhookVerifyRequest) (*WebhookVerifyResponse, error)
    HandleWebhook(ctx context.Context, payload interface{}) error
    
    // Message sending
    SendMessage(ctx context.Context, req SendMessageRequest) (*SendMessageResponse, error)
    
    // Health check
    HealthCheck(ctx context.Context) (*HealthStatus, error)
}

// Universal Message Models
type UniversalMessage struct {
    ChannelID       int       `json:"channel_id"`
    ExternalID      string    `json:"external_id"`
    SenderID        string    `json:"sender_id"`
    SenderName      string    `json:"sender_name"`
    MessageText     string    `json:"message_text"`
    MessageType     string    `json:"message_type"` // text, image, sticker, etc.
    Timestamp       time.Time `json:"timestamp"`
    RawPayload      json.RawMessage `json:"raw_payload"`
}

type ProcessMessageRequest struct {
    Message     UniversalMessage `json:"message"`
    ChannelInfo ChannelInfo      `json:"channel_info"`
    Context     MessageContext   `json:"context"`
}

type ProcessMessageResponse struct {
    ShouldReply     bool    `json:"should_reply"`
    Response        string  `json:"response"`
    MatchedKeyword  string  `json:"matched_keyword,omitempty"`
    Confidence      float64 `json:"confidence"`
    ResponseType    string  `json:"response_type"` // text, template, quick_reply
}
```

### B. **Service Endpoints**

#### 🔌 **FacebookConnect** (Port 8085)
```
GET  /webhook/facebook        # Webhook verification
POST /webhook/facebook        # Webhook handler
GET  /health                  # Health check
POST /api/send-message        # Send message
```

#### 🔌 **LineConnect** (Port 8089)
```
GET  /webhook/line           # Webhook verification
POST /webhook/line           # Webhook handler
GET  /health                 # Health check
POST /api/send-message       # Send message
```

#### 🔌 **TikTokConnect** (Port 8087)
```
GET  /webhook/tiktok         # Webhook verification
POST /webhook/tiktok         # Webhook handler
GET  /health                 # Health check
POST /api/send-message       # Send message
```

#### 🏢 **AutoReplyManagement** (Port 8086)
```
POST /api/process-message    # Process incoming message
GET  /api/keywords           # Get keywords
POST /api/keywords           # Save keywords
GET  /api/analytics          # Get analytics data
GET  /health                 # Health check
```

#### 🏢 **ChannelManagement** (Port 8092)
```
GET  /api/channels           # List channels
POST /api/channels           # Create channel
PUT  /api/channels/:id       # Update channel
GET  /api/channels/:id/health # Channel health status
POST /api/channels/:id/credentials # Update credentials
```

#### 🏢 **MessageRouter** (Port 8093)
```
POST /api/route-message      # Route incoming message
GET  /api/queue-status       # Queue status
GET  /health                 # Health check
```

---

## 🎛️ **Configuration Management**

### A. **Environment Variables**
```env
# Facebook
FACEBOOK_PAGE_TOKEN=xxx
FACEBOOK_VERIFY_TOKEN=xxx
FACEBOOK_API_VERSION=v18.0

# LINE  
LINE_CHANNEL_ACCESS_TOKEN=xxx
LINE_CHANNEL_SECRET=xxx

# TikTok
TIKTOK_CLIENT_KEY=xxx
TIKTOK_CLIENT_SECRET=xxx
TIKTOK_WEBHOOK_VERIFY_TOKEN=xxx

# Lazada
LAZADA_APP_KEY=xxx
LAZADA_APP_SECRET=xxx
LAZADA_WEBHOOK_SECRET=xxx

# Shopee
SHOPEE_PARTNER_ID=xxx
SHOPEE_PARTNER_KEY=xxx
SHOPEE_WEBHOOK_SECRET=xxx

# Database
DB_HOST=postgres
DB_PORT=5432
DB_NAME=loyverse_cache
DB_USER=postgres
DB_PASSWORD=password

# Service Ports
FACEBOOK_CONNECT_PORT=8085
LINE_CONNECT_PORT=8089
TIKTOK_CONNECT_PORT=8087
LAZADA_CONNECT_PORT=8090
SHOPEE_CONNECT_PORT=8091
AUTO_REPLY_PORT=8086
CHANNEL_MGMT_PORT=8092
MESSAGE_ROUTER_PORT=8093
```

### B. **Docker Compose Structure**
```yaml
version: '3.8'

services:
  # Database
  postgres:
    image: postgres:13
    ports:
      - "5433:5432"
    
  # External Services
  facebook-connect:
    build: ./backend/external/FacebookConnect
    ports:
      - "8085:8085"
    depends_on:
      - postgres
      - auto-reply-management
      
  line-connect:
    build: ./backend/external/LineConnect
    ports:
      - "8089:8089"
    depends_on:
      - postgres
      - auto-reply-management
      
  tiktok-connect:
    build: ./backend/external/TikTokConnect
    ports:
      - "8087:8087"
    depends_on:
      - postgres
      - auto-reply-management
      
  lazada-connect:
    build: ./backend/external/LazadaConnect
    ports:
      - "8090:8090"
    depends_on:
      - postgres
      - auto-reply-management
      
  shopee-connect:
    build: ./backend/external/ShopeeConnect
    ports:
      - "8091:8091"
    depends_on:
      - postgres
      - auto-reply-management
  
  # Internal Services
  auto-reply-management:
    build: ./backend/internal/AutoReplyManagement
    ports:
      - "8086:8086"
    depends_on:
      - postgres
      
  channel-management:
    build: ./backend/internal/ChannelManagement
    ports:
      - "8092:8092"
    depends_on:
      - postgres
      
  message-router:
    build: ./backend/internal/MessageRouter
    ports:
      - "8093:8093"
    depends_on:
      - postgres
      - auto-reply-management
      
  # Frontend
  frontend:
    build: ./frontend
    ports:
      - "3000:3000"
    depends_on:
      - auto-reply-management
      - channel-management
```

---

## 🎯 **Benefits of Multi-Channel Architecture**

### A. **Scalability**
- ✅ แต่ละ channel scale ได้อิสระ
- ✅ เพิ่ม channel ใหม่ได้ง่าย
- ✅ Load balancing ได้ตาม channel

### B. **Maintainability**  
- ✅ แยก API integration จาก business logic
- ✅ Update API ไม่กระทบ business logic
- ✅ Test แต่ละ channel ได้อิสระ

### C. **Flexibility**
- ✅ กำหนด keyword ต่าง channel ได้
- ✅ รองรับ message type ที่แตกต่างกัน
- ✅ Custom business rules ต่อ channel

### D. **Monitoring**
- ✅ Analytics แยกตาม channel
- ✅ Health monitoring ครอบคลุม
- ✅ Performance metrics ต่อ channel

---

## 🚀 **Implementation Phases**

### Phase 1: **Core Infrastructure**
1. สร้าง AutoReplyManagement (Internal)
2. สร้าง ChannelManagement (Internal)  
3. สร้าง MessageRouter (Internal)
4. Setup Database Schema

### Phase 2: **Facebook Integration**
1. สร้าง FacebookConnect (External)
2. Integration testing
3. Deploy และ test

### Phase 3: **Additional Channels**
1. LineConnect (existing, migrate)
2. TikTokConnect (new)
3. LazadaConnect (new)
4. ShopeeConnect (new)

### Phase 4: **Advanced Features**
1. AI/ML integration
2. Advanced analytics
3. A/B testing
4. Performance optimization

---

## 📊 **Service Summary**

| Type | Service | Port | Technology | Status |
|------|---------|------|------------|--------|
| 🔌 **External** | FacebookConnect | 8085 | Go 1.22 | ⚠️ New |
| 🔌 **External** | LineConnect | 8089 | Go 1.22 | ✅ Exists |
| 🔌 **External** | TikTokConnect | 8087 | Go 1.22 | ⚠️ New |
| 🔌 **External** | LazadaConnect | 8090 | Go 1.22 | ⚠️ New |
| 🔌 **External** | ShopeeConnect | 8091 | Go 1.22 | ⚠️ New |
| 🏢 **Internal** | AutoReplyManagement | 8086 | Go 1.22 | ⚠️ New |
| 🏢 **Internal** | ChannelManagement | 8092 | Go 1.22 | ⚠️ New |
| 🏢 **Internal** | MessageRouter | 8093 | Go 1.22 | ⚠️ New |

---

*This multi-channel architecture provides a scalable, maintainable foundation for supporting multiple social media and e-commerce platforms with centralized auto-reply management.*
