# Farmsville

Neighobrs claim items. Enables record keeping and organization for small scale local famers.

## Architecture

The project uses a dual-application architecture:

- **Django Admin** (port 8000): Content management system for managing orders
- **Go Web Application** (port 3000): Customer-facing interface for browsing and claiming items
- **PostgreSQL**
- **Tailwind CSS**

## Quick Start

### Prerequisites

- Docker and Docker Compose
- Go
- Node.js (for Tailwind CSS)
- Python with uv (for Django admin)

### Development Setup

The project uses a Makefile for development. Simply run:

```bash
make run
```

This command will:

- Start PostgreSQL in Docker Compose
- Run database migrations
- Create seed data for testing
- Start Django admin on `http://localhost:8000` (credentials: admin/admin)
- Build and watch Tailwind CSS
- Start the Go web application on `http://localhost:3000`

To stop all services:

```bash
make stop
```

## Contributing

Contributions are welcome!
