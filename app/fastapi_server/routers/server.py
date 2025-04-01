from fastapi import APIRouter, Depends, status, BackgroundTasks
from fastapi.responses import JSONResponse
import datetime

from fastapi_server.core.security import verify_token
from fastapi_server.routers.server_models import ServerCreationReq, ServerCreationResp, CreationStatus, ErrorDetail, ServerStartReq, ServerStartResp, ServerData, StandardResp

from config.supabase import supabase
from scripts.server.handler import start_server, stop_server
from scripts.server.info import get_server_status

router = APIRouter()

@router.post("/server", status_code=status.HTTP_201_CREATED, response_model=ServerCreationResp)
async def create_new_server(request: ServerCreationReq, user = Depends(verify_token)):
    try:
        existing_server = supabase.table("servers").select("*").eq("user_id", user["sub"]).eq("name", request.name).execute()
        if existing_server.data:
            raise Exception("Server with the same name already exists")

        resp = (supabase.table("servers").insert({
            "user_id": user["sub"],
            "name": request.name,
            "version": request.version,
            "type": request.type
        }).execute())
        id = resp.data[0]['id']

        return ServerCreationResp(
            data=CreationStatus(server_id=id),
            success=True,
        )
    except Exception as e:
        print(f"Error creating server: {e}")
        return JSONResponse(
            status_code=status.HTTP_500_INTERNAL_SERVER_ERROR,
            content=ServerCreationResp(
                success=False,
                error=ErrorDetail(
                    type="SERVER_CREATION_ERROR",
                    message=f"{e}"
                )
            ).dict()
        )


@router.post("/server/{server_id}/start", status_code=status.HTTP_200_OK, response_model=ServerStartResp)
async def start_server_(request: ServerStartReq, user = Depends(verify_token)):
    try:
        resp = (supabase.table("servers").select("*").eq("id", request.server_id).single().execute())
        server_type = resp.data["type"]
        server_version = resp.data["version"]
        server_name = resp.data["name"]
    except Exception as e:
        print(f"Error fetching server: {e}")
        return JSONResponse(
            status_code=status.HTTP_500_INTERNAL_SERVER_ERROR,
            content=ServerStartResp(
                success=False,
                error=ErrorDetail(
                    type="SERVER_FETCH_ERROR",
                    message=f"{e}"
                )
            ).dict()
        )

    try:
        port = start_server(request.server_id, server_name, server_type, server_version)

        if (port == -1):
            raise Exception("Failed to start server")

        return ServerStartResp(
                    success=True,
                    data=ServerData(
                        id=request.server_id,
                        url=f"129.213.144.81:{port}",
                        name=server_name,
                        type=server_type,
                        version=server_version,
                        status="RUNNING"
                    )
                )
    except Exception as e:
        print(f"Error starting server: {e}")
        return JSONResponse(
            status_code=status.HTTP_500_INTERNAL_SERVER_ERROR,
            content=ServerStartResp(
                success=False,
                error=ErrorDetail(
                    type="SERVER_START_ERROR",
                    message="Server start failed: "+str(e)
                )
            ).dict()
        )

@router.post("/server/{server_id}/stop")
async def stop_server_(server_id: str, background_tasks: BackgroundTasks, user = Depends(verify_token)):
    try:
        background_tasks.add_task(stop_server, server_id)

        return {"message": "Server stopped"}
    except Exception as e:
        print(f"Error stopping server: {e}")
        return {"error": "Error stopping server"}

@router.post("/server/{server_id}/delete")
async def delete_server_(server_id: str, user = Depends(verify_token)):
    try:
        response = (supabase.table("servers")
            .update({"deleted_at": datetime.datetime.now(datetime.timezone.utc).isoformat()})
            .eq("id", server_id)
            .eq("user_id", user["sub"])
            .execute())

        return response
    except Exception as e:
        print(f"Error fetching server: {e}")
        return {"error": "Error fetching server"}

@router.get("/server")
async def get_all_servers(user = Depends(verify_token)):
    try:
        response = (supabase.table("servers")
            .select("*")
            .eq("user_id", user["sub"])
            .is_("deleted_at", None)
            .execute())

        servers = []
        for server in response.data:
            servers.append({
                "server": server,
                "status": get_server_status(server["id"]),
            })
        return StandardResp(
            success=True,
            data=servers
        )

    except Exception as e:
        print(f"Error fetching servers: {e}")
        return StandardResp(
            success=False,
            data=[]
        )

@router.post("/server/{server_id}")
async def get_server_(server_id: str, user = Depends(verify_token)):
    try:
        resp = (supabase.table("servers").select("*").eq("id", server_id).single().execute())
    except Exception as e:
        print(f"Error fetching server: {e}")
        return {"error": "Error fetching server"}

    try:
        return StandardResp(
            success=True,
            data={
                "status": get_server_status(server_id),
                "server": resp.data
            }
        )
    except Exception as e:
        print(f"Error fetching server: {e}")
        return JSONResponse(
            status_code=status.HTTP_500_INTERNAL_SERVER_ERROR,
            content=ServerCreationResp(
                success=False,
                error=ErrorDetail(
                    type="SERVER_START_ERROR",
                    message="Server start failed: " + str(e)
                )
            ).dict()
        )
