# Troubleshooting

<note>
    Note the lack of a trailing slash in some env vars like FRONTEND_URL. This is important.
</note>

<tip>
    Always check the container and browser logs for more specific error messages
</tip>

<procedure title="I can't log in with OAuth/OIDC?" id="procedure-i-cannot-log-in-with-oauth">
    <step>Verify your OAuth provider's configuration. <a href="authentication-setup.md" anchor="openid-connect-settings-auth-openid-connect">See the OAuth documentation</a></step>
    <step>Check the frontend url in your config file. It should match the URL you use to access MediaManager.</step>
</procedure>

<procedure title="I cannot log in?" id="procedure-i-cannot-log-in">
   <step>Make sure you are logging in, not signing up.</step>
   <step>Try logging in with the following credentials: 
        <list>
          <li>Email: admin@mediamanager.local or admin@example.com</li>
          <li>Password: admin</li>
        </list>
   </step>
</procedure>

<procedure title="My hardlinks don't work?" id="procedure-my-hardlinks-dont-work">
   <step>Make sure you are using only one volumes for TV, Movies and Downloads. <a href="https://raw.githubusercontent.com/maxdorninger/MediaManager/refs/heads/master/docker-compose.yaml"> See the configuration in the example <code>docker-compose.yaml</code> file.</a></step>
</procedure>

<note>If it still doesn't work, <a href="https://github.com/maxdorninger/MediaManager/issues">please open an Issue.</a> It is possible that a bug is causing the issue.</note>