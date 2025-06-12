# 🚀 AI ChatBot MVP Specification

## 📋 **MVP Overview**
สร้าง AI ChatBot แทน zwiz ที่ใช้งานได้จริงในระยะสั้น แต่วางโครงสร้างให้ขยายได้ในอนาคต

---

## 🎯 **MVP Scope (Phase 1 - ใช้งานได้ใน 2-3 สัปดาห์)**

### 📥 **Channel Integration**
✅ **Facebook Messenger เท่านั้น** (ใช้ของเดิมที่มี)
- Webhook รับข้อความ
- ส่งข้อความตอบกลับ
- รองรับทั้ง Comment และ Private Message

### 🧠 **AI Engine (Minimal)**
✅ **Keyword + AI Hybrid**
- **Rule-based**: Keywords สำหรับคำถามพื้นฐาน
- **AI-powered**: เรียกใช้ OpenAI API สำหรับคำถามซับซ้อน
- **Fallback**: ส่งต่อให้แอดมินเมื่อไม่รู้คำตอบ

### 💬 **Core Features**
✅ **ต้องมี**
- รับข้อความจาก Facebook
- ตอบด้วย AI หรือ Keywords
- เก็บ Chat History
- แอดมินดูแชทได้
- ปิด/เปิด auto-reply ได้

❌ **ไม่มีใน MVP**
- LINE Integration
- TikTok, Shopee, Lazada
- Advanced Analytics
- AI Training Interface
- Multi-agent support

### 🖥️ **Admin Dashboard (Simple)**
✅ **หน้าพื้นฐาน**
- ดู Chat List + Chat History
- เปิด/ปิด AI Bot
- จัดการ Keywords/Responses
- ตั้งค่า OpenAI API Key

---

## 🏗️ **MVP Architecture (Simplified)**

### 📁 **Directory Structure**
```
ruaymanagement/backend/
├── 🔌 external/
│   └── FacebookConnect/          # Port 8085
│       ├── handlers/webhook.go   # รับ webhook
│       └── services/api.go       # ส่งข้อความ
│
└── 🏢 internal/
    └── ChatBotCore/              # Port 8086 (All-in-one)
        ├── handlers/
        │   ├── chat.go           # Chat management
        │   ├── keywords.go       # Keyword management  
        │   └── ai.go             # AI processing
        ├── services/
        │   ├── chatbot.go        # Main chatbot logic
        │   ├── ai_engine.go      # OpenAI integration
        │   └── keyword_matcher.go # Keyword matching
        └── repository/
            ├── chat_repo.go      # Chat history
            └── keyword_repo.go   # Keywords
```

### 🔄 **Simplified Data Flow**
```
Facebook → FacebookConnect → ChatBotCore → FacebookConnect → Facebook
           (Webhook)        (AI + Keywords)   (Send Reply)
```

### 📊 **MVP Database Schema**
```sql
-- Single schema for MVP
CREATE SCHEMA chatbot_mvp;

-- Chat conversations
CREATE TABLE chatbot_mvp.conversations (
    id SERIAL PRIMARY KEY,
    facebook_user_id VARCHAR(255) NOT NULL,
    user_name VARCHAR(255),
    status VARCHAR(20) DEFAULT 'active', -- active, closed, admin_mode
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- Chat messages
CREATE TABLE chatbot_mvp.messages (
    id SERIAL PRIMARY KEY,
    conversation_id INTEGER REFERENCES chatbot_mvp.conversations(id),
    message_text TEXT NOT NULL,
    sender_type VARCHAR(10) NOT NULL, -- 'user', 'bot', 'admin'
    message_type VARCHAR(20) DEFAULT 'text', -- text, image, quick_reply
    ai_used BOOLEAN DEFAULT false,
    matched_keyword VARCHAR(255),
    processing_time_ms INTEGER,
    facebook_message_id VARCHAR(255),
    created_at TIMESTAMP DEFAULT NOW()
);

-- Keywords (simple)
CREATE TABLE chatbot_mvp.keywords (
    id SERIAL PRIMARY KEY,
    keyword TEXT NOT NULL,
    response TEXT NOT NULL,
    is_active BOOLEAN DEFAULT true,
    priority INTEGER DEFAULT 0,
    created_at TIMESTAMP DEFAULT NOW()
);

-- Bot settings
CREATE TABLE chatbot_mvp.bot_settings (
    id SERIAL PRIMARY KEY,
    setting_key VARCHAR(100) UNIQUE NOT NULL,
    setting_value TEXT,
    updated_at TIMESTAMP DEFAULT NOW()
);

-- Initial settings
INSERT INTO chatbot_mvp.bot_settings (setting_key, setting_value) VALUES
('ai_enabled', 'true'),
('openai_api_key', ''),
('openai_model', 'gpt-3.5-turbo'),
('fallback_to_admin', 'true'),
('auto_reply_enabled', 'true');
```

---

## 🔧 **MVP Implementation Plan**

### **Week 1: Core Infrastructure**
1. ✅ **Setup ChatBotCore service** (Port 8086)
2. ✅ **Database schema** และ repositories
3. ✅ **Basic AI integration** (OpenAI API)
4. ✅ **Keyword matching** system

