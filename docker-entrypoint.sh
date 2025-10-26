#!/bin/sh

# Wait for PostgreSQL to be ready
echo "Waiting for PostgreSQL..."
until PGPASSWORD=$DATABASE_PASSWORD psql -h $DATABASE_HOST -U $DATABASE_USER -d $DATABASE_DBNAME -c '\q' 2>/dev/null; do
  echo "PostgreSQL is unavailable - sleeping"
  sleep 1
done

echo "PostgreSQL is up - running migrations..."

# Function to run migration (extract content between -- +migrate Up and -- +migrate Down)
run_migration() {
    local file=$1
    echo "Processing migration: $(basename $file)"
    
    # Simple approach: Use Python to extract between markers if available, otherwise use grep
    # For now, just execute the file directly - the "Up" parts will run, "Down" parts will be ignored by grep comments
    
    # Extract only lines between "-- +migrate Up" and "-- +migrate Down"  
    # Using grep with context (but this needs more complex logic)
    # Simpler: Use a Python one-liner if available, or just grep for non-comment non-Down lines
    awk '
    /-- \+migrate Up/ {start=1; next}
    /-- \+migrate Down/ {start=0}
    start && !/-- \+migrate Down/
    ' "$file" > /tmp/migration_up.sql
    
    if [ -s /tmp/migration_up.sql ]; then
        echo "Running migration: $(basename $file)"
        PGPASSWORD=$DATABASE_PASSWORD psql -h $DATABASE_HOST -U $DATABASE_USER -d $DATABASE_DBNAME -f /tmp/migration_up.sql
        if [ $? -eq 0 ]; then
            echo "Migration $(basename $file) completed successfully"
        else
            echo "Migration $(basename $file) failed"
        fi
        rm -f /tmp/migration_up.sql
    else
        echo "Migration $(basename $file) is empty or invalid"
    fi
}

# Run all migrations in order
if [ -d "/root/db/migrations" ]; then
    for migration in $(ls /root/db/migrations/*.sql 2>/dev/null | sort); do
        run_migration "$migration"
    done
    echo "All migrations completed"
else
    echo "No migrations directory found"
fi

# Run the application
echo "Starting application..."
exec ./gopilot
