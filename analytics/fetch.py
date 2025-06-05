from google.cloud import firestore
from google.cloud.firestore_v1 import FieldFilter
from datetime import datetime
import pandas as pd
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
            data = []
            for doc in docs:
                doc_dict = doc.to_dict()
                doc_dict['id'] = doc.id  
                data.append(doc_dict)
            return pd.DataFrame(data)
        except Exception as e:
            logger.exception("Failed to fetch Firestore documents.")
            raise

    def prepare_entradas_data(self, data: pd.DataFrame) -> pd.DataFrame:
        try:
            data = data[[
                "id", "BodegaRecepcion", "Cantidad", "Cliente", "FechaRecepcion", "FechaAjusteASN",
                "TipoDelivery", "PersonaRecepcion", "ProveedorRecepcion", "Type"
            ]]

            expected_columns = {"id", "BodegaRecepcion", "Cantidad", "Cliente", 
                                "FechaRecepcion", "FechaAjusteASN", "TipoDelivery", "PersonaRecepcion", 
                                "ProveedorRecepcion", "Type"}
            missing = expected_columns - set(data.columns)
            if missing:
                logger.error(f"Missing expected columns: {missing}")
                return pd.DataFrame()

            data.rename(columns={
                "id": "landing_movement_id",
                "BodegaRecepcion": "bodega",
                "Cantidad": "cantidad",
                "Cliente": "cliente",
                "FechaRecepcion": "fecha_movimiento",
                "FechaAjusteASN": "fecha_ajuste_asn",
                "TipoDelivery": "tipo_delivery",
                "PersonaRecepcion": "operador",
                "ProveedorRecepcion": "proveedor",
                "Type": "tipo"
            }, inplace=True)
            
            return data
        except Exception as e:
            logger.exception("Error")
            return pd.DataFrame()
        
    def prepare_salidas_data(self, data: pd.DataFrame) -> pd.DataFrame:
        try:
            data = data[[
                "id", "BodegaSalida", "Cliente", "FechaSalida", "PersonaEntrega", "ProveedorSalida",
                "Type"
            ]]

            expected_columns = {"id", "BodegaSalida", "Cliente", "FechaSalida",
                                "PersonaEntrega", "ProveedorSalida", "Type"}
            missing = expected_columns - set(data.columns)
            if missing:
                logger.error(f"Missing expected columns: {missing}")
                return pd.DataFrame()

            data.rename(columns={
                "id": "landing_movement_id",
                "BodegaSalida": "bodega",
                "Cliente": "cliente",
                "FechaSalida": "fecha_movimiento",
                "PersonaEntrega": "operador",
                "ProveedorSalida": "proveedor",
                "Type": "tipo"
            }, inplace=True)

            data["cantidad"] = 0 #Not tracked
            data["fecha_ajuste_asn"] = None #Does not apply
            data["tipo_delivery"] = None #Does not apply

            return data
        except Exception as e:
            logger.exception("Error")
            return pd.DataFrame()
