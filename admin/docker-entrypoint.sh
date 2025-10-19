#!/bin/bash
set -e

echo "Starting Django admin deployment..."

# Wait for database to be ready
echo "Waiting for database to be ready..."
until python -c "import psycopg2; psycopg2.connect(host='${DB_HOST}', port='${DB_PORT}', user='${DB_USER}', password='${DB_PASSWORD}', dbname='${DB_NAME}', sslmode='require')" 2>/dev/null; do
  echo "Database is unavailable - sleeping"
  sleep 2
done

echo "Database is ready!"

# Run migrations
echo "Running database migrations..."
python manage.py migrate --noinput

# Create superuser if credentials are provided
if [ -n "$DJANGO_SUPERUSER_USERNAME" ] && [ -n "$DJANGO_SUPERUSER_PASSWORD" ]; then
  echo "Creating superuser..."
  python manage.py shell <<EOF
from django.contrib.auth.models import User
import os

username = os.environ.get('DJANGO_SUPERUSER_USERNAME')
password = os.environ.get('DJANGO_SUPERUSER_PASSWORD')
email = os.environ.get('DJANGO_SUPERUSER_EMAIL', '')

if not User.objects.filter(username=username).exists():
    User.objects.create_superuser(username, email, password)
    print(f'Superuser {username} created successfully')
else:
    print(f'Superuser {username} already exists')
EOF
else
  echo "Skipping superuser creation - no credentials provided"
fi

# Note: We do NOT run seed_db in production

# Start Django with gunicorn (production WSGI server)
echo "Starting Django with gunicorn on port ${PORT:-8000}..."
exec gunicorn farmsville.wsgi:application \
    --bind 0.0.0.0:${PORT:-8000} \
    --workers 1 \
    --threads 2 \
    --timeout 60 \
    --access-logfile - \
    --error-logfile -
