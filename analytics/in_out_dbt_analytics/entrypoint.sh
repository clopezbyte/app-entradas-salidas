#!/bin/sh

# Print env
echo "Running dbt with:"
echo "  DBT_MODEL: $DBT_MODEL"
echo "  DBT_TARGET: $DBT_TARGET"


# Run dbt model with env vars
echo "Running dbt run for $DBT_MODEL , $DBT_TARGET ..."
dbt run --select "$DBT_MODEL" --target "$DBT_TARGET"
echo "dbt run completed"

# Test dbt model with env vars
echo "Running dbt Test for $DBT_MODEL , $DBT_TARGET ..."
dbt test --select "$DBT_MODEL" --target "$DBT_TARGET"
echo "dbt test completed"