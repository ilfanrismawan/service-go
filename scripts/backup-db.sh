#!/bin/bash
# Automated database backup script for PostgreSQL
# Run this as a cron job (e.g., daily at 2 AM)

set -e

# Configuration
BACKUP_DIR="${BACKUP_DIR:-/var/backups/postgresql}"
DB_NAME="${DB_NAME:-iphone_service}"
DB_HOST="${DB_HOST:-localhost}"
DB_PORT="${DB_PORT:-5432}"
DB_USER="${DB_USER:-postgres}"
RETENTION_DAYS="${RETENTION_DAYS:-30}"

# Create backup directory if it doesn't exist
mkdir -p "$BACKUP_DIR"

# Generate backup filename with timestamp
TIMESTAMP=$(date +%Y%m%d_%H%M%S)
BACKUP_FILE="$BACKUP_DIR/${DB_NAME}_${TIMESTAMP}.sql"

# Perform backup
echo "Starting database backup: $BACKUP_FILE"
PGPASSWORD="$DB_PASSWORD" pg_dump -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$DB_NAME" -F c -f "$BACKUP_FILE.dump" 2>/dev/null || \
PGPASSWORD="$DB_PASSWORD" pg_dump -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$DB_NAME" > "$BACKUP_FILE"

if [ $? -eq 0 ]; then
    echo "✅ Backup completed successfully: $BACKUP_FILE"
    
    # Compress backup
    gzip "$BACKUP_FILE" 2>/dev/null || true
    
    # Remove old backups (older than retention days)
    find "$BACKUP_DIR" -name "${DB_NAME}_*.sql*" -type f -mtime +$RETENTION_DAYS -delete
    
    echo "✅ Old backups cleaned (keeping last $RETENTION_DAYS days)"
else
    echo "❌ Backup failed!"
    exit 1
fi

