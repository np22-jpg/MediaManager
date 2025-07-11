# Reverse Proxy

MediaManager is probably unlike any other service that you have deployed before, because it has a separate frontend and backend container.
This means that setting up a reverse proxy is a bit different.

When deploying MediaManager, you have two choices:

- Use two different hostnames, e.g. `mediamanager.example.com` for the frontend and `api-mediamanager.example.com` for the backend.
- Use a single hostname with a base path, e.g. `mediamanager.example.com` for the frontend and `mediamanager.example.com/api/v1` for the backend.

If you choose the first option, you can set up your reverse proxy as usual, forwarding requests to the appropriate containers based on the hostname.
If you choose the second option, you need to ensure that your reverse proxy is configured to handle the base path correctly.

## Example Caddy Configuration

```
mm.my-domain.com {
    @api path /api/*
    reverse_proxy @api 10.0.0.7:8000

    reverse_proxy 10.0.0.7:3000
}
```

## Example Traefik Configuration
This example assumes you use Traefik with Docker labels.

```yaml

services:
  backend:
    image: ghcr.io/maxdorninger/mediamanager/backend:latest
    ports:
      - "8000:8000"
    labels:
      - "traefik.enable=true"
      - "traefik.http.routers.mediamanager-backend.rule=Host(`media.example`)&&PathPrefix(`/api/v1`)"
      - "traefik.http.routers.mediamanager-backend.tls=true"
      - "traefik.http.routers.mediamanager-backend.tls.certresolver=letsencrypt"
      - "traefik.http.routers.mediamanager-backend.entrypoints=websecure"
    environment:
      - MISC_FRONTEND_URL=https://media.example/
  frontend:
    image: ghcr.io/maxdorninger/mediamanager/frontend:latest
    labels:
      - "traefik.enable=true"
      - "traefik.http.routers.mediamanager-frontend.rule=Host(`media.example`)"
      - "traefik.http.routers.mediamanager-frontend.tls=true"
      - "traefik.http.routers.mediamanager-frontend.tls.certresolver=letsencrypt"
      - "traefik.http.routers.mediamanager-frontend.entrypoints=websecure"
    ports:
      - "3000:3000"
    environment:
      - PUBLIC_API_URL=https://media.example/api/v1
```
