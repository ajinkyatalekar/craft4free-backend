# https://craft4free.online 
This repo contains the python codebase for the backend for craft4free. This is hosted on an oracle cloud instance. It runs traefik inside a docker container, which routes traffic to a fastapi server.

## Technologies
- Fastapi
- Python-on-whales (Docker management)
- Supabase

## Structure
- `app/fastapi_server` contains the fastapi server which is copied to a docker container when run.
- `app/scripts` contains the scripts for handling docker deployment, management, and removal of minecraft server containers.
