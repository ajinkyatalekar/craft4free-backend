import os
from fastapi import HTTPException, Security
from fastapi.security import HTTPBearer, HTTPAuthorizationCredentials
import jwt
from jwt.exceptions import PyJWTError
import dotenv

dotenv.load_dotenv()

JWT_SECRET = os.getenv("SUPABASE_JWT_SECRET")
JWT_ISSUER = os.getenv("SUPABASE_URL_v1")
JWT_AUDIENCE = "authenticated"

security = HTTPBearer()

async def verify_token(credentials: HTTPAuthorizationCredentials = Security(security)):
    token = credentials.credentials
    try:
        # Properly verify the token with signature validation
        payload = jwt.decode(
            token,
            JWT_SECRET,
            algorithms=["HS256"],
            options={
                "verify_signature": True,
                "verify_exp": True,
                "verify_aud": True,
            },
            issuer=JWT_ISSUER,
            audience=JWT_AUDIENCE
        )

        if not payload.get("sub"):
            raise HTTPException(status_code=401, detail="Invalid user in token")

        return payload
    except PyJWTError as e:
        print(f"Error verifying token: {e}")
        raise HTTPException(status_code=401, detail="Invalid authentication token")
