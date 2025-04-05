from python_on_whales import docker
import logging
from typing import Tuple, Optional, Dict, Any
from scripts.server.models.server import ServerStatus
from scripts.server.config import MINECRAFT_PORT, MINECRAFT_IMAGE

logger = logging.getLogger(__name__)

class DockerService:
    """Service for interacting with Docker containers running the Minecraft server. Consumed by ServerService."""

    @staticmethod
    def container_exists(container_id: str) -> bool:
        """Check if a container exists."""
        try:
            return docker.container.exists(container_id)
        except Exception as e:
            logger.error(f"Error checking container existence: {e}")
            return False

    @staticmethod
    def get_container_status(container_id: str) -> Tuple[ServerStatus, str]:
        """Get the status of a container."""
        try:
            if not DockerService.container_exists(container_id):
                return ServerStatus.STOPPED, ""

            container = docker.container.inspect(container_id)

            if container.state.status == "exited":
                return ServerStatus.STOPPED, ""

            if container.state.status == "running" and container.state.health:
                health_status = container.state.health.status
                if health_status == "starting":
                    return ServerStatus.STARTING, ""
                if health_status == "healthy":
                    return ServerStatus.RUNNING, ""

            return ServerStatus.UNKNOWN, "Unknown container health status"

        except Exception as e:
            logger.error(f"Error getting container status: {e}")
            return ServerStatus.UNKNOWN, str(e)

    @staticmethod
    def get_container_port(container_id: str) -> Optional[str]:
        """Get the host port for a container."""
        try:
            if not DockerService.container_exists(container_id):
                return None

            container = docker.container.inspect(container_id)
            if container.state.status != "running":
                return None

            ports = container.network_settings.ports
            if ports and f"{MINECRAFT_PORT}/tcp" in ports:
                return ports[f"{MINECRAFT_PORT}/tcp"][0]["HostPort"]
            return None

        except Exception as e:
            logger.error(f"Error getting container port: {e}")
            return None

    @staticmethod
    def run_container(container_id: str, env_vars: Dict[str, Any], volume_path: str) -> Optional[str]:
        """Run a new container with the given parameters."""
        try:
            if DockerService.container_exists(container_id):
                docker.remove(container_id)

            docker.run(
                MINECRAFT_IMAGE,
                detach=True,
                interactive=True,
                tty=True,
                publish=[(0, MINECRAFT_PORT)],
                envs=env_vars,
                name=container_id,
                volumes=[(volume_path, "/data")]
            )

            return DockerService.get_container_port(container_id)

        except Exception as e:
            logger.error(f"Error running container: {e}")
            return None

    @staticmethod
    def stop_container(container_id: str) -> bool:
        """Stop and remove a container."""
        try:
            if DockerService.container_exists(container_id):
                docker.stop(container_id)
                docker.remove(container_id)
            return True
        except Exception as e:
            logger.error(f"Error stopping container: {e}")
            return False
