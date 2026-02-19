# Host Lotus
Host Lotus is a minecraft server hosting platform. This repository is the GO-based backend for the platform.

## Controller
The controller accepts API calls from the frontend, manages the persistent database, manages the workers, and sends server CRUD requests to the workers.

## Worker
The worker accepts commands from the controller. It runs servers on the local machine, and sends updates and heartbeats to the controller.