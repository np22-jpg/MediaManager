# Configuration

MediaManager uses a TOML configuration file (`config.toml`) for all backend settings.
This centralized configuration approach makes it easier to manage, backup, and share your MediaManager setup.

Frontend settings are configured through environment variables in your `docker-compose.yaml` file.

## Configuration File Location

Your `config.toml` file should be mounted to `/data/config.toml` inside the container:

```yaml
volumes:
  - ./config.toml:/data/config.toml
```

Though you can change the location of the configuration file by setting the `CONFIG_FILE` env variable to another path,
e.g. `/etc/mm/config.toml`.

## Configuration Sections

The configuration is organized into the following sections:

- `[misc]` - General settings
- `[database]` - Database settings
- `[auth]` - Authentication settings
- `[notifications]` - Notification settings (Email, Gotify, Ntfy, Pushover)
- `[torrents]` - Download client settings (qBittorrent, SABnzbd)
- `[indexers]` - Indexer settings (Prowlarr and Jackett )
- `[metadata]` - TMDB and TVDB settings

## Configuring Secrets

For sensitive information like API keys, passwords, and secrets, you _should_ use environment variables.
You can actually set every configuration value through environment variables.
For example, to set the `token_secret` value for authentication, with a .toml file you would use:

```toml
[auth]
token_secret = "your_super_secret_key_here"
```

But you can also set it through an environment variable:

```
AUTH_TOKEN_SECRET = "your_super_secret_key_here"
```

or another example with the OIDC client secret:

```toml
[auth]
...
[auth.openid_connect]
client_secret = "your_client_secret_from_provider"
```

env variable:

```
AUTH_OPENID_CONNECT_CLIENT_SECRET = "your_client_secret_from_provider"
```

So for every config "level", you basically have to take the name of the value and prepend it with the section names in
uppercase with underscores as delimiters.


