#!/bin/sh

# Print env
echo "Running dbt with:"
echo "  DBT_MODEL: $DBT_MODEL"
echo "  DBT_TARGET: $DBT_TARGET"
echo "  DBT_VARS: $DBT_VARS"

# Run dbt model
echo "Running dbt run for $DBT_MODEL , $DBT_TARGET ..."

if [ -z "$DBT_VARS" ]; then #if DBT_VARS is not set
    echo "Running: dbt run --select \"$DBT_MODEL\" --target \"$DBT_TARGET\""
    dbt run --select "$DBT_MODEL" --target "$DBT_TARGET"
else # if DBT_VARS is set
    echo "Running: dbt run --select \"$DBT_MODEL\" --target \"$DBT_TARGET\" --vars '$DBT_VARS'"
    dbt run --select "$DBT_MODEL" --target "$DBT_TARGET" --vars "$DBT_VARS"
fi

echo "dbt run completed"

# Run dbt test
echo "Running dbt test for $DBT_MODEL , $DBT_TARGET ..."

if [ -z "$DBT_VARS" ]; then
    dbt test --select "$DBT_MODEL" --target "$DBT_TARGET"
else
    dbt test --select "$DBT_MODEL" --target "$DBT_TARGET" --vars "$DBT_VARS"
fi

echo "dbt test completed"