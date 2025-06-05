from google.cloud import bigquery
import pandas as pd
import logging
import sys

# Configure Logging
logging.basicConfig(level=logging.INFO, format='%(asctime)s - %(levelname)s - %(message)s', stream=sys.stdout)
logger = logging.getLogger('in-out-analytics')

class LoadToBigQuery:
    def __init__(self, project_id:str):
        self.project_id = project_id
        self.bigquery_client = bigquery.Client(project=project_id)

    def load_dataframe_to_bigquery(self, dataframe: pd.DataFrame, dataset_name: str, table_name: str) -> None:
        dataset_ref = self.bigquery_client.dataset(dataset_name)
        table_ref = dataset_ref.table(table_name)

        job_config = bigquery.LoadJobConfig(
            # write_disposition=bigquery.WriteDisposition.WRITE_TRUNCATE, #leave as default to APPEND
            schema=[
                bigquery.SchemaField("landing_movement_id", bigquery.enums.SqlTypeNames.STRING, mode="REQUIRED"),
                bigquery.SchemaField("bodega", bigquery.enums.SqlTypeNames.STRING, mode="NULLABLE"),
                bigquery.SchemaField("cantidad", bigquery.enums.SqlTypeNames.INTEGER, mode="NULLABLE"),
                bigquery.SchemaField("cliente", bigquery.enums.SqlTypeNames.STRING, mode="NULLABLE"),
                bigquery.SchemaField("fecha_movimiento", bigquery.enums.SqlTypeNames.TIMESTAMP, mode="NULLABLE"),
                bigquery.SchemaField("fecha_ajuste_asn", bigquery.enums.SqlTypeNames.TIMESTAMP, mode="NULLABLE"),
                bigquery.SchemaField("tipo_delivery", bigquery.enums.SqlTypeNames.STRING, mode="NULLABLE"),
                bigquery.SchemaField("operador", bigquery.enums.SqlTypeNames.STRING, mode="NULLABLE"),
                bigquery.SchemaField("proveedor", bigquery.enums.SqlTypeNames.STRING, mode="NULLABLE"),
                bigquery.SchemaField("tipo", bigquery.enums.SqlTypeNames.STRING, mode="NULLABLE"),
            ]
        )

        try:
            if not dataframe.empty:
                load_job = self.bigquery_client.load_table_from_dataframe(
                    dataframe, table_ref, job_config=job_config
                )
                load_job.result()
                logger.info(f"BigQuery load job ID: {load_job.job_id}")
                logger.info(f"Loaded {load_job.output_rows} rows into {dataset_name}.{table_name}.")
            else:
                logger.warning("No data to load. DataFrame is empty.")
                return 
        except Exception as e:
            logger.exception(f"Failed to load data to BigQuery: {e}")
            raise
