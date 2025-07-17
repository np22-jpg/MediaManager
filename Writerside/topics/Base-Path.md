# Base Path

Unfortunately you need to build the docker container, to configure a url base path. This is because
SvelteKit needs to know the base path at build time.

## Docker Build Arguments

When building the Docker image for the frontend, you can pass build arguments to set the base path and version.

### `BASE_PATH`

Sets the base url path, it must begin with a slash and not end with one. Example (in build command):
- clone the repo
- cd into the repo's root directory
- `docker build --build-arg BASE_URL=/media -f Dockerfile .`

### `VERSION`

Sets the version variable in the container. This isn't strictly necessary, but it can be useful for debugging or versioning purposes.
Example: `docker build --build-arg VERSION=locally-built -f web/Dockerfile .`