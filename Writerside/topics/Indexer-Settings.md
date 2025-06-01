# Indexer Settings

## Prowlarr

| Variable           | Description                         | Default                 | Example                |
|--------------------|-------------------------------------|-------------------------|------------------------|
| `PROWLARR_ENABLED` | Set to `True` to enable Prowlarr.   | `False`                 | `true`                 |
| `PROWLARR_API_KEY` | Your Prowlarr API key.              | -                       | `prowlarr_api_key`     |
| `PROWLARR_URL`     | Base URL of your Prowlarr instance. | `http://localhost:9696` | `http://prowlarr:9696` |

## Jackett

| Variable           | Description                                        | Default                 | Example                |
|--------------------|----------------------------------------------------|-------------------------|------------------------|
| `JACKETT_ENABLED`  | Set to `True` to enable Jackett.                   | `False`                 | `true`                 |
| `JACKETT_API_KEY`  | Your Prowlarr API key.                             | -                       | `jackett_api_key`      |
| `JACKETT_URL`      | Base URL of your Prowlarr instance.                | `http://localhost:9117` | `http://prowlarr:9117` |
| `JACKETT_INDEXERS` | list of all indexers for Jackett to search through | `all`                   | `["1337x","0magnet"]`  |

<note>
    <include from="notes.topic" element-id="list-format"/>
</note>