from scripts.server.services.docker_service import DockerService
from scripts.server.models.server import ServerInfo, ServerConfig
from scripts.server.config import SERVER_HOST_IP, DATA_DIR_TEMPLATE, HOST_PWD
from typing import Tuple, Optional
import os

class ServerService:
    """Service for managing the Minecraft server."""

    @staticmethod
    def get_server_info(server_id: str) -> ServerInfo:
        """Get comprehensive information about a server."""
        status, error = DockerService.get_container_status(server_id)
        port = DockerService.get_container_port(server_id)

        return ServerInfo(
            status=status,
            port=port,
            url=f"{SERVER_HOST_IP}:{port}" if port else None,
            error=error
        )

    @staticmethod
    def start_server(config: ServerConfig) -> Tuple[bool, Optional[str]]:
        """Start a server with the given configuration."""

        # Create data directory
        data_dir = DATA_DIR_TEMPLATE.format(base_dir=HOST_PWD, server_id=config.id)
        os.makedirs(data_dir, exist_ok=True)

        # Prepare environment variables
        env_vars = {
            "EULA": "TRUE",
            "MEMORY": config.memory,
            "MOTD": config.motd or config.name,
            "VERSION": config.version,
            "TYPE": config.type,
            "ONLINE_MODE": str(config.online_mode).lower()
        }

        # Start container
        port = DockerService.run_container(config.id, env_vars, data_dir)
        if port:
            return True, port
        return False, "Failed to start server"

    @staticmethod
    def stop_server(server_id: str) -> bool:
        """Stop a server."""
        return DockerService.stop_container(server_id)
