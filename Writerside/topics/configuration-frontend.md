# Frontend

| Variable             | Description                                                    | Default                        | Example                                   |
|----------------------|----------------------------------------------------------------|--------------------------------|-------------------------------------------|
| `PUBLIC_WEB_SSR`     | Enables/disables Server-Side Rendering. (this is experimental) | `false`                        | `true`                                    |
| `PUBLIC_API_URL`     | You (the browser) mut reach the backend from this url.         | `http://localhost:8000/api/v1` | `https://mediamanager.example.com/api/v1` |
| `PUBLIC_SSR_API_URL` | The frontent container must reach the backend from this url.   | `http://localhost:8000/api/v1` | `http://backend:8000/api/v1`              |

## Build Arguments (web/Dockerfile)

| Argument  | Description                                                                                               | Example (in build command)                                   |
|-----------|-----------------------------------------------------------------------------------------------------------|--------------------------------------------------------------|
| `VERSION` | Sets the `PUBLIC_VERSION` environment variable at runtime in the frontend container. Passed during build. | `docker build --build-arg VERSION=1.2.3 -f web/Dockerfile .` |

