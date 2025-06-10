
{{ config(
    materialized ='incremental', 
    unique_key = 'silver_movement_id',
    on_schema_change = 'fail',
    description ='This model transforms the bronze in and out of warehouse data and appends into a silver format, calculating metrics for analytics dashboard.',
    partition_by = {
        "field": "fecha_movimiento",
        "data_type": "timestamp",
        "granularity": "month"
    },
    cluster_by=["tipo","bodega","cliente"]
) }}

WITH new_in_out_data AS(
    SELECT
        landing_movement_id as silver_movement_id,
        CASE
            WHEN bodega IS NULL THEN 'UNKNOWN' 
            WHEN bodega = 'Bodega MTY Santa Catarina' THEN 'Bodega MTY'
            ELSE bodega
        END AS bodega,
        COALESCE(cantidad, 0) AS cantidad, 
        UPPER(
            CASE
                WHEN cliente IS NULL THEN 'UNKNOWN' 
                WHEN cliente = 'NA' THEN 'N/A'
                ELSE cliente
            END 
        ) AS cliente,
        fecha_movimiento,
        fecha_ajuste_asn,
        CASE
            WHEN tipo_delivery IS NULL THEN 'N/A' 
            ELSE tipo_delivery
        END AS tipo_delivery,
        CASE 
            WHEN operador IS NULL THEN 'UNKNOWN'
            WHEN operador = 'NA' THEN 'N/A'
            ELSE operador
        END AS operador,
        UPPER(
            CASE 
                WHEN proveedor IS NULL THEN 'UNKNOWN'
                WHEN proveedor = 'NA' THEN 'N/A'
                ELSE proveedor
            END 
        )AS proveedor,
        UPPER(COALESCE(tipo, 'UNKNOWN')) AS tipo 
    FROM 
        `b-materials.in_out_bronze.landing_in_out_movements`
    WHERE
        -- TIMESTAMP_TRUNC(fecha_movimiento, MONTH) = TIMESTAMP("2025-05-01") -- For backfilling purposes (try to do dynamically)
        TIMESTAMP_TRUNC(fecha_movimiento, MONTH) = TIMESTAMP_TRUNC(CURRENT_TIMESTAMP(), MONTH)
)


SELECT 
    *
FROM 
    new_in_out_data
{% if is_incremental() %}
WHERE 
    silver_movement_id NOT IN (
        SELECT silver_movement_id FROM {{ this }}
)
{% endif %}