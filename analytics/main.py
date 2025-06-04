from fetch import FetchDocuments
from dotenv import load_dotenv
import os



def run_elt():
    #Load environment variables
    load_dotenv()
    project_id = os.getenv("PROJECT_ID")
    fetcher = FetchDocuments(project_id)
    
    #Fetch Entradas
    firestore_entradas_data = fetcher.fetch_firestore_documents("entradas") 
    print("Firestore Entradas Data:")
    print(firestore_entradas_data.head())

    #Fetch Salidas
    firestore_salidas_data = fetcher.fetch_firestore_documents("salidas") 
    print("Firestore Salidas Data:")
    print(firestore_salidas_data.head())

    #Load Entradas to BigQuery


    #load Salidas to Bigquery
    
