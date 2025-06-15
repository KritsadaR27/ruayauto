# 🔧 Facebook App Configuration Guide

## 📋 Required App Setup

### 1. Facebook App Permissions
ใน Meta Developers Console (https://developers.facebook.com) ตั้งค่า:

**Basic Permissions:**
- `pages_show_list` - เพื่อดูรายการ Pages ที่ user admin
- `pages_read_engagement` - เพื่ออ่านการมีส่วนร่วมบน Pages
- `pages_manage_metadata` - เพื่อจัดการ metadata ของ Pages

**❌ ไม่ต้องใช้ (deprecated):**
- `pages_manage_comments` - Scope นี้ deprecated แล้ว

### 2. Valid OAuth Redirect URIs
เพิ่ม URLs เหล่านี้ใน Facebook App Settings > Basic > Valid OAuth Redirect URIs:

**Development (IMPORTANT):**
```
http://localhost:3008/api/auth/facebook/callback
```

**Production:**
```
https://yourdomain.com/api/auth/facebook/callback
```

**⚠️ สำคัญ:** Facebook App ต้อง:
1. อยู่ใน **Development Mode** สำหรับ localhost testing
2. เพิ่ม **Test Users** หรือ **App Roles** สำหรับ developers
3. **ไม่ต้อง HTTPS** สำหรับ localhost development

### 3. Webhook Configuration

**Webhook URL:**
```
https://yourdomain.com/webhook/facebook
```

**Verify Token:**
```
${FACEBOOK_VERIFY_TOKEN}
```

**Subscribed Fields:**
- `comments`
- `feed`
- `posts`

### 4. Page Access Token
หลังจาก OAuth success แล้ว:
1. เลือก Page ที่ต้องการ
2. ระบบจะขอ Page Access Token อัตโนมัติ
3. Token จะถูกเก็บใน database

## 🔧 Environment Variables Required

```bash
# Facebook App Configuration
FACEBOOK_APP_ID=your_app_id_here
NEXT_PUBLIC_FACEBOOK_APP_ID=your_app_id_here
FACEBOOK_APP_SECRET=your_app_secret_here
FACEBOOK_VERIFY_TOKEN=your_verify_token_here

# OAuth Configuration
NEXTAUTH_URL=http://localhost:3008
NEXTAUTH_SECRET=your_secret_here
```

## 🚨 Common Issues & Solutions

### Issue 1: "Invalid Scope" Error
**Solution:** ใช้ scope ที่ถูกต้อง:
```javascript
scope: 'pages_show_list,pages_read_engagement,pages_manage_metadata'
```

### Issue 2: "App not approved" Error
**Solution:** 
1. ไปที่ App Review ใน Meta Developers Console
2. Submit สำหรับ permissions ที่ต้องการ
3. อธิบายการใช้งานอย่างชัดเจน

### Issue 3: "Invalid Redirect URI" Error
**Solution:**
1. ตรวจสอบ Valid OAuth Redirect URIs ใน App Settings
2. ต้องตรงกับ `NEXTAUTH_URL/api/auth/facebook/callback`

### Issue 4: Webhook Verification Failed
**Solution:**
1. ตรวจสอบ HTTPS certificate (production)
2. ตรวจสอบ FACEBOOK_VERIFY_TOKEN ตรงกัน
3. ตรวจสอบ webhook endpoint response ถูกต้อง

## 📝 Testing Steps

1. **Test OAuth Flow:**
   ```bash
   curl -X GET "http://localhost:3008"
   # คลิก "Connect with Facebook" button
   ```

2. **Test Webhook:**
   ```bash
   curl -X GET "http://localhost:8091/webhook/facebook?hub.mode=subscribe&hub.challenge=test&hub.verify_token=${FACEBOOK_VERIFY_TOKEN}"
   ```

3. **Check Page Connection:**
   ```bash
   curl -X GET "http://localhost:3008/api/facebook/pages"
   ```

## 🔗 Useful Links

- [Meta Developers Console](https://developers.facebook.com)
- [Facebook API Documentation](https://developers.facebook.com/docs/graph-api)
- [Webhook Setup Guide](https://developers.facebook.com/docs/graph-api/webhooks)
- [Page Access Tokens](https://developers.facebook.com/docs/pages/access-tokens)