### **Week 2: Facebook Integration**
1. ✅ **Refactor FacebookConnect** (Port 8085)
2. ✅ **Connect FacebookConnect ↔ ChatBotCore**
3. ✅ **Chat history** และ conversation management
4. ✅ **Basic testing**

### **Week 3: Admin Dashboard**
1. ✅ **Simple dashboard** (Next.js)
2. ✅ **Chat list** และ chat history viewer
3. ✅ **Bot settings** page
4. ✅ **Keywords management**

---

## 🎛️ **MVP Configuration**

### **Environment Variables**
```env
# Facebook (existing)
FACEBOOK_PAGE_TOKEN=xxx
FACEBOOK_VERIFY_TOKEN=xxx

# OpenAI
OPENAI_API_KEY=xxx
OPENAI_MODEL=gpt-3.5-turbo
OPENAI_MAX_TOKENS=150

# Database (existing)
DB_HOST=postgres
DB_PORT=5432
DB_NAME=loyverse_cache

# Services
FACEBOOK_CONNECT_PORT=8085
CHATBOT_CORE_PORT=8086
```

### **Docker Compose (MVP)**
```yaml
version: '3.8'

services:
  postgres:
    image: postgres:13
    ports:
      - "5433:5432"
    
  facebook-connect:
    build: ./backend/external/FacebookConnect
    ports:
      - "8085:8085"
    depends_on:
      - postgres
      - chatbot-core
      
  chatbot-core:
    build: ./backend/internal/ChatBotCore
    ports:
      - "8086:8086"
    depends_on:
      - postgres
      
  frontend:
    build: ./frontend-next
    ports:
      - "3000:3000"
    depends_on:
      - chatbot-core
```

---

## 🧠 **MVP AI Logic**

### **Decision Tree**
```
Incoming Message
       ↓
   Check Keywords
       ↓
   ✅ Match → Send predefined response
       ↓
   ❌ No Match → Call OpenAI API
       ↓
   ✅ AI Response → Send AI response
       ↓
   ❌ AI Failed → Send fallback message
       ↓
   📝 Log everything
```

### **OpenAI Integration (Simple)**
```go
type AIRequest struct {
    Message string `json:"message"`
    Context string `json:"context,omitempty"`
}

type AIResponse struct {
    Response    string  `json:"response"`
    Confidence  float64 `json:"confidence"`
    ShouldReply bool    `json:"should_reply"`
}

// Simple prompt template
func buildPrompt(message, context string) string {
    return fmt.Sprintf(`
คุณเป็นผู้ช่วยลูกค้าของร้านอาหาร
คำถาม: %s
บริบท: %s

ตอบแบบสั้นๆ เป็นกันเอง ถ้าไม่แน่ใจให้บอกว่าจะส่งต่อให้แอดมิน
`, message, context)
}
```

---

## 📱 **MVP User Stories**

### **ลูกค้า (Facebook User)**
1. ลูกค้าส่งข้อความ "หมูปิ้งราคาเท่าไหร่"
2. บอทตอบ "หมูปิ้งไม้ละ 15 บาท ค่ะ สั่งกี่ไม้ดีคะ?"
3. ลูกค้าถาม "ส่งได้มั้ย"
4. บอทตอบผ่าน AI "ส่งได้ค่ะ ในรัศมี 5 กม. ค่าส่ง 20 บาท"
5. ลูกค้าถาม "วันนี้ปิดกี่โมง"
6. บอทตอบผ่าน keyword "เปิดทุกวัน 10:00-21:00 ค่ะ"

### **แอดมิน (Dashboard)**
1. เข้า Dashboard เห็น chat list
2. คลิกดู chat history ของลูกค้า
3. เห็นว่าบอทตอบอะไรไปบ้าง
4. ปิด auto-reply สำหรับลูกค้าคนนี้ได้
5. เพิ่ม keyword ใหม่ได้

---

## 🎯 **MVP Success Metrics**

### **Technical**
- ✅ Response time < 3 วินาที
- ✅ 99% uptime
- ✅ Handle 100 messages/day

### **Business**
- ✅ ลดเวลาตอบลูกค้า 80%
- ✅ ตอบได้ 70% ของคำถามพื้นฐาน
- ✅ แอดมินใช้งาน dashboard ได้

---

## 🚀 **Future Roadmap (Post-MVP)**

### **Phase 2** (Month 2-3)
- LINE Integration
- Advanced Analytics
- AI Training Interface
- Multi-language support

### **Phase 3** (Month 4-6)
- TikTok, Shopee Integration
- Advanced AI (RAG, Fine-tuning)
- A/B Testing
- Performance optimization

### **Phase 4** (Month 6+)
- Voice messages
- Image recognition
- E-commerce integration
- Multi-tenant support

---

## ✅ **MVP Acceptance Criteria**

1. **Facebook Integration**: รับส่งข้อความได้
2. **AI Response**: ตอบคำถามผ่าน OpenAI ได้
3. **Keyword System**: ตอบคำถามพื้นฐานได้
4. **Chat History**: เก็บและดูประวัติได้
5. **Admin Dashboard**: จัดการพื้นฐานได้
6. **Deployment**: Deploy ใช้งานจริงได้

---

*MVP นี้จะให้คุณมีระบบ AI ChatBot ที่ใช้งานได้จริงแทน zwiz ในเวลาไม่เกิน 3 สัปดาห์ พร้อมโครงสร้างที่ขยายได้ในอนาคต*
