# Facebook Webhook Keyword Reply (Go + Gin)

## Overview
This project is a Go backend using Gin that listens for Facebook Webhook events (feed.comment) and replies to comments based on keyword matching.

## Setup
1. Copy `.env` and fill in your `FB_PAGE_TOKEN` and `FB_VERIFY_TOKEN`.
2. Build and run with Docker Compose:
   ```sh
   docker-compose up --build
   ```
3. The server will listen on port 8100.

## Development
- POST route: `/webhook/facebook`
- Add your keywords and replies in `utils/reply.go`.

## No database or frontend included.
