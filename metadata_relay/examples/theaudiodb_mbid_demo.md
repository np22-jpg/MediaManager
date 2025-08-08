# TheAudioDB MBID Support Demo

This demonstrates the new MBID (MusicBrainz ID) support in TheAudioDB endpoints with prioritized caching.

## New Endpoints

### Artist by MBID
```bash
# Get artist information by MusicBrainz ID
curl "http://localhost:8080/theaudiodb/artist/5b11f4ce-a62d-471e-81fc-a69a8278c7da"
```

**Response example:**
```json
{
  "name": "Nirvana",
  "biography": "Nirvana was an American rock band...",
  "website": "https://www.nirvana.com",
  "genre": "Grunge",
  "formedYear": "1987",
  "mbid": "5b11f4ce-a62d-471e-81fc-a69a8278c7da",
  "thumb": "https://www.theaudiodb.com/images/media/artist/thumb/rvvnqt1347913617.jpg",
  "fanart": "https://www.theaudiodb.com/images/media/artist/fanart/spvrvu1347913675.jpg",
  "country": "United States"
}
```

### Album by MBID
```bash
# Get album information by MusicBrainz release group ID
curl "http://localhost:8080/theaudiodb/album/1b022e01-4da6-387b-8658-8678046e4cef"
```

**Response example:**
```json
{
  "name": "Nevermind",
  "artist": "Nirvana",
  "year": "1991",
  "genre": "Grunge",
  "thumb": "https://www.theaudiodb.com/images/media/album/thumb/uxrqxy1347913577.jpg",
  "description": "Nevermind is the second studio album by American rock band Nirvana...",
  "mbid": "1b022e01-4da6-387b-8658-8678046e4cef"
}
```

### Track by MBID
```bash
# Get track information by MusicBrainz recording ID
curl "http://localhost:8080/theaudiodb/track/f1b10b1e-c2c6-4ff1-bd0a-b3e36f57e0d1"
```

**Response example:**
```json
{
  "name": "Smells Like Teen Spirit",
  "artist": "Nirvana",
  "album": "Nevermind",
  "duration": "301000",
  "genre": "Grunge",
  "description": "Smells Like Teen Spirit is a song by American rock band Nirvana...",
  "thumb": "https://www.theaudiodb.com/images/media/track/thumb/abc123.jpg",
  "mbid": "f1b10b1e-c2c6-4ff1-bd0a-b3e36f57e0d1"
}
```

### Existing Name Search (for comparison)
```bash
# Traditional name-based search
curl "http://localhost:8080/theaudiodb/artist?name=Nirvana"
```

## Caching Strategy

The MBID endpoints use **prioritized caching** with different TTL values:

- **MBID Artist Lookup**: 7-day TTL (168 hours) - High priority caching
- **MBID Album Lookup**: 7-day TTL (168 hours) - High priority caching  
- **MBID Track Lookup**: 7-day TTL (168 hours) - High priority caching
- **Name Search**: 1-day TTL (24 hours) - Standard caching

### Why Prioritized Caching?

1. **MBID Stability**: MusicBrainz IDs are stable identifiers that rarely change
2. **Data Quality**: MBID-based lookups return more accurate and complete data
3. **Performance**: Longer cache TTL reduces API calls to TheAudioDB
4. **Integration**: MBID-based endpoints are ideal for integration with MusicBrainz data

## Integration Example

Here's how you might use the MBID endpoints in conjunction with MusicBrainz:

```bash
# 1. Search for an artist in MusicBrainz
curl "http://localhost:8080/musicbrainz/artists/search?query=Radiohead&limit=1"

# 2. Extract the MBID from the MusicBrainz response
# Example MBID: a74b1b7f-71a5-4011-9441-d0b5e4122711

# 3. Get additional information from TheAudioDB using the MBID
curl "http://localhost:8080/theaudiodb/artist/a74b1b7f-71a5-4011-9441-d0b5e4122711"

# 4. Get album information if you have a release group MBID
curl "http://localhost:8080/theaudiodb/album/{release-group-mbid}"

# 5. Get track information if you have a recording MBID
curl "http://localhost:8080/theaudiodb/track/{recording-mbid}"
```

## Cache Performance Monitoring

You can monitor cache performance via the metrics endpoint:

```bash
curl "http://localhost:8080/metrics" | grep theaudiodb
```

Look for metrics like:
- `cache_hits_total{provider="theaudiodb",operation="artist_mbid"}`
- `cache_hits_total{provider="theaudiodb",operation="album_mbid"}`
- `cache_hits_total{provider="theaudiodb",operation="track_mbid"}`
- `cache_misses_total{provider="theaudiodb",operation="artist_mbid"}`
- `cache_misses_total{provider="theaudiodb",operation="album_mbid"}`
- `cache_misses_total{provider="theaudiodb",operation="track_mbid"}`
- `api_requests_total{provider="theaudiodb"}`
