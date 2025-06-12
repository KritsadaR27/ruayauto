# Service Architecture Reorganization Plan
# Microservices: FacebookConnect (External) + ChatBotCore (Internal)

## 🎯 TARGET ARCHITECTURE

### 📂 /backend/external/FacebookConnect/ (Customer-Facing Integration)
```
external/FacebookConnect/
├── cmd/server/main.go                    # FacebookConnect server (port 8085)
├── config/config.go                      # FB-specific config
├── handlers/
│   ├── webhook.go                        # Facebook webhook receiver
│   ├── verification.go                   # Facebook verification
│   └── health.go                         # Health check
├── models/
│   ├── facebook_message.go               # Facebook API models
│   └── webhook_payload.go                # Webhook payload models
├── services/
│   ├── facebook_api.go                   # Facebook API client
│   ├── message_transformer.go            # FB format ↔ Internal format
│   └── webhook_processor.go              # Webhook processing
├── utils/
│   ├── facebook_auth.go                  # FB token management
│   └── rate_limiter.go                   # API rate limiting
└── Dockerfile                            # FacebookConnect container
```

### 📂 /backend/internal/ChatBotCore/ (Business Logic Core)
```
internal/ChatBotCore/
├── cmd/server/main.go                    # ChatBotCore server (port 8086)
├── config/config.go                      # Core business config
├── models/
│   ├── conversation.go                   # Business conversation model
│   ├── message.go                        # Internal message model
│   ├── keyword.go                        # Keyword business model
│   ├── response_template.go              # Response templates
│   └── analytics.go                      # Analytics models
├── repository/
│   ├── interfaces.go                     # Repository contracts
│   ├── conversation.go                   # Conversation data access
│   ├── message.go                        # Message data access
│   ├── keyword.go                        # Keyword data access
│   ├── response_template.go              # Template data access
│   └── analytics.go                      # Analytics data access
├── services/
│   ├── chatbot.go                        # Core chatbot business logic
│   ├── keyword_matching.go               # Keyword matching algorithms
│   ├── conversation_manager.go           # Conversation lifecycle
│   ├── analytics.go                      # Business analytics
│   └── message_processor.go              # Message processing pipeline
├── handlers/
│   ├── message.go                        # Internal message API
│   ├── keyword.go                        # Keyword management API
│   ├── analytics.go                      # Analytics API
│   └── health.go                         # Health check
└── Dockerfile                            # ChatBotCore container
```

## 🔄 DATA FLOW DESIGN

### 📥 Incoming Message Flow:
```
Facebook User Comment
    ↓
[FacebookConnect:8085] webhook.go
    ↓ Transform FB → Internal format
[FacebookConnect] message_transformer.go
    ↓ HTTP POST /api/messages/process
[ChatBotCore:8086] message.go handler
    ↓ Business logic processing
[ChatBotCore] chatbot.go service
    ↓ Keyword matching + Response generation
[ChatBotCore] keyword_matching.go
    ↓ HTTP Response with reply
[FacebookConnect] webhook.go
    ↓ Transform Internal → FB format
[FacebookConnect] facebook_api.go
    ↓ POST to Facebook API
Facebook User sees reply
```

### 📤 Admin Management Flow:
```
Frontend Admin Panel
    ↓ HTTP POST /api/keywords
[ChatBotCore:8086] keyword.go handler
    ↓ Business validation
[ChatBotCore] keyword matching service
    ↓ Database operations
[ChatBotCore] keyword repository
    ↓ PostgreSQL
Database
```

## ⚙️ SERVICE COMMUNICATION

### 🔌 FacebookConnect → ChatBotCore API:
```go
// POST http://chatbot-core:8086/api/messages/process
type ProcessMessageRequest struct {
    Content        string    `json:"content"`
    SenderID       string    `json:"sender_id"`
    ConversationID string    `json:"conversation_id"`
    Timestamp      time.Time `json:"timestamp"`
    MessageType    string    `json:"message_type"` // "text", "image", etc.
}

type ProcessMessageResponse struct {
    Success        bool     `json:"success"`
    Reply          string   `json:"reply,omitempty"`
    MatchedKeyword string   `json:"matched_keyword,omitempty"`
    ShouldReply    bool     `json:"should_reply"`
    Error          string   `json:"error,omitempty"`
}
```

### 🔄 ChatBotCore → FacebookConnect (Future):
```go
// POST http://facebook-connect:8085/api/messages/send
type SendMessageRequest struct {
    RecipientID string `json:"recipient_id"`
    Content     string `json:"content"`
    MessageType string `json:"message_type"`
}
```

## 🏗️ IMPLEMENTATION STRATEGY

### Phase 1: Clean Separation ✅
1. Move Facebook-specific code to `external/FacebookConnect/`
2. Keep only business logic in `internal/ChatBotCore/`
3. Define clear API contracts between services

### Phase 2: Service Communication 🔄
1. Implement HTTP API in ChatBotCore for message processing
2. Update FacebookConnect to call ChatBotCore API
3. Remove direct database access from FacebookConnect

### Phase 3: Containerization 🐳
1. Separate Docker containers for each service
2. Docker Compose orchestration
3. Service discovery and health checks

## 📝 KEY PRINCIPLES COMPLIANCE

✅ **External Service (FacebookConnect):**
- ONLY handles Facebook API integration
- NO business logic (keyword matching)
- NO direct database access
- Transforms FB format ↔ Internal format
- Handles authentication & rate limiting

✅ **Internal Service (ChatBotCore):**
- ONLY business logic (keyword matching, conversation management)
- NO external API calls (Facebook)
- Platform-agnostic message processing
- Database operations through repository pattern

✅ **Replaceability:**
- Can replace FacebookConnect with TelegramConnect
- Can upgrade ChatBotCore (add AI) without touching FacebookConnect
- Clear API boundaries prevent tight coupling

✅ **Scalability:**
- Each service scales independently
- Stateless service design
- Database connection pooling

This architecture ensures clean separation of concerns, maintainability, and future extensibility for multi-channel support and AI integration.
