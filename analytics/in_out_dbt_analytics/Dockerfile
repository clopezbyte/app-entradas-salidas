# Use official dbt image
FROM ghcr.io/dbt-labs/dbt-bigquery:1.7.6

# Set working directory
WORKDIR /app

# Copy your dbt project into the container (including profiles.yml and keys)
COPY . .

# Move profiles.yml to the expected dbt location
RUN mkdir -p /root/.dbt && cp profiles.yml /root/.dbt/profiles.yml

# Add entrypoint script
COPY entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh

# Use bash script as entrypoint
ENTRYPOINT ["/entrypoint.sh"]