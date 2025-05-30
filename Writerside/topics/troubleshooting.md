# Troubleshooting

<note>
    Note the lack of a trailing slash in some env vars like FRONTEND_URL. This is important.
</note>

<tip>
    Always check the container logs for more specific error messages
</tip>

## Authentication Issues

    * Double-check `AUTH_TOKEN_SECRET`. If it changes, existing sessions/tokens will be invalidated.
    * For OpenID:
        * Verify `OPENID_CLIENT_ID`, `OPENID_CLIENT_SECRET`, and `OPENID_CONFIGURATION_ENDPOINT` are correct.
        * Ensure the `FRONTEND_URL` is accurate and that your OpenID provider has the correct redirect URI whitelisted (
          e.g., `http://your-frontend-url/api/v1/auth/cookie/Authentik/callback`).
        * Check backend logs for errors from `httpx_oauth` or `fastapi-users`.

## CORS Errors

    * Ensure `FRONTEND_URL` is correctly set.
    * Ensure your frontend's url is listed in `CORS_URLS`.

## Data Not Appearing / File Issues

    * Verify that the volume mounts for `IMAGE_DIRECTORY`, `TV_DIRECTORY`, `MOVIE_DIRECTORY`, and `TORRENT_DIRECTORY` in
      your `docker-compose.yml` are correctly pointing to your media folders on the host machine.
    * Check file and directory permissions for the user running the Docker container (or the `node` user inside the
      containers).
