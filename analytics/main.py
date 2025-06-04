from fetch import FetchDocuments
from dotenv import load_dotenv
from datetime import datetime
import os
import logging
import sys


# Configure Logging
logging.basicConfig(level=logging.INFO, format='%(asctime)s - %(levelname)s - %(message)s', stream=sys.stdout)
logger = logging.getLogger('in-out-analytics')

def parse_int_env(var_name: str, default: int) -> int:
    val = os.getenv(var_name)
    if val and val.strip().isdigit():
        return int(val)
    return default

def run_elt(year: int, month: int):
    # Load environment variables
    load_dotenv()
    project_id = os.getenv("PROJECT_ID")

    if not project_id:
        logger.error("Missing PROJECT_ID environment variable.")
        sys.exit(1)

    fetcher = FetchDocuments(project_id)

    logger.info(f"Fetching data for year={year}, month={month}")

    # Fetch Entradas
    entradas_df = fetcher.fetch_firestore_documents("entradas", "FechaRecepcion", year, month)
    logger.info(f"Fetched {len(entradas_df)} 'entradas' records.")
    # print(entradas_df.head())

    # Fetch Salidas
    salidas_df = fetcher.fetch_firestore_documents("salidas", "FechaSalida", year, month)
    logger.info(f"Fetched {len(salidas_df)} 'salidas' records.")
    # print(salidas_df.head())

    # TODO: Load to BigQuery
    # load_to_bigquery(entradas_df, ...)
    # load_to_bigquery(salidas_df, ...)


if __name__ == "__main__":
    load_dotenv()

    try:
        year = int(sys.argv[1]) if len(sys.argv) > 1 else parse_int_env("YEAR", datetime.now().year)
        month = int(sys.argv[2]) if len(sys.argv) > 2 else parse_int_env("MONTH", datetime.now().month)
    except Exception as e:
        logger.error(f"Failed to parse year/month: {e}")
        sys.exit(1)

    run_elt(year, month)
