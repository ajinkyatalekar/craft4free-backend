from pydantic import BaseModel
from typing import Optional, Any

class StandardResponse(BaseModel):
    success: bool
    data: Optional[Any] = None
    error: Optional[str] = None

class ServerCreateRequest(BaseModel):
    name: str
    version: str
    type: str
