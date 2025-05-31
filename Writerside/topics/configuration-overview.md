# Installation Guide

The recommended way to install and run Media Manager is using Docker and Docker Compose.

## Prerequisites

* Ensure Docker and Docker Compose are installed on your system.
* If you plan to use OAuth 2.0 / OpenID Connect for authentication, you will need an account and client credentials
  from an OpenID provider (e.g., Authentik, Pocket ID).

## Setup

* Download the docker-compose.yaml from the MediaManager repo with the following command:
  ```
  wget -o docker-compose.yaml https://raw.githubusercontent.com/maxdorninger/MediaManager/refs/heads/master/docker-compose.yaml
  ```

* Configure the necessary environment variables in your `docker-compose.yaml` file.
* For more information on the available configuration options, see the [Configuration section](Configuration.md) of the
  documentation.

<note>
   It is good practice to put API keys and other sensitive information in a separate `.env` file and reference them in your
  `docker-compose.yaml`.
</note>


