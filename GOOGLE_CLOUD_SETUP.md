# Google Cloud Storage Setup for RuayChatBot

## 1. Setup Google Cloud Project

```bash
# Install Google Cloud CLI
curl https://sdk.cloud.google.com | bash
exec -l $SHELL

# Login and create project
gcloud auth login
gcloud projects create ruay-chatbot --name="RuayChatBot"
gcloud config set project ruay-chatbot

# Enable required APIs
gcloud services enable storage-json.googleapis.com
gcloud services enable storage-component.googleapis.com
```

## 2. Create Storage Bucket

```bash
# Create bucket (must be globally unique)
gsutil mb gs://ruay-chatbot-media

# Set public access for uploaded files
gsutil iam ch allUsers:objectViewer gs://ruay-chatbot-media

# Set CORS for web access
echo '[
    {
      "origin": ["*"],
      "method": ["GET", "POST", "PUT", "DELETE"],
      "responseHeader": ["Content-Type"],
      "maxAgeSeconds": 3600
    }
]' > cors.json

gsutil cors set cors.json gs://ruay-chatbot-media
```

## 3. Create Service Account

```bash
# Create service account
gcloud iam service-accounts create ruay-chatbot-storage \
    --description="Service account for RuayChatBot file uploads" \
    --display-name="RuayChatBot Storage"

# Grant storage permissions
gcloud projects add-iam-policy-binding ruay-chatbot \
    --member="serviceAccount:ruay-chatbot-storage@ruay-chatbot.iam.gserviceaccount.com" \
    --role="roles/storage.objectAdmin"

# Create and download service account key
gcloud iam service-accounts keys create ./gcloud-service-key.json \
    --iam-account=ruay-chatbot-storage@ruay-chatbot.iam.gserviceaccount.com
```

## 4. Environment Variables

Add to your `.env.local`:

```bash
GOOGLE_CLOUD_PROJECT_ID=ruay-chatbot
GOOGLE_CLOUD_BUCKET_NAME=ruay-chatbot-media
GOOGLE_APPLICATION_CREDENTIALS=./gcloud-service-key.json
```

## 5. Install Dependencies

```bash
npm install @google-cloud/storage
```

## 6. Usage Example

The media upload API will automatically:
- Upload files to Google Cloud Storage
- Generate public URLs
- Store metadata in media library
- Handle file type validation
- Manage file size limits

## 7. File Organization

Files will be organized as:
```
gs://ruay-chatbot-media/
  uploads/
    rule-{ruleId}/
      {timestamp}-{filename}
  avatars/
    {userId}/
      {filename}
  temp/
    {sessionId}/
      {filename}
```

## 8. Security Notes

- Service account key should be kept secure
- Consider using Workload Identity in production
- Set up proper IAM roles for different environments
- Enable logging for audit trail
