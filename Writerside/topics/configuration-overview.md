# Installation Guide

The recommended way to install and run Media Manager is using Docker and Docker Compose.

1. **Prerequisites:**
    * Ensure Docker and Docker Compose are installed on your system.
    * If you plan to use OAuth 2.0 / OpenID Connect for authentication, you will need an account and client credentials
      from an OpenID provider (e.g., Authentik, Pocket ID).

2. **Setup:**
    * Copy the docker-compose.yml from the MediaManager repo.
    * Configure the necessary environment variables in your `docker-compose.yml` file.
    * (Optional) Create a `.env` file in the root directory for backend environment variables and/or a `web/.env` for
      frontend environment variables if you prefer to manage them separately from `docker-compose.yml`.

3. **Running the Application:**
    * Execute the command `docker-compose up -d` from the root directory. This will build the Docker images (if not
      already built) and start all the services (backend, frontend, and potentially a database if configured in your
      compose file).
    * The backend will typically be available at `http://localhost:8000` and the frontend at `http://localhost:3000` (or
      as configured).

# Configuration Overview

Media Manager is configured primarily through environment variables. These can be set in your `docker-compose.yml` file,
a `.env` file.

Detailed configuration options are split into backend and frontend sections:

* [Backend Configuration](configuration-backend.md)
* [Frontend Configuration](configuration-frontend.md)

Build arguments are also used during the Docker image build process, primarily for versioning.

