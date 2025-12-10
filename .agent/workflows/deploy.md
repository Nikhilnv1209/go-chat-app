---
description: Build and deploy the application using Docker Compose
---

1. Stop any running local instances (server or db)
// turbo
2. Build and start the containers
   ```bash
   podman-compose up --build -d
   ```
3. Check status
   ```bash
   podman-compose ps
   ```
4. View logs
   ```bash
   podman-compose logs -f app
   ```
