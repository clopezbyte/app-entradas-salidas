from google.cloud import bigquery

class LoadToBigQuery:
    def __init__(self, project_id:str):
        self.project_id = project_id
        self.bigquery_client = bigquery.Client(project=project_id)

    def load_dataframe_to_bigquery(self, dataframe, dataset_name: str, table_name: str):
        dataset_ref = self.bigquery_client.dataset(dataset_name)
        table_ref = dataset_ref.table(table_name)

        job_config = bigquery.LoadJobConfig(
            write_disposition=bigquery.WriteDisposition.WRITE_TRUNCATE,
            source_format=bigquery.SourceFormat.PARQUET,
        )

        load_job = self.bigquery_client.load_table_from_dataframe(
            dataframe, table_ref, job_config=job_config
        )

        load_job.result()