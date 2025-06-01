# Indexer Settings

## Prowlarr

### `PROWLARR_ENABLED`

Set to `True` to enable Prowlarr. Default is `False`. Example: `true`.

### `PROWLARR_API_KEY`

This is your Prowlarr API key. Example: `prowlarr_api_key`.

### `PROWLARR_URL`

Base URL of your Prowlarr instance. Default is `http://localhost:9696`. Example: `http://prowlarr:9696`.

## Jackett

### `JACKETT_ENABLED`

Set to `True` to enable Jackett. Default is `False`. Example: `true`.

### `JACKETT_API_KEY`

This is your Prowlarr API key. Example: `jackett_api_key`.

### `JACKETT_URL`

Base URL of your Prowlarr instance. Default is `http://localhost:9117`. Example: `http://prowlarr:9117`.

### `JACKETT_INDEXERS`

List of all indexers for Jackett to search through. Default is `all`. Example: `["1337x","0magnet"]`.

<note>
    <include from="notes.topic" element-id="list-format"/>
</note>