# Frontend

## Environment Variables

### `PUBLIC_API_URL`

You (the browser) must reach the backend from this url. Default is `http://localhost:8000/api/v1`. Example:
`https://mediamanager.example.com/api/v1`.

## Build Arguments (web/Dockerfile)

**TODO: expand on this section**

Unfortunately you need to build the frontend docker container, to configure a url base path for the frontend. This is because
SvelteKit needs to know the base path at build time.

### `BASE_URL`

Sets the base url path, it must begin with a slash and not end with one. Example (in build command):
- clone the repo
- cd into the repo's root directory
- `docker build --build-arg BASE_URL=/media -f web/Dockerfile .`

### `VERSION`

Sets the `PUBLIC_VERSION` environment variable at runtime in the frontend container. Passed during build. Example (in
build command): `docker build --build-arg VERSION=1.2.3 -f web/Dockerfile .`