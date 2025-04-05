from fastapi import APIRouter, Depends, status, WebSocket, WebSocketDisconnect, BackgroundTasks
from fastapi.responses import JSONResponse
import datetime
import asyncio
from python_on_whales import docker

from scripts.server.models.server import ServerConfig
from config.supabase import supabase
from fastapi_server.core.security import verify_token
from fastapi_server.models.server import ServerCreateRequest, StandardResponse

from scripts.server.services.server_service import ServerService

router = APIRouter()

@router.post("/new/servers/{server_id}")
async def get_server(server_id: str, user = Depends(verify_token)):
    try:
        resp = (supabase.table("servers").select("*").eq("id", server_id).single().execute())
    except Exception as e:
        return StandardResponse(
            success=False,
            error=str(e)
        )

    try:
        return StandardResponse(
            success=True,
            data={
                "status": ServerService.get_server_info(server_id),
                "server": resp.data
            }
        )
    except Exception as e:
        return JSONResponse(
            status_code=status.HTTP_500_INTERNAL_SERVER_ERROR,
            content=StandardResponse(
                success=False,
                error=str(e)
            ).dict()
        )

@router.get("/new/servers")
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
                "status": ServerService.get_server_info(server["id"]),
            })
        return StandardResponse(
            success=True,
            data=servers
        )
    except Exception as e:
        return StandardResponse(
            success=False,
            error=str(e)
        )

@router.post("/new/servers")
async def create_server(request: ServerCreateRequest, user = Depends(verify_token)):
    try:
        existing_server = supabase.table("servers").select("*").eq("user_id", user["sub"]).eq("name", request.name).execute()
        if existing_server.data:
            raise Exception("Server with the same name already exists")

        (supabase.table("servers").insert({
                "user_id": user["sub"],
                "name": request.name,
                "version": request.version,
                "type": request.type
            }).execute())

        return StandardResponse(
            success=True
        )
    except Exception as e:
        return JSONResponse(
            status_code=status.HTTP_500_INTERNAL_SERVER_ERROR,
            content=StandardResponse(
                success=False,
                error=str(e)
            ).dict()
        )

@router.post("/new/servers/{server_id}/delete")
async def delete_server(server_id: str, user = Depends(verify_token)):
    try:
        (supabase.table("servers")
            .update({"deleted_at": datetime.datetime.now(datetime.timezone.utc).isoformat()})
            .eq("id", server_id)
            .eq("user_id", user["sub"])
            .execute())

        return StandardResponse(
            success=True
        )
    except Exception as e:
        return StandardResponse(
            success=False,
            error=str(e)
        )

@router.post("/new/servers/{server_id}/start")
async def start_server(server_id: str, user = Depends(verify_token)):
    try:
        resp = (supabase.table("servers").select("*").eq("id", server_id).single().execute())

        config = ServerConfig(
            id = resp.data["id"],
            name = resp.data["name"],
            type = resp.data["type"],
            version = resp.data["version"]
        )
    except Exception as e:
        print(f"Error fetching server: {e}")
        return JSONResponse(
            status_code=status.HTTP_500_INTERNAL_SERVER_ERROR,
            content=StandardResponse(
                success=False,
                error=str(e)
            ).dict()
        )

    try:
        ok, port = ServerService.start_server(config)

        if not ok:
            raise Exception("Failed to start server")

        return StandardResponse(
            success=True,
        )
    except Exception as e:
        print(f"Error starting server: {e}")
        return JSONResponse(
            status_code=status.HTTP_500_INTERNAL_SERVER_ERROR,
            content=StandardResponse(
                success=False,
                error=str(e)
            ).dict()
        )

@router.post("/new/servers/{server_id}/stop")
async def stop_server(server_id: str, background_tasks: BackgroundTasks, user = Depends(verify_token)):
    try:
        background_tasks.add_task(ServerService.stop_server, server_id)

        return StandardResponse(
            success=True,
        )
    except Exception as e:
        print(f"Error stopping server: {e}")
        return JSONResponse(
            status_code=status.HTTP_500_INTERNAL_SERVER_ERROR,
            content=StandardResponse(
                success=False,
                error=str(e)
            ).dict()
        )


server_processes = {}
@router.websocket("/new/ws/console/{server_id}")
async def websocket_endpoint(websocket: WebSocket, server_id: str):
    await websocket.accept()

    # Initialize server entry if it doesn't exist
    if server_id not in server_processes:
        server_processes[server_id] = {
            'connected_clients': set()
        }

    # Add this client to the set of connected clients
    server_processes[server_id]['connected_clients'].add(websocket)

    try:
        # Get the container ID for this server
        container_id = server_id
        if not container_id:
            await websocket.send_text("[SERVER_NOT_RUNNING]")
        else:
            try:
                # Get the container
                container = docker.container.inspect(container_id)

                # Send historical logs first (last 100 lines)
                logs = container.logs(tail=100).splitlines()
                for line in logs:
                    line_text = line.decode('utf-8').strip() if isinstance(line, bytes) else line.strip()
                    await websocket.send_text(line_text)

                # For python-on-whales, we need to use asyncio subprocess to stream logs
                async def stream_logs():
                    # Use docker CLI directly for streaming
                    process = await asyncio.create_subprocess_shell(
                        f"docker logs -f {container_id}",
                        stdout=asyncio.subprocess.PIPE,
                        stderr=asyncio.subprocess.PIPE
                    )

                    # Check if process started correctly
                    if process.stdout is None:
                        await websocket.send_text("[ERROR] Failed to start log streaming")
                        return

                    while True:
                        line = await process.stdout.readline()
                        if not line:
                            break
                        line_text = line.decode('utf-8').strip()
                        try:
                            await websocket.send_text(line_text)
                        except Exception:
                            # WebSocket closed
                            process.kill()
                            break

                # Start stream logs task
                task = asyncio.create_task(stream_logs())

                # Keep the websocket open
                try:
                    while True:
                        await asyncio.sleep(1)  # Simple keep-alive approach
                except WebSocketDisconnect:
                    task.cancel()

            except Exception as e:
                await websocket.send_text(f"[ERROR] {str(e)}")
                await websocket.send_text("[SERVER_NOT_RUNNING]")
    except WebSocketDisconnect:
        # Remove this client from the set of connected clients
        if server_id in server_processes:
            server_processes[server_id]['connected_clients'].discard(websocket)
