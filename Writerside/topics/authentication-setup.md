# Authentication

MediaManager supports multiple authentication methods. Email/password authentication is the default, but you can also
enable OpenID Connect (OAuth 2.0) for integration with external identity providers.

<note>
   Note the lack of a trailing slash in some env vars like FRONTEND_URL. This is important.
</note>

## General Authentication Settings

### `AUTH_TOKEN_SECRET`

Strong secret key for signing JWTs (create with `openssl rand -hex 32`). This is a required field. Example:
`AUTH_TOKEN_SECRET=your_super_secret_key`.

### `AUTH_SESSION_LIFETIME`

Lifetime of user sessions in seconds. Default is `86400` (1 day). Example: `AUTH_SESSION_LIFETIME=604800` (1 week).

### `AUTH_ADMIN_EMAIL`

A list of email addresses for administrator accounts. This is a required field. Example:
`AUTH_ADMIN_EMAIL=admin@example.com`.

### `FRONTEND_URL`

The URL the frontend will be accessed from. This is a required field. Example: `https://mediamanager.example`.

<note>
On login/registration, every user whose email is in <code>AUTH_ADMIN_EMAIL</code> will be granted admin privileges.
Users whose email is not in <code>AUTH_ADMIN_EMAIL</code> will be regular users and will need to be verified by an administrator,
this can be done in the settings page.
</note>
<tip>
    <include from="notes.topic" element-id="list-format"/>
</tip>

## OpenID Connect (OAuth 2.0)

### `OPENID_ENABLED`

Enables OpenID authentication. Default is `FALSE`. Example: `TRUE`.

### `OPENID_CLIENT_ID`

Client ID from your OpenID provider.

### `OPENID_CLIENT_SECRET`

Client Secret from your OpenID provider.

### `OPENID_CONFIGURATION_ENDPOINT`

URL of your OpenID provider's discovery document (e.g., `.../.well-known/openid-configuration`). Example:
`https://authentik.example.com/application/o/mediamanager/.well-known/openid-configuration`.

### `OPENID_NAME`

Display name for this OpenID provider. Default is `OpenID`. Example: `Authentik`.

### Configuring OpenID Connect

1. Set `OPENID_ENABLED=TRUE`
2. Configure the following environment variables:
    * `OPENID_CLIENT_ID`
    * `OPENID_CLIENT_SECRET`
    * `OPENID_CONFIGURATION_ENDPOINT`
    * `OPENID_NAME` (optional)
    * `FRONTEND_URL` (it is important that this is set correctly, as it is used for the redirect URIs)
3. Your OpenID server will likely want a redirect URI. This URL will be like:
   `{FRONTEND_URL}/api/v1/auth/cookie/{OPENID_NAME}/callback`. The exact path depends on the `OPENID_NAME`.

4. Example URL: `https://mediamanager.example/api/v1/auth/cookie/Authentik/callback`

