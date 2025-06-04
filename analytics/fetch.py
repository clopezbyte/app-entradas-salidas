from google.cloud import firestore
from google.cloud.firestore_v1 import FieldFilter
from datetime import datetime
import pandas as pd
import os
import logging
import sys


# Configure Logging
logging.basicConfig(level=logging.INFO, format='%(asctime)s - %(levelname)s - %(message)s', stream=sys.stdout)
logger = logging.getLogger('in-out-analytics')


class FetchDocuments():
    def __init__(self, project_id: str):
        self.project_id = project_id
        self.firestore_client = firestore.Client(
            project=project_id,
            database="app-in-out-good"
        )

    def fetch_firestore_documents(self, collection_name: str, timestamp_field:str, year: int, month: int) -> pd.DataFrame:
        """Query a Firestore NoSQL database collection.
        
        Keyword arguments:
        collection_name: a str representing the collection name.
        timestamp_field: a str representing the name of the field containing a timestamp in the document.
        year: integer representing the target year.
        month: integer representing the target month.
        Return: pd.Dataframe
        """
        
        start = datetime(year, month, 1)
        end = datetime(year + (month == 12), (month % 12) + 1, 1)

        query = self.firestore_client.collection(collection_name)\
        .where(filter=FieldFilter(timestamp_field, ">=", start))\
        .where(filter=FieldFilter(timestamp_field, "<", end))

        try:
            docs = query.stream()
            data = [doc.to_dict() for doc in docs]
            return pd.DataFrame(data)
        except Exception as e:
            logger.exception("Failed to fetch Firestore documents.")
            raise