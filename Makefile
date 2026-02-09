start-db:
	cd data && docker compose up -d postgres && echo "Waiting for PostgreSQL to be ready..." && until docker compose exec postgres pg_isready -U admin -d farmsville; do sleep 1; done && echo "PostgreSQL is ready!"

run: start-db
	@echo "Starting Django admin..."
	cd admin && DB_HOST=localhost DB_PORT=5432 DB_NAME=farmsville DB_USER=admin DB_PASSWORD=admin DEBUG=True uv run python manage.py migrate && \
	DB_HOST=localhost DB_PORT=5432 DB_NAME=farmsville DB_USER=admin DB_PASSWORD=admin DEBUG=True uv run python manage.py shell -c "from django.contrib.auth.models import User; import os; User.objects.create_superuser(os.environ.get('ADMIN_USERNAME', 'admin'), '', os.environ.get('ADMIN_PASSWORD', 'admin')) if not User.objects.filter(username=os.environ.get('ADMIN_USERNAME', 'admin')).exists() else print('Superuser already exists')" && \
	DB_HOST=localhost DB_PORT=5432 DB_NAME=farmsville DB_USER=admin DB_PASSWORD=admin DEBUG=True uv run python manage.py seed_db && \
	DB_HOST=localhost DB_PORT=5432 DB_NAME=farmsville DB_USER=admin DB_PASSWORD=admin DEBUG=True uv run python manage.py runserver &
	@echo "Building Tailwind CSS..."
	cd web && npm run build:css
	@echo "Starting Tailwind CSS watcher..."
	cd web && npm run watch:css &
	@echo "Starting Go backend..."
	cd web && air

stop:
	@echo "Stopping Django admin..."
	-pkill -f "manage.py runserver"
	@echo "Stopping Go backend..."
	-pkill -f "tmp/main"
	@echo "Cleaning up Go build artifacts..."
	-rm -f web/tmp/main
	@echo "Stopping Tailwind CSS watcher..."
	-pkill -f "tailwindcss.*--watch"
	@echo "Stopping PostgreSQL..."
	cd data && docker compose down -v
