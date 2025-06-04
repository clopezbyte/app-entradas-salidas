#!/bin/sh

# Print env for debugging
echo "Running dbt with:"
echo "  DBT_MODEL: $DBT_MODEL"
echo "  DBT_TARGET: $DBT_TARGET"

# Run dbt with env vars
dbt run --select "$DBT_MODEL" --target "$DBT_TARGET"
