from google.cloud import firestore
import pandas as pd


class FetchDocuments():
    def __init__(self, project_id: str):
        self.project_id = project_id
        self.firestore_client = firestore.Client(project=project_id)

    def fetch_firestore_documents(self, collection_name: str) -> pd.DataFrame:
        collection_ref = self.firestore_client.collection(collection_name)
        docs = collection_ref.stream()
        data = [doc.to_dict() for doc in docs]
        return pd.DataFrame(data)