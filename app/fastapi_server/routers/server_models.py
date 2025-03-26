from pydantic import BaseModel
from typing import Optional, Any

class ErrorDetail(BaseModel):
    type: str
    message: str

class StandardResp(BaseModel):
    success: bool
    data: Optional[Any] = None
    error: Optional[ErrorDetail] = None


class ServerData(BaseModel):
    id: str
    url: str
    name: str
    type: str
    version: str
    status: str

class CreationStatus(BaseModel):
    server_id: str

class ServerCreationReq(BaseModel):
    name: str
    version: str
    type: str

class ServerCreationResp(StandardResp):
    data: Optional[CreationStatus] = None

class ServerFetchReq(BaseModel):
    server_id: str

class ServerFetchResp(StandardResp):
    data: Optional[ServerData] = None

class ServerStartReq(BaseModel):
    server_id: str

class ServerStartResp(StandardResp):
    data: Optional[ServerData] = None
