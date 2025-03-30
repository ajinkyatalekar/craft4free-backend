from fastapi import FastAPI
from fastapi.middleware.cors import CORSMiddleware
from fastapi_server.routers import server

# You might need to add this to your startup script if you have permission issues
import subprocess

# Check if Docker socket permissions need to be adjusted (if not running as root)
try:
    subprocess.run(["docker", "info"], check=True)
except:
    # Try to fix permissions if you have sudo access
    subprocess.run(["sudo", "chmod", "666", "/var/run/docker.sock"])

app = FastAPI(openapi_url=None)

app.add_middleware(
    CORSMiddleware,
    allow_credentials=True,
    allow_origins=["*"],
    allow_methods=["*"],
    allow_headers=["*"],
)

app.include_router(server.router)

@app.get("/")
def read_root():
    return {
        "status": "online",
        "api_version": "0.0.1",
        "docs": "/docs",
    }
