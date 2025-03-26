import os
from supabase import create_client, Client
import dotenv

dotenv.load_dotenv()

url: str = os.getenv("SUPABASE_URL", "")
key: str = os.environ.get("SUPABASE_SERVICE_ROLE_KEY", "")

supabase: Client = create_client(url, key)
