# Troubleshooting

<note>
    Note the lack of a trailing slash in some env vars like FRONTEND_URL. This is important.
</note>

<tip>
    Always check the container and browser logs for more specific error messages
</tip>

## Authentication Issues (OIDC)

* Verify `OPENID_CLIENT_ID`, `OPENID_CLIENT_SECRET`, and `OPENID_CONFIGURATION_ENDPOINT` are correct.
* Ensure the `FRONTEND_URL` is accurate and that your OpenID provider has the correct redirect URI whitelisted (
  e.g., `http://your-frontend-url/api/v1/auth/cookie/Authentik/callback`).


## Data Not Appearing / File Issues

* For hardlinks to work, you must not use different docker volumes for TV, Torrents, etc.