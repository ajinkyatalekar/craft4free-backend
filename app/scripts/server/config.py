import os

# Environment variables
HOST_PWD = os.environ.get('HOST_PWD', os.getcwd())
SERVER_HOST_IP = os.environ.get('SERVER_HOST_IP', '129.213.144.81')

# Docker settings
MINECRAFT_IMAGE = "itzg/minecraft-server"
DEFAULT_MEMORY = "512M"
MINECRAFT_PORT = 25565

# Server settings
DATA_DIR_TEMPLATE = "{base_dir}/data/servers/{server_id}"
