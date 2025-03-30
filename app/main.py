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
    allow_origins=[
        # Prod
        "https://craft4free.online",
         "https://www.craft4free.online",
        # Dev
        "https://dev.d2w3788e6h0dd9.amplifyapp.com"
    ],
    allow_methods=["GET", "POST", "PUT", "DELETE", "OPTIONS"],
    allow_headers=["Content-Type", "Authorization"],
)

app.include_router(server.router)

@app.get("/")
def read_root():
    return {
        "status": "online",
        "api_version": "0.0.1"
    }
