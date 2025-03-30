from python_on_whales import docker
import os

host_pwd = os.environ.get('HOST_PWD')

def start_server(server_id: str, name: str, type: str, version: str):
    '''
    Run a docker container with the minecraft server image
    '''

    data_dir = f"{host_pwd}/data/servers/{server_id}"
    os.makedirs(data_dir, exist_ok=True)

    if docker.container.exists(server_id):
        docker.remove(server_id)

    docker.run(
        "itzg/minecraft-server",
        detach=True,
        interactive=True,
        tty=True,
        publish=[(0, 25565)],
        envs={
            "EULA": "TRUE",
            "MEMORY": "512M",
            "MOTD": f"{name}",
            "VERSION": version,
            "TYPE": type,
            "ONLINE_MODE": "true"
        },
        name=server_id,
        volumes=[(data_dir, "/data")]
    )

    ports = docker.container.inspect(server_id).network_settings.ports

    if ports:
        return ports["25565/tcp"][0]["HostPort"]

    return -1

def stop_server(server_id: str):
    '''
    Stop and remove container
    '''
    if docker.container.exists(server_id):
        docker.stop(server_id)
        docker.remove(server_id)
