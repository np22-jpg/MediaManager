<br />
<div align="center">
  <a href="https://maxdorninger.github.io/MediaManager">
    <img src="https://raw.githubusercontent.com/maxdorninger/MediaManager/refs/heads/master/Writerside/images/logo.svg" alt="Logo" width="260" height="260">
  </a>

<h3 align="center">MediaManager</h3>

  <p align="center">
    Modern management system for your media library
    <br />
    <a href="https://maxdorninger.github.io/MediaManager/introduction.html"><strong>Explore the docs »</strong></a>
    <br />
    <a href="https://maxdorninger.github.io/MediaManager/issues/new?labels=bug&template=bug-report---.md">Report Bug</a>
    &middot;
    <a href="https://maxdorninger.github.io/MediaManager/issues/new?labels=enhancement&template=feature-request---.md">Request Feature</a>
  </p>
</div>


MediaManager is modern software to manage your TV and movie library. It is designed to be a replacement for Sonarr,
Radarr, Overseer, and Jellyseer.
It supports TVDB and TMDB for metadata, supports OIDC and OAuth 2.0 for authentication and supports Prowlarr and
Jackett.
MediaManager is built first and foremost for deployment with Docker, making it easy to set up.

It also provides an API to interact with the software programmatically, allowing for automation and integration with
other services.

## Quick Start

```
   wget -O docker-compose.yaml https://raw.githubusercontent.com/maxdorninger/MediaManager/refs/heads/master/docker-compose.yaml   
   # Edit docker-compose.yaml to set the environment variables!
   docker compose up -d
```

### [View the docs for installation instructions and more](https://maxdorninger.github.io/MediaManager/configuration-overview.html#configuration-overview)

<!-- ROADMAP -->
## Roadmap

- [x] support for more torrent indexers
- [x] fully automatic downloads
- [x] add tests
- [x] add more logs/errors
- [x] make API return proper error codes
- [x] optimize images for web in the backend
- [x] responsive ui
- [x] automatically update metadata of shows
- [x] automatically download new seasons/episodes of shows
- [x] add fallback to just copy files if hardlinks don't work
- [x] add check at startup if hardlinks work
- [x] create separate metadata relay service, so that api keys for TMDB and TVDB are not strictly needed
- [x] support for movies
- [x] expand README with more information and a quickstart guide
- [x] improve reliability of scheduled tasks
- [ ] add notification system
- [ ] add in-depth documentation on the architecture of the codebase
- [ ] make indexer module multithreaded
- [ ] add support for deluge and transmission
- [ ] add delete button for movies/TV shows
- [ ] rework prowlarr module (select which indexers to use, etc.)
- [ ] add sequence diagrams to the documentation
- [ ] _maybe_ rework the logo
- [ ] _maybe_ add support for configuration via toml config file

See the [open issues](hhttps://maxdorninger.github.io/MediaManager/issues) for a full list of proposed features (and known issues).


<!-- LICENSE -->
## License

Distributed under the GPL 3.0. See `LICENSE.txt` for more information.


<!-- ACKNOWLEDGMENTS -->
## Acknowledgments

* [Thanks to Pawel Czerwinski for the image on the login screen](https://unsplash.com/@pawel_czerwinski)

