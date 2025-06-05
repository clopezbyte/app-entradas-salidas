

{{ config(
    materialized='incremental', 
    description='This model transforms the bronze in and out of warehouse data and appends into a silver format, calculating metrics for analytics dashboard.',
    partition_by={
        "field": "fecha_movimiento",
        "data_type": "timestamp",
        "granularity": "month"
    },
    cluster_by=["tipo","bodega","cliente"]
)}}

