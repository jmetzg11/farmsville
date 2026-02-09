# Build stage for CSS
FROM node:20-alpine AS css-builder

WORKDIR /app/web

# Copy package files and install dependencies
COPY web/package*.json ./
RUN npm install

# Copy UI files and build Tailwind CSS
COPY web/ui ./ui
COPY web/tailwind.config.js ./
RUN npm run build:css

# Build stage for Go
FROM golang:1.25-alpine AS go-builder

WORKDIR /app/web

# Copy go mod files
COPY web/go.mod web/go.sum ./
RUN go mod download

# Copy source code
COPY web/*.go ./

# Copy UI files including the built CSS from css-builder
COPY --from=css-builder /app/web/ui ./ui

# Build the Go binary with embedded files
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

# Runtime stage
FROM alpine:latest

WORKDIR /app

RUN apk add --no-cache tzdata

# Copy only the binary (UI files are embedded in it)
COPY --from=go-builder /app/web/main .

# Expose port 3000 (matching fly.toml)
EXPOSE 3000

# Run the application with -prod flag
CMD ["./main", "-prod"]