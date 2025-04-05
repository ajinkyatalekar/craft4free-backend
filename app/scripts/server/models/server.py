from enum import Enum
from pydantic import BaseModel
from typing import Optional

class ServerStatus(str, Enum):
    RUNNING = "running"
    STARTING = "starting"
    STOPPED = "stopped"
    UNKNOWN = "unknown"

class ServerInfo(BaseModel):
    port: Optional[str] = None
    url: Optional[str] = None
    status: ServerStatus = ServerStatus.STOPPED
    error: Optional[str] = None

class ServerConfig(BaseModel):
    id: str
    name: str
    type: str
    version: str
    memory: str = "512M"
    motd: Optional[str] = None
    online_mode: bool = True
