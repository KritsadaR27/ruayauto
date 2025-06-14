# Development Dockerfile for Next.js
FROM node:20-alpine

# Install dependencies for development
RUN apk add --no-cache libc6-compat

WORKDIR /app

# Copy package files
COPY web/package.json web/package-lock.json* ./

# Install all dependencies (including dev dependencies)
RUN npm ci

# Expose the development port
EXPOSE 3008

ENV PORT 3008
ENV HOSTNAME "0.0.0.0"
ENV NODE_ENV development
ENV NEXT_TELEMETRY_DISABLED 1

# Start the development server with hot reload
CMD ["npm", "run", "dev"]
