# Backend

These settings configure the core backend application through the `config.toml` file. All backend configuration is now
centralized in this TOML file instead of environment variables.

## General Settings (`[misc]`)

- `frontend_url`

The URL the frontend will be accessed from. This is a required field and must include the trailing slash. The default
path is `http://localhost:8000/web/`. Make sure to change this to match your actual frontend URL.

- `cors_urls`

A list of origins you are going to access the API from. Note the lack of trailing slashes.

- `development`

Set to `true` to enable development mode. Default is `false`.

## Example Configuration

Here's a complete example of the general settings section in your `config.toml`:

```toml
[misc]
# REQUIRED: Change this to match your actual frontend URL
frontend_url = "http://localhost:8000/web/"

cors_urls = ["http://localhost:8000"]

# Optional: Development mode (set to true for debugging)
development = false
```

<note>
    The <code>frontend_url</code> is the most important settings to configure correctly. Make sure it matches your actual deployment URLs.
</note>