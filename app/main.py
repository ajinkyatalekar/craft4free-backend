from fastapi import FastAPI
from fastapi.middleware.cors import CORSMiddleware
from fastapi_server.routers import server

app = FastAPI()

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
