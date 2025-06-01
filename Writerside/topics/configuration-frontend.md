# Frontend

## Environment Variables

### `PUBLIC_WEB_SSR`

Enables/disables Server-Side Rendering. (this is experimental). Default is `false`. Example: `true`.

### `PUBLIC_API_URL`

You (the browser) must reach the backend from this url. Default is `http://localhost:8000/api/v1`. Example:
`https://mediamanager.example.com/api/v1`.

### `PUBLIC_SSR_API_URL`

The frontend container must reach the backend from this url. Default is `http://localhost:8000/api/v1`. Example:
`http://backend:8000/api/v1`.

## Build Arguments (web/Dockerfile)

**TODO: expand on this section**

To configure a url base path for the frontend, you need to build the frontend docker container, this is because
unfortunately SvelteKit needs to know the base path at build time.

### `VERSION`

Sets the `PUBLIC_VERSION` environment variable at runtime in the frontend container. Passed during build. Example (in
build command): `docker build --build-arg VERSION=1.2.3 -f web/Dockerfile .`

### `BASE_URL`

Sets the base url path, it must begin with a slash and not end with one. Example (in build command):
`docker build --build-arg BASE_URL=/media -f web/Dockerfile .`
