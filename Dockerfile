FROM node:alpine AS frontend-builder
WORKDIR /app/frontend
# Copy package files first to leverage caching
COPY frontend/package*.json ./
RUN npm install
# Only then copy source files (this helps with caching when only source changes)
COPY frontend/src ./src
# COPY frontend/public ./public
# Copy only necessary config files
COPY frontend/*.js frontend/*.json ./
RUN npm run build

FROM golang:1.23-alpine AS go-builder
WORKDIR /go/src/farmsville
# Install build dependencies
RUN apk add --no-cache gcc musl-dev
# Copy go.mod and go.sum first to leverage Docker cache
COPY go.mod go.sum ./
RUN go mod download
# Copy only what's needed for building
COPY main.go ./
COPY backend/ ./backend/
# Build the Go application
RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -o main .

# Final stage - use a lightweight Alpine image
FROM alpine:latest
WORKDIR /app
# Install SQLite
RUN apk add --no-cache sqlite

# Copy only the compiled Go binary from the go-builder stage
COPY --from=go-builder /go/src/farmsville/main /app/main
# Copy only the frontend build from the frontend-builder stage
COPY --from=frontend-builder /app/frontend/build /app/frontend/build

# Expose the port
EXPOSE 3000

# Command to run the application
CMD ["/app/main"]
