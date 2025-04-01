from python_on_whales import docker
import os

host_pwd = os.environ.get('HOST_PWD')

def get_server_port(server_id: str):
    if docker.container.exists(server_id) and docker.container.inspect(server_id).state.status == "running":
        ports = docker.container.inspect(server_id).network_settings.ports
        if ports:
            return ports["25565/tcp"][0]["HostPort"]
    return ""

def get_server_status(server_id: str):
    resp = {
        "port": "",
        "url": "",
        "status": "",
        "error": ""
    }

    # Server does not exist
    if not docker.container.exists(server_id):
        resp["status"] = "stopped"
        return resp

    container = docker.container.inspect(server_id)

    # Server not running
    if container.state.status == "exited":
        resp["status"] = "stopped"
        return resp

    if container.state.status == "running" and container.state.health:
        resp["port"] = get_server_port(server_id)
        resp["url"] = f"129.213.144.81:{resp['port']}"

        # Server starting
        if container.state.health.status == "starting":
            resp["status"] = "starting"
            return resp

        # Server running
        if container.state.health.status == "healthy":
            resp["status"] = "running"
            return resp

    resp["error"] = "Unknown status"
    return resp
