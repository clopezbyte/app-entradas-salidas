#!/bin/sh

# Print env
echo "Running dbt with:"
echo "  DBT_MODEL: $DBT_MODEL"
echo "  DBT_TARGET: $DBT_TARGET"
echo "  DBT_VARS: $DBT_VARS"

# Build the vars for overrides #DBT_VARS='{ "backfill_month": "2024-06-01" }' #This is supposed to be passed manually 
# as env vars override
if [ -z "$DBT_VARS" ]; then
    VARS_OPTION=""
else
    VARS_OPTION="--vars $DBT_VARS"
fi


# Run dbt model
echo "Running dbt run for $DBT_MODEL , $DBT_TARGET ..."
dbt run --select "$DBT_MODEL" --target "$DBT_TARGET" $VARS_OPTION
echo "dbt run completed"

# Run dbt test
echo "Running dbt test for $DBT_MODEL , $DBT_TARGET ..."
dbt test --select "$DBT_MODEL" --target "$DBT_TARGET" $VARS_OPTION
echo "dbt test completed"
