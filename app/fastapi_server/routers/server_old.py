from fastapi import APIRouter
from fastapi import APIRouter, Depends, BackgroundTasks, status
from fastapi import status
from fastapi.responses import JSONResponse
from fastapi_server.core.security import verify_token
from routers.server_models import CreationStatus, ServerCreationReq, ServerCreationResp, ServerData, ServerStartResp, ServerStartReq, ErrorDetail
from config.supabase import supabase
# from scripts.server.handler import start_server, stop_server
# import datetime

router = APIRouter()

@router.post("/server", status_code=status.HTTP_201_CREATED, response_model=ServerCreationResp)
async def create_new_server(request: ServerCreationReq, user = Depends(verify_token)):
    try:

        existing_server = supabase.table("servers").select("*").eq("user_id", "1").eq("name", request.name).execute()
        if existing_server.data:
            raise Exception("Server with the same name already exists")

        resp = (supabase.table("servers").insert({
            "user_id": "1",
            "name": request.name,
            "version": request.version,
            "type": request.type
        }).execute())
        resp = resp.model_dump()
        id = resp['data'][0]['id']

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

# @router.get("/server/{server_id}")
# async def get_server(server_id: str, user = Depends(verify_token)):
#     try:
#         response = (supabase.table("servers")
#             .select("*")
#             .eq("id", server_id)
#             .eq("user_id", user["sub"])
#             .execute())

#         return response
#     except Exception as e:
#         print(f"Error fetching server: {e}")
#         return {"error": "Error fetching server"}

# @router.post("/server/{server_id}/start", status_code=status.HTTP_200_OK, response_model=ServerStartResp)
# async def start_server_(request: ServerStartReq, user = Depends(verify_token)):
#     try:
#         port = start_server(request.server_id, "test", "VANILLA", "1.21.4")

#         if (port == -1):
#             raise Exception("Failed to start server")

#         return ServerStartResp(
#                     success=True,
#                     data=ServerData(
#                         id=request.server_id,
#                         url=f"129.213.144.81:{port}",
#                         name="",
#                         type="VANILLA",
#                         version="1.21.4",
#                         status="RUNNING"
#                     )
#                 )

#     except Exception as e:
#         print(f"Error starting server: {e}")
#         return JSONResponse(
#             status_code=status.HTTP_500_INTERNAL_SERVER_ERROR,
#             content=ServerCreationResp(
#                 success=False,
#                 error=ErrorDetail(
#                     type="SERVER_START_ERROR",
#                     message="Server start failed"
#                 )
#             ).dict()
#         )

# @router.post("/server/{server_id}/stop")
# async def stop_server_(server_id: str, background_tasks: BackgroundTasks, user = Depends(verify_token)):
#     try:
#         background_tasks.add_task(stop_server, server_id)

#         return {"message": "Server stopped"}
#     except Exception as e:
#         print(f"Error stopping server: {e}")
#         return {"error": "Error stopping server"}

# @router.post("/server/{server_id}/delete")
# async def delete_server_(server_id: str, user = Depends(verify_token)):
#     try:
#         response = (supabase.table("servers")
#             .update({"deleted_at": datetime.datetime.now(datetime.timezone.utc).isoformat()})
#             .eq("id", server_id)
#             .eq("user_id", user["sub"])
#             .execute())

#         return response
#     except Exception as e:
#         print(f"Error fetching server: {e}")
#         return {"error": "Error fetching server"}
