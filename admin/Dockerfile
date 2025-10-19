# Dockerfile for Django Admin
FROM python:3.12-slim

# Install system dependencies
RUN apt-get update && apt-get install -y \
    postgresql-client \
    && rm -rf /var/lib/apt/lists/*

# Set working directory
WORKDIR /app

# Copy dependency files
COPY pyproject.toml uv.lock ./

# Install uv and Python dependencies
RUN pip install --no-cache-dir uv && \
    uv pip install --system --no-cache django psycopg2-binary python-dotenv supabase gunicorn whitenoise

# Copy application code
COPY . .

# Make entrypoint script executable
RUN chmod +x docker-entrypoint.sh

# Collect static files
RUN python manage.py collectstatic --noinput || true

# Expose port
EXPOSE 8000

# Set production environment
ENV DJANGO_ENV=production
ENV PYTHONUNBUFFERED=1
ENV PORT=8000

# Run the startup script
CMD ["./docker-entrypoint.sh"]
