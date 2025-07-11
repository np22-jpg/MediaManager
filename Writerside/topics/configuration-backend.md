# Backend
These settings configure the core backend application through the `config.toml` file. All backend configuration is now centralized in this TOML file instead of environment variables.

## General Settings (`[misc]`)

- `frontend_url`

The URL the frontend will be accessed from. This is a required field and must include the trailing slash.

- `cors_urls`

A list of origins you are going to access the API from. Note the lack of trailing slashes.

- `api_base_path`

The URL base path of the backend API. Default is `/api/v1`. Note the lack of a trailing slash.

- `development`

Set to `true` to enable development mode. Default is `false`.

## Example Configuration

Here's a complete example of the general settings section in your `config.toml`:

```toml
[misc]
# REQUIRED: Change this to match your actual frontend URL
frontend_url = "http://localhost:3000/"

# REQUIRED: List all origins that will access the API
cors_urls = ["http://localhost:3000", "http://localhost:8000"]

# Optional: API base path (rarely needs to be changed)
api_base_path = "/api/v1"

# Optional: Development mode (set to true for debugging)
development = false
```

<note>
    The <code>frontend_url</code> and <code>cors_urls</code> are the most important settings to configure correctly. Make sure they match your actual deployment URLs.
</note>
