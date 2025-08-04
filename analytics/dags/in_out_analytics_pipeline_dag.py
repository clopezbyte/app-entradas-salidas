from airflow import DAG
from airflow.operators.empty import EmptyOperator
from airflow.operators.python import PythonOperator
from datetime import datetime, timedelta
from airflow.providers.google.cloud.operators.cloud_run import CloudRunExecuteJobOperator


default_args = {
    'owner': 'clopez',
    'start_date': datetime(2025, 8, 1),
    'retries': 1,
    'retry_delay': timedelta(minutes=10),
    'project_id': 'b-materials',
}

# Custom operator to allow templating in overrides and push env vars to Cloud Run Jobs
# class TemplatedCloudRunExecuteJobOperator(CloudRunExecuteJobOperator):
#     template_fields = CloudRunExecuteJobOperator.template_fields + ("overrides",)

def get_target_date(**kwargs):
    """
    Push previous month/year to XCom.
    """
    excec_date = kwargs['ds']
    excec_date_obj = datetime.strptime(excec_date, "%Y-%m-%d")
    first_of_month = excec_date_obj.replace(day=1)
    previous_month_date = first_of_month - timedelta(days=1)
    kwargs['ti'].xcom_push(key='prev_year', value=previous_month_date.year)
    kwargs['ti'].xcom_push(key='prev_month', value=previous_month_date.month)


with DAG(
    dag_id='in_out_analytics_pipeline',
    default_args=default_args,
    schedule_interval='30 15 1 * *', #NOTE: UTC Timestamp == 09:30 on first day of the month AM CST
    catchup= False,
    tags=['in_out_pipeline', 'elt'],
    description= 'This DAG orchestrates the ELT process for the In & Out Transactional System' \
    'using the Cloud Run Jobs: in-out-analytics-pipeline and the in-out-analytics-dbt-job.'
) as dag:

    start = EmptyOperator(task_id='start')

    #Get month and year task for EL process
    get_target_dag_date = PythonOperator(
        task_id='get_target_dag_date',
        python_callable=get_target_date,
        provide_context=True,
    )

    #Cloud Run Job
    # Trigger EL process (extract from Firestore, load to Bronze Landing Table in BigQuery)
    # This job can be overriden with a $YEAR and $MONTH env vars **For backfilling purposes**
    # If not set, it will use the current year and month like (2025, 6)
    trigger_el_job = CloudRunExecuteJobOperator(
        task_id='trigger_el_job',
        region='us-central1',
        project_id='b-materials',
        job_name='in-out-analytics-pipeline',
        overrides={
            "container_overrides": [
                {
                    "env": [
                        {
                            "name": "YEAR",
                            "value": "{{ ti.xcom_pull(task_ids='get_target_dag_date', key='prev_year') }}"
                        },
                        {
                            "name": "MONTH",
                            "value": "{{ ti.xcom_pull(task_ids='get_target_dag_date', key='prev_month') }}"
                        }
                    ]
                }
            ]
        }
    )

    #CLoud Run Job
    # Trigger DBT process (tranform data from Bronze to Silver in BigQuery)
    # This job can be overriden with a $DBT_MODEL and $DBT_TARGET env vars 
    # like ($DBT_MODEL='in_out_bronze', $DBT_TARGET='bronze') **For running other models**
    # If not set, it will use the default env vars (in_out_silver, silver)
    trigger_dbt_job_silver = CloudRunExecuteJobOperator(
        task_id='trigger_dbt_job_silver',
        region='us-central1',
        project_id='b-materials',
        job_name='in-out-analytics-dbt-job',
        overrides={
            "container_overrides":[
                {
                    "env": [
                        {"name": "DBT_MODEL", "value": "in_out_silver"},
                        {"name": "DBT_TARGET", "value": "silver"}
                    ]
                }
            ]
        }
    )

    end = EmptyOperator(task_id='end')

start >> get_target_dag_date >> trigger_el_job >> trigger_dbt_job_silver >> end
