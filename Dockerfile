# Use the official PostgreSQL image as the base image
FROM postgres:15

# Install postgresql-contrib (needed for uuid-ossp extension)
RUN apt-get update && apt-get install -y postgresql-contrib

# Set environment variables for PostgreSQL
ENV POSTGRES_USER=chirp_user
ENV POSTGRES_PASSWORD=chirp_password
ENV POSTGRES_DB=chirp_db

# Expose the default PostgreSQL port
EXPOSE 5432

# Add a health check to ensure the database is ready
HEALTHCHECK --interval=10s --timeout=5s --start-period=30s --retries=3 \
  CMD pg_isready -U $POSTGRES_USER -d $POSTGRES_DB
