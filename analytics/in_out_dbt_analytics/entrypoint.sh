#!/bin/sh

# Print env for debugging
echo "Running dbt with:"
echo "  DBT_MODEL: $DBT_MODEL"
echo "  DBT_TARGET: $DBT_TARGET"

# Test dbt model with env vars
echo "Running dbt Test for $DBT_MODEL , $DBT_TARGET ..."
dbt test --select "$DBT_MODEL" --target "$DBT_TARGET"
echo "dbt test completed"

# Run dbt model with env vars
echo "Running dbt run for $DBT_MODEL , $DBT_TARGET ..."
dbt run --select "$DBT_MODEL" --target "$DBT_TARGET"
echo "dbt run completed"
