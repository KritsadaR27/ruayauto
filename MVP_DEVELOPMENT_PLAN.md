# 🚀 MVP Development Plan - Simplified ChatBot

## 📋 **MVP Scope Summary**

### ✅ **What We'll Build**
- **Facebook Integration Only** - เฉพาะ Facebook ก่อน
- **Keyword Detection & Auto-Reply** - ดักคำและตอบอัตโนมัติ  
- **Comment Monitoring** - ดึงคอมเม้นต์เข้า inbox
- **Chat Storage** - เก็บประวัติแชทไว้ดูภายหลัง
- **Multiple Reply Types** - ตอบได้หลายรูปแบบ (ข้อความ, รูป, ฯลฯ)

### ❌ **What We'll Skip (Phase 2)**
- ❌ AI Integration 
- ❌ Admin Chat Interface
- ❌ LINE/Shopee/TikTok

---

## 🏗️ **Architecture Overview**

```
Frontend (Next.js) - Keyword/Response Management
        ↓
ChatBotCore (8086) - Message Logic & Storage  
        ↓
FacebookConnect (8085) - Facebook API Integration
        ↓
PostgreSQL - Data Storage
```

---

## 📅 **Development Phases**

### **Phase 1: Service Separation & Foundation** (Week 1)

#### **Chunk 1.1: Create ChatBotCore Service**
- ✅ Setup ChatBotCore service structure
- ✅ Database schema for conversations & messages
- ✅ Basic REST API endpoints
- ✅ Repository pattern for chat data

#### **Chunk 1.2: Refactor FacebookConnect** 
- ✅ Extract Facebook API logic to separate service
- ✅ Clean webhook handling
- ✅ Message sending functionality
- ✅ Service-to-service communication

#### **Chunk 1.3: Integration Testing**
- ✅ Test FacebookConnect ↔ ChatBotCore communication
- ✅ End-to-end message flow
- ✅ Error handling & logging

### **Phase 2: Enhanced Features** (Week 2)

#### **Chunk 2.1: Multiple Reply Types**
- ✅ Text responses
- ✅ Image responses  
- ✅ Quick replies
- ✅ Template messages

#### **Chunk 2.2: Comment Monitoring**
- ✅ Facebook comment webhook
- ✅ Comment-to-inbox conversion
- ✅ Comment reply functionality

#### **Chunk 2.3: Chat History & Analytics**
- ✅ Conversation threading
- ✅ Message history storage
- ✅ Basic analytics endpoints

### **Phase 3: Frontend Integration** (Week 3)

#### **Chunk 3.1: Chat Management UI**
- ✅ Chat list view
- ✅ Chat history viewer
- ✅ Message timeline

#### **Chunk 3.2: Enhanced Keyword Management**
- ✅ Response templates
- ✅ Multiple response types
- ✅ Priority & conditions

#### **Chunk 3.3: Testing & Deployment**
- ✅ Full system testing
- ✅ Production deployment
- ✅ Performance optimization

---

## 📊 **Database Schema (MVP)**

```sql
-- Chat management schema
CREATE SCHEMA chatbot_mvp;

-- Conversations (Facebook users)
CREATE TABLE chatbot_mvp.conversations (
    id SERIAL PRIMARY KEY,
    facebook_user_id VARCHAR(255) UNIQUE NOT NULL,
    user_name VARCHAR(255),
    status VARCHAR(20) DEFAULT 'active', -- active, paused, closed
    last_message_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- Messages (all messages in conversations)
CREATE TABLE chatbot_mvp.messages (
    id SERIAL PRIMARY KEY,
    conversation_id INTEGER REFERENCES chatbot_mvp.conversations(id),
    facebook_message_id VARCHAR(255),
    sender_type VARCHAR(10) NOT NULL, -- 'user' or 'bot'
    message_type VARCHAR(20) DEFAULT 'text', -- text, image, quick_reply, template
    message_text TEXT,
    message_data JSONB, -- For rich message data
    matched_keyword VARCHAR(255),
    response_template VARCHAR(255),
    processing_time_ms INTEGER,
    created_at TIMESTAMP DEFAULT NOW()
);

-- Enhanced keywords with response types
CREATE TABLE chatbot_mvp.keywords (
    id SERIAL PRIMARY KEY,
    keyword TEXT NOT NULL,
    response_type VARCHAR(20) DEFAULT 'text', -- text, image, quick_reply, template
    response_content TEXT NOT NULL,
    response_data JSONB, -- Additional data for rich responses
    priority INTEGER DEFAULT 0,
    is_active BOOLEAN DEFAULT true,
    match_type VARCHAR(20) DEFAULT 'exact', -- exact, partial, regex
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- Response templates for reuse
CREATE TABLE chatbot_mvp.response_templates (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    template_type VARCHAR(20) NOT NULL, -- text, image, quick_reply, generic
    template_data JSONB NOT NULL,
    description TEXT,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT NOW()
);

-- Analytics (simple)
CREATE TABLE chatbot_mvp.conversation_stats (
    id SERIAL PRIMARY KEY,
    date DATE NOT NULL,
    total_conversations INTEGER DEFAULT 0,
    total_messages INTEGER DEFAULT 0,
    bot_responses INTEGER DEFAULT 0,
    avg_response_time_ms DECIMAL(10,2),
    created_at TIMESTAMP DEFAULT NOW()
);
```

---

## 🔌 **Service APIs**

### **ChatBotCore (Port 8086)**
```
POST /api/process-message     # Process incoming message
GET  /api/conversations       # List conversations  
GET  /api/conversations/:id   # Get conversation history
POST /api/keywords           # Manage keywords
GET  /api/keywords           # Get keywords
POST /api/response-templates # Manage templates
GET  /api/analytics          # Get stats
GET  /health                 # Health check
```

### **FacebookConnect (Port 8085)**
```
GET  /webhook/facebook       # Webhook verification
POST /webhook/facebook       # Webhook handler
POST /api/send-message       # Send message to Facebook
GET  /health                 # Health check
```

---

## 🎯 **Success Criteria**

### **Technical**
- ✅ Response time < 2 seconds
- ✅ Handle 50+ concurrent conversations
- ✅ 99% uptime
- ✅ All messages stored correctly

### **Functional** 
- ✅ Detect keywords accurately
- ✅ Send appropriate responses
- ✅ Store chat history
- ✅ Support multiple response types
- ✅ Comments converted to inbox messages

### **User Experience**
- ✅ Admins can view all chats
- ✅ Easy keyword management
- ✅ Real-time message updates
- ✅ Clean, intuitive interface

---

## 🚀 **Ready to Start**

**Next Step**: Begin Phase 1, Chunk 1.1 - Create ChatBotCore Service

ความพร้อม:
- ✅ Database connection established
- ✅ Repository pattern in place  
- ✅ Clean folder structure
- ✅ Development environment ready

**Estimated Timeline**: 3 weeks to fully functional MVP
**Target**: Replace existing zwiz functionality with improved architecture

---

*This MVP will provide a solid foundation for future AI integration while delivering immediate value for Facebook message management.*
